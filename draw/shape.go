// Package draw provides SVG primitive shapes, text elements, and styling options.
// It supports circles, ellipses, rectangles, lines, and paths with a fluent builder.
package draw

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// shapeCounter generates unique IDs for inline defs (filters, gradients).
var shapeCounter atomic.Uint64

func nextID() string {
	return fmt.Sprintf("fin-%d", shapeCounter.Add(1))
}

// ── Shared option types ───────────────────────────────────────────────────────

// fillProvider is the interface for types that can provide an SVG fill reference and definition.
// LinearGradient and RadialGradient implement this interface.
type fillProvider interface {
	Ref() string
	Def() string
}

type shapeStyle struct {
	fill        string   // hex color or url(#id)
	fillDef     string   // optional gradient/filter def to inline
	stroke      string
	strokeWidth float64
	opacity     float64
	filter      *Filter
	transform   string // raw SVG transform attribute value
}

// buildDefs returns the inline <defs> block if the shape needs one.
func (s *shapeStyle) buildDefs() string {
	var parts []string
	if s.fillDef != "" {
		parts = append(parts, s.fillDef)
	}
	if s.filter != nil {
		parts = append(parts, s.filter.Def())
	}
	if len(parts) == 0 {
		return ""
	}
	return "<defs>" + strings.Join(parts, "") + "</defs>"
}

// attrString builds common SVG attribute string.
func (s *shapeStyle) attrString() string {
	var b strings.Builder
	if s.fill != "" {
		b.WriteString(fmt.Sprintf(` fill="%s"`, s.fill))
	}
	if s.stroke != "" {
		b.WriteString(fmt.Sprintf(` stroke="%s" stroke-width="%.4f"`, s.stroke, s.strokeWidth))
	}
	if s.opacity > 0 && s.opacity != 1.0 {
		b.WriteString(fmt.Sprintf(` opacity="%.4f"`, s.opacity))
	}
	if s.filter != nil {
		b.WriteString(fmt.Sprintf(` filter="%s"`, s.filter.Ref()))
	}
	if s.transform != "" {
		b.WriteString(fmt.Sprintf(` transform="%s"`, s.transform))
	}
	return b.String()
}

// wrap wraps SVG content in a <g> with optional inline defs.
func wrap(defs, content string) string {
	if defs == "" {
		return "<g>" + content + "</g>"
	}
	return "<g>" + defs + content + "</g>"
}

// ShapeOption defines a functional option for configuring a shape's visual properties.
type ShapeOption func(*shapeStyle)

// Fill sets a solid hex color fill for the shape.
func Fill(color string) ShapeOption {
	return func(s *shapeStyle) {
		s.fill = color
		s.fillDef = ""
	}
}

// FillGradient sets a gradient fill for the shape. The gradient definition is automatically
// inlined into the document's definitions or the element's own definition block.
func FillGradient(g fillProvider) ShapeOption {
	return func(s *shapeStyle) {
		s.fill = g.Ref()
		s.fillDef = g.Def()
	}
}

// Stroke sets the stroke color and width for the shape.
func Stroke(color string, width float64) ShapeOption {
	return func(s *shapeStyle) {
		s.stroke = color
		s.strokeWidth = width
	}
}

// Opacity sets the overall opacity (0.0 to 1.0) of the shape.
func Opacity(v float64) ShapeOption {
	return func(s *shapeStyle) { s.opacity = v }
}

// WithFilter attaches an SVG filter (e.g., Blur, Glow, Shadow) to the shape.
func WithFilter(f *Filter) ShapeOption {
	return func(s *shapeStyle) { s.filter = f }
}

// Transform applies a raw SVG transform string to the shape.
func Transform(t string) ShapeOption {
	return func(s *shapeStyle) { s.transform = t }
}

// Circle represents an SVG <circle> element.
type Circle struct {
	cx, cy, r float64
	style     shapeStyle
}

// NewCircle creates a new circle centered at (cx, cy) with radius r.
func NewCircle(cx, cy, r float64, opts ...ShapeOption) *Circle {
	c := &Circle{cx: cx, cy: cy, r: r}
	c.style.opacity = 1.0
	for _, opt := range opts {
		opt(&c.style)
	}
	return c
}

// Render generates the SVG XML for the circle.
func (c *Circle) Render() (string, error) {
	el := fmt.Sprintf(`<circle cx="%.4f" cy="%.4f" r="%.4f"%s/>`,
		c.cx, c.cy, c.r, c.style.attrString())
	return wrap(c.style.buildDefs(), el), nil
}

// Ellipse represents an SVG <ellipse> element.
type Ellipse struct {
	cx, cy, rx, ry float64
	style          shapeStyle
}

// NewEllipse creates a new ellipse centered at (cx, cy) with x-radius rx and y-radius ry.
func NewEllipse(cx, cy, rx, ry float64, opts ...ShapeOption) *Ellipse {
	e := &Ellipse{cx: cx, cy: cy, rx: rx, ry: ry}
	e.style.opacity = 1.0
	for _, opt := range opts {
		opt(&e.style)
	}
	return e
}

