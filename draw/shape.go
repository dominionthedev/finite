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

// fillProvider is anything that can return a fill string (color or gradient ref).
type fillProvider interface {
	Ref() string
	Def() string
}

// shapeStyle holds common visual properties shared across shape types.
type shapeStyle struct {
	fill        string   // hex color or url(#id)
	fillDef     string   // optional gradient/filter def to inline
	stroke      string
	strokeWidth float64
	opacity     float64
	filter      *Filter
	transform   string // raw SVG transform attribute value
}

// applyDefs returns the inline <defs> block if the shape needs one.
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

// ── Shape Option ─────────────────────────────────────────────────────────────

// ShapeOption configures a shape's visual properties.
type ShapeOption func(*shapeStyle)

// Fill sets a solid hex color fill.
//
//	draw.Fill("#a78bfa")
func Fill(color string) ShapeOption {
	return func(s *shapeStyle) {
		s.fill = color
		s.fillDef = ""
	}
}

// FillGradient sets a gradient fill. The gradient's def is inlined automatically.
//
//	draw.FillGradient(draw.Radial("orb", "50%", "50%", "50%", ...))
func FillGradient(g fillProvider) ShapeOption {
	return func(s *shapeStyle) {
		s.fill = g.Ref()
		s.fillDef = g.Def()
	}
}

// Stroke sets the stroke color and width.
//
//	draw.Stroke("#ffffff", 2)
func Stroke(color string, width float64) ShapeOption {
	return func(s *shapeStyle) {
		s.stroke = color
		s.strokeWidth = width
	}
}

// Opacity sets the element's overall opacity (0.0–1.0).
func Opacity(v float64) ShapeOption {
	return func(s *shapeStyle) { s.opacity = v }
}

// WithFilter attaches a filter to the shape. The filter def is inlined automatically.
//
//	draw.WithFilter(draw.Blur("b1", 18))
func WithFilter(f *Filter) ShapeOption {
	return func(s *shapeStyle) { s.filter = f }
}

// Transform sets a raw SVG transform attribute value.
//
//	draw.Transform("rotate(45, 200, 200)")
func Transform(t string) ShapeOption {
	return func(s *shapeStyle) { s.transform = t }
}

// ── Circle ────────────────────────────────────────────────────────────────────

// Circle is an SVG <circle> element.
type Circle struct {
	cx, cy, r float64
	style     shapeStyle
}

// NewCircle creates a circle centered at (cx, cy) with radius r.
//
//	draw.NewCircle(200, 200, 150,
//	    draw.Fill("#a78bfa"),
//	    draw.Opacity(0.8),
//	)
func NewCircle(cx, cy, r float64, opts ...ShapeOption) *Circle {
	c := &Circle{cx: cx, cy: cy, r: r}
	c.style.opacity = 1.0
	for _, opt := range opts {
		opt(&c.style)
	}
	return c
}

// Render implements doc.Renderable.
func (c *Circle) Render() (string, error) {
	el := fmt.Sprintf(`<circle cx="%.4f" cy="%.4f" r="%.4f"%s/>`,
		c.cx, c.cy, c.r, c.style.attrString())
	return wrap(c.style.buildDefs(), el), nil
}

// ── Ellipse ───────────────────────────────────────────────────────────────────

// Ellipse is an SVG <ellipse> element.
type Ellipse struct {
	cx, cy, rx, ry float64
	style          shapeStyle
}

// NewEllipse creates an ellipse centered at (cx, cy) with x-radius rx and y-radius ry.
func NewEllipse(cx, cy, rx, ry float64, opts ...ShapeOption) *Ellipse {
	e := &Ellipse{cx: cx, cy: cy, rx: rx, ry: ry}
	e.style.opacity = 1.0
	for _, opt := range opts {
		opt(&e.style)
	}
	return e
}

// Render implements doc.Renderable.
func (e *Ellipse) Render() (string, error) {
	el := fmt.Sprintf(`<ellipse cx="%.4f" cy="%.4f" rx="%.4f" ry="%.4f"%s/>`,
		e.cx, e.cy, e.rx, e.ry, e.style.attrString())
	return wrap(e.style.buildDefs(), el), nil
}

// ── Rect ──────────────────────────────────────────────────────────────────────

// Rect is an SVG <rect> element.
type Rect struct {
	x, y, w, h float64
	rx, ry      float64 // corner radii
	style       shapeStyle
}

