// Package draw provides SVG primitive shapes, text elements, and styling options.
// It supports circles, ellipses, rectangles, lines, and paths with a fluent builder.
package draw

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/dominionthedev/finite/geom"
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
	fill        string // hex color or url(#id)
	fillDef     string // optional gradient/pattern def to inline
	stroke      string
	strokeWidth float64
	opacity     float64
	filter      *Filter
	clip        *ClipPath
	mask        *Mask
	transform   string // raw SVG transform attribute value
	title       string // accessible name (<title> child)
	desc        string // accessible description (<desc> child)
	ariaLabel   string // aria-label attribute
	role        string // ARIA role
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
	if s.clip != nil {
		parts = append(parts, s.clip.Def())
	}
	if s.mask != nil {
		parts = append(parts, s.mask.Def())
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
	if s.clip != nil {
		b.WriteString(fmt.Sprintf(` clip-path="%s"`, s.clip.Ref()))
	}
	if s.mask != nil {
		b.WriteString(fmt.Sprintf(` mask="%s"`, s.mask.Ref()))
	}
	if s.transform != "" {
		b.WriteString(fmt.Sprintf(` transform="%s"`, s.transform))
	}
	b.WriteString(a11yAttrs(s.role, s.ariaLabel))
	return b.String()
}

// wrap wraps SVG content in a <g> with optional inline defs.
func wrap(defs, content string) string {
	if defs == "" {
		return "<g>" + content + "</g>"
	}
	return "<g>" + defs + content + "</g>"
}

// elementMarkup builds a shape element; uses open/close form when title/desc are set.
func elementMarkup(tag, attrs, title, desc string) string {
	body := a11yContent(title, desc)
	if body == "" {
		return "<" + tag + attrs + "/>"
	}
	return "<" + tag + attrs + ">" + body + "</" + tag + ">"
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

// Transform applies a geom.Matrix transform to the shape.
// For convenience, a raw SVG transform string is still accepted via TransformString.
func Transform(m geom.Matrix) ShapeOption {
	return func(s *shapeStyle) { s.transform = m.String() }
}

// TransformString applies a raw SVG transform string to the shape.
// Prefer Transform(geom.Matrix) for type-safe composition.
func TransformString(t string) ShapeOption {
	return func(s *shapeStyle) { s.transform = t }
}

// WithClip attaches a clipPath to the shape.
func WithClip(c *ClipPath) ShapeOption {
	return func(s *shapeStyle) { s.clip = c }
}

// WithMask attaches a mask to the shape.
func WithMask(m *Mask) ShapeOption {
	return func(s *shapeStyle) { s.mask = m }
}

// FillPattern sets a pattern fill. The pattern definition is inlined like gradients.
func FillPattern(p *Pattern) ShapeOption {
	return func(s *shapeStyle) {
		s.fill = p.Ref()
		s.fillDef = p.Def()
	}
}

// Title sets the accessible name as a <title> child element.
func Title(text string) ShapeOption {
	return func(s *shapeStyle) { s.title = text }
}

// Desc sets the accessible description as a <desc> child element.
func Desc(text string) ShapeOption {
	return func(s *shapeStyle) { s.desc = text }
}

// AriaLabel sets the aria-label attribute.
func AriaLabel(label string) ShapeOption {
	return func(s *shapeStyle) { s.ariaLabel = label }
}

// Role sets the ARIA role attribute (e.g. "img", "presentation").
func Role(role string) ShapeOption {
	return func(s *shapeStyle) { s.role = role }
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
	attrs := fmt.Sprintf(` cx="%.4f" cy="%.4f" r="%.4f"%s`, c.cx, c.cy, c.r, c.style.attrString())
	el := elementMarkup("circle", attrs, c.style.title, c.style.desc)
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
	attrs := fmt.Sprintf(` cx="%.4f" cy="%.4f" rx="%.4f" ry="%.4f"%s`, e.cx, e.cy, e.rx, e.ry, e.style.attrString())
	el := elementMarkup("ellipse", attrs, e.style.title, e.style.desc)
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
	attrs := fmt.Sprintf(` x="%.4f" y="%.4f" width="%.4f" height="%.4f"%s%s`,
		r.x, r.y, r.w, r.h, rx, r.style.attrString())
	el := elementMarkup("rect", attrs, r.style.title, r.style.desc)
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
	attrs := fmt.Sprintf(` x1="%.4f" y1="%.4f" x2="%.4f" y2="%.4f"%s`, l.x1, l.y1, l.x2, l.y2, l.style.attrString())
	el := elementMarkup("line", attrs, l.style.title, l.style.desc)
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
	attrs := fmt.Sprintf(` d="%s"%s`, strings.TrimSpace(p.d.String()), p.style.attrString())
	el := elementMarkup("path", attrs, p.style.title, p.style.desc)
	return wrap(p.style.buildDefs(), el), nil
}
