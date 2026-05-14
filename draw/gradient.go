package draw

import (
	"fmt"
	"strings"
)

// Stop is a single color stop in a gradient.
type Stop struct {
	Offset  string  // "0%", "50%", "100%"
	Color   string  // hex color, e.g. "#a78bfa"
	Opacity float64 // 0.0–1.0, default 1.0
}

// S is shorthand for constructing a Stop.
//
//	draw.S("0%", "#a78bfa", 1.0)
func S(offset, color string, opacity float64) Stop {
	return Stop{Offset: offset, Color: color, Opacity: opacity}
}

// GradientOption configures a gradient.
type GradientOption func(*gradientBase)

type gradientBase struct {
	id    string
	stops []Stop
}

// ── Linear Gradient ───────────────────────────────────────────────────────────

// LinearGradient defines a linear gradient along a vector.
type LinearGradient struct {
	gradientBase
	x1, y1, x2, y2 string // "0%"/"100%" or absolute values
}

// Linear creates a LinearGradient from (x1,y1) to (x2,y2).
// id must be unique within the document.
//
//	draw.Linear("brand-grad", "0%", "0%", "100%", "100%",
//	    draw.S("0%", "#a78bfa", 1),
//	    draw.S("100%", "#ec4899", 1),
//	)
func Linear(id, x1, y1, x2, y2 string, stops ...Stop) *LinearGradient {
	return &LinearGradient{
		gradientBase: gradientBase{id: id, stops: stops},
		x1:           x1, y1: y1, x2: x2, y2: y2,
	}
}

// Def returns the SVG <linearGradient> element for use in a <defs> block.
func (g *LinearGradient) Def() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<linearGradient id="%s" x1="%s" y1="%s" x2="%s" y2="%s">`,
		g.id, g.x1, g.y1, g.x2, g.y2,
	))
	for _, s := range g.stops {
		sb.WriteString(fmt.Sprintf(
			`<stop offset="%s" stop-color="%s" stop-opacity="%.3f"/>`,
			s.Offset, s.Color, s.Opacity,
		))
	}
	sb.WriteString(`</linearGradient>`)
	return sb.String()
}

// Ref returns the fill/stroke reference string: "url(#id)".
func (g *LinearGradient) Ref() string {
	return fmt.Sprintf("url(#%s)", g.id)
}

// ── Radial Gradient ───────────────────────────────────────────────────────────

// RadialGradient defines a radial gradient emanating from a focal point.
type RadialGradient struct {
	gradientBase
	cx, cy, r  string // center and radius
	fx, fy     string // focal point (optional, defaults to cx/cy)
}

// Radial creates a RadialGradient.
// id must be unique within the document.
//
//	draw.Radial("orb-grad", "40%", "35%", "60%",
//	    draw.S("0%", "#ffffff", 0.95),
//	    draw.S("100%", "#7c3aed", 0.0),
//	)
func Radial(id, cx, cy, r string, stops ...Stop) *RadialGradient {
	return &RadialGradient{
		gradientBase: gradientBase{id: id, stops: stops},
		cx:           cx, cy: cy, r: r,
	}
}

// WithFocalPoint sets the gradient's focal point (fx, fy).
func (g *RadialGradient) WithFocalPoint(fx, fy string) *RadialGradient {
	g.fx = fx
	g.fy = fy
	return g
}

// Def returns the SVG <radialGradient> element for use in a <defs> block.
func (g *RadialGradient) Def() string {
	var sb strings.Builder
	attrs := fmt.Sprintf(`id="%s" cx="%s" cy="%s" r="%s"`, g.id, g.cx, g.cy, g.r)
	if g.fx != "" {
		attrs += fmt.Sprintf(` fx="%s" fy="%s"`, g.fx, g.fy)
	}
	sb.WriteString(fmt.Sprintf(`<radialGradient %s>`, attrs))
	for _, s := range g.stops {
		sb.WriteString(fmt.Sprintf(
			`<stop offset="%s" stop-color="%s" stop-opacity="%.3f"/>`,
			s.Offset, s.Color, s.Opacity,
		))
	}
	sb.WriteString(`</radialGradient>`)
	return sb.String()
}

// Ref returns the fill/stroke reference string: "url(#id)".
func (g *RadialGradient) Ref() string {
	return fmt.Sprintf("url(#%s)", g.id)
}