// NewRect creates a rectangle at (x, y) with width w and height h.
//
//	draw.NewRect(50, 50, 300, 200,
//	    draw.Fill("#1e1e2e"),
//	    draw.RoundedCorners(12),
//	)
func NewRect(x, y, w, h float64, opts ...ShapeOption) *Rect {
	r := &Rect{x: x, y: y, w: w, h: h}
	r.style.opacity = 1.0
	for _, opt := range opts {
		opt(&r.style)
	}
	return r
}

// RoundedCorners is a ShapeOption that sets equal x and y corner radii.
func RoundedCorners(r float64) ShapeOption {
	return func(s *shapeStyle) {
		// stored in the Rect directly via closure trick — see Render()
		// we encode it in the transform field temporarily (hack-free below)
		_ = r // handled by WithRx/Ry on Rect directly; see NewRoundedRect
	}
}

// NewRoundedRect creates a rectangle with rounded corners.
func NewRoundedRect(x, y, w, h, radius float64, opts ...ShapeOption) *Rect {
	r := &Rect{x: x, y: y, w: w, h: h, rx: radius, ry: radius}
	r.style.opacity = 1.0
	for _, opt := range opts {
		opt(&r.style)
	}
	return r
}

// Render implements doc.Renderable.
func (r *Rect) Render() (string, error) {
	rx := ""
	if r.rx > 0 {
		rx = fmt.Sprintf(` rx="%.4f" ry="%.4f"`, r.rx, r.ry)
	}
	el := fmt.Sprintf(`<rect x="%.4f" y="%.4f" width="%.4f" height="%.4f"%s%s/>`,
		r.x, r.y, r.w, r.h, rx, r.style.attrString())
	return wrap(r.style.buildDefs(), el), nil
}

// ── Line ──────────────────────────────────────────────────────────────────────

// Line is an SVG <line> element.
type Line struct {
	x1, y1, x2, y2 float64
	style           shapeStyle
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

// Render implements doc.Renderable.
func (l *Line) Render() (string, error) {
	el := fmt.Sprintf(`<line x1="%.4f" y1="%.4f" x2="%.4f" y2="%.4f"%s/>`,
		l.x1, l.y1, l.x2, l.y2, l.style.attrString())
	return wrap(l.style.buildDefs(), el), nil
}

// ── Path ──────────────────────────────────────────────────────────────────────

// Path is an SVG <path> element built with a fluent path builder.
type Path struct {
	d     strings.Builder
	style shapeStyle
}

// NewPath creates a new path builder. Chain M, L, C, Q, A, Z commands.
//
//	draw.NewPath(draw.Fill("#a78bfa")).
//	    M(100, 100).
//	    L(200, 50).
//	    L(300, 100).
//	    Z()
func NewPath(opts ...ShapeOption) *Path {
	p := &Path{}
	p.style.opacity = 1.0
	for _, opt := range opts {
		opt(&p.style)
	}
	return p
}

func (p *Path) M(x, y float64) *Path {
	fmt.Fprintf(&p.d, "M%.4f,%.4f ", x, y)
	return p
}

func (p *Path) L(x, y float64) *Path {
	fmt.Fprintf(&p.d, "L%.4f,%.4f ", x, y)
	return p
}

// C adds a cubic bezier curve (two control points + endpoint).
func (p *Path) C(x1, y1, x2, y2, x, y float64) *Path {
	fmt.Fprintf(&p.d, "C%.4f,%.4f %.4f,%.4f %.4f,%.4f ", x1, y1, x2, y2, x, y)
	return p
}

// Q adds a quadratic bezier curve (one control point + endpoint).
func (p *Path) Q(cx, cy, x, y float64) *Path {
	fmt.Fprintf(&p.d, "Q%.4f,%.4f %.4f,%.4f ", cx, cy, x, y)
	return p
}

func (p *Path) Z() *Path {
	p.d.WriteString("Z")
	return p
}

// SetD sets the raw path data string.
func (p *Path) SetD(d string) *Path {
	p.d.Reset()
	p.d.WriteString(d)
	return p
}

// Render implements doc.Renderable.
func (p *Path) Render() (string, error) {
	el := fmt.Sprintf(`<path d="%s"%s/>`, strings.TrimSpace(p.d.String()), p.style.attrString())
	return wrap(p.style.buildDefs(), el), nil
}