// Render generates the SVG XML for the ellipse.
func (e *Ellipse) Render() (string, error) {
	el := fmt.Sprintf(`<ellipse cx="%.4f" cy="%.4f" rx="%.4f" ry="%.4f"%s/>`,
		e.cx, e.cy, e.rx, e.ry, e.style.attrString())
	return wrap(e.style.buildDefs(), el), nil
}

// Rect represents an SVG <rect> element.
type Rect struct {
	x, y, w, h float64
	rx, ry     float64 // corner radii
	style      shapeStyle
}

// NewRect creates a new rectangle at (x, y) with width w and height h.
func NewRect(x, y, w, h float64, opts ...ShapeOption) *Rect {
	r := &Rect{x: x, y: y, w: w, h: h}
	r.style.opacity = 1.0
	for _, opt := range opts {
		opt(&r.style)
	}
	return r
}

// RoundedCorners is a legacy ShapeOption. Use NewRoundedRect instead.
func RoundedCorners(r float64) ShapeOption {
	return func(s *shapeStyle) {
		_ = r
	}
}

// NewRoundedRect creates a rectangle with rounded corners of the specified radius.
func NewRoundedRect(x, y, w, h, radius float64, opts ...ShapeOption) *Rect {
	r := &Rect{x: x, y: y, w: w, h: h, rx: radius, ry: radius}
	r.style.opacity = 1.0
	for _, opt := range opts {
		opt(&r.style)
	}
	return r
}

// Render generates the SVG XML for the rectangle.
func (r *Rect) Render() (string, error) {
	rx := ""
	if r.rx > 0 {
		rx = fmt.Sprintf(` rx="%.4f" ry="%.4f"`, r.rx, r.ry)
	}
	el := fmt.Sprintf(`<rect x="%.4f" y="%.4f" width="%.4f" height="%.4f"%s%s/>`,
		r.x, r.y, r.w, r.h, rx, r.style.attrString())
	return wrap(r.style.buildDefs(), el), nil
}

// Line represents an SVG <line> element.
type Line struct {
	x1, y1, x2, y2 float64
	style          shapeStyle
}

// NewLine creates a line from (x1, y1) to (x2, y2).
func NewLine(x1, y1, x2, y2 float64, opts ...ShapeOption) *Line {
	l := &Line{x1: x1, y1: y1, x2: x2, y2: y2}
	l.style.opacity = 1.0
	for _, opt := range opts {
		opt(&l.style)
	}
	return l
}

// Render generates the SVG XML for the line.
func (l *Line) Render() (string, error) {
	el := fmt.Sprintf(`<line x1="%.4f" y1="%.4f" x2="%.4f" y2="%.4f"%s/>`,
		l.x1, l.y1, l.x2, l.y2, l.style.attrString())
	return wrap(l.style.buildDefs(), el), nil
}

// Path represents an SVG <path> element. It provides a fluent builder API.
type Path struct {
	d     strings.Builder
	style shapeStyle
}

// NewPath initializes a new path builder.
func NewPath(opts ...ShapeOption) *Path {
	p := &Path{}
	p.style.opacity = 1.0
	for _, opt := range opts {
		opt(&p.style)
	}
	return p
}

// M adds a MoveTo command to the path.
func (p *Path) M(x, y float64) *Path {
	fmt.Fprintf(&p.d, "M%.4f,%.4f ", x, y)
	return p
}

// L adds a LineTo command to the path.
func (p *Path) L(x, y float64) *Path {
	fmt.Fprintf(&p.d, "L%.4f,%.4f ", x, y)
	return p
}

// C adds a Cubic Bezier Curve command to the path.
func (p *Path) C(x1, y1, x2, y2, x, y float64) *Path {
	fmt.Fprintf(&p.d, "C%.4f,%.4f %.4f,%.4f %.4f,%.4f ", x1, y1, x2, y2, x, y)
	return p
}

// Q adds a Quadratic Bezier Curve command to the path.
func (p *Path) Q(cx, cy, x, y float64) *Path {
	fmt.Fprintf(&p.d, "Q%.4f,%.4f %.4f,%.4f ", cx, cy, x, y)
	return p
}

// Z closes the current subpath.
func (p *Path) Z() *Path {
	p.d.WriteString("Z")
	return p
}

// SetD sets the raw path data string. This resets any previous builder commands.
func (p *Path) SetD(d string) *Path {
	p.d.Reset()
	p.d.WriteString(d)
	return p
}

// Render generates the SVG XML for the path.
func (p *Path) Render() (string, error) {
	el := fmt.Sprintf(`<path d="%s"%s/>`, strings.TrimSpace(p.d.String()), p.style.attrString())
	return wrap(p.style.buildDefs(), el), nil
}
