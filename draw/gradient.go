// Package draw provides SVG primitive shapes, text elements, and styling options.
package draw

import (
	"fmt"
	"strings"
)

// Stop represents a single color stop within an SVG gradient.
type Stop struct {
	Offset  string  // Offset as percentage (e.g., "0%", "100%")
	Color   string  // Hex color code (e.g., "#a78bfa")
	Opacity float64 // Opacity value (0.0 to 1.0)
}

// S is a shorthand constructor for a gradient Stop.
func S(offset, color string, opacity float64) Stop {
	return Stop{Offset: offset, Color: color, Opacity: opacity}
}

// GradientOption is a functional option for configuring gradients.
type GradientOption func(*gradientBase)

type gradientBase struct {
	id    string
	stops []Stop
}

// LinearGradient defines a linear color transition along a vector (x1,y1) to (x2,y2).
type LinearGradient struct {
	gradientBase
	x1, y1, x2, y2 string // Positions as percentages or absolute units
}

// Linear creates a new LinearGradient with the specified ID and coordinates.
func Linear(id, x1, y1, x2, y2 string, stops ...Stop) *LinearGradient {
	return &LinearGradient{
		gradientBase: gradientBase{id: id, stops: stops},
		x1:           x1, y1: y1, x2: x2, y2: y2,
	}
}

// Def returns the SVG <linearGradient> XML definition.
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

// Ref returns the SVG fill/stroke reference string, e.g., "url(#my-linear-grad)".
func (g *LinearGradient) Ref() string {
	return fmt.Sprintf("url(#%s)", g.id)
}

// RadialGradient defines a radial color transition emanating from a focal point.
type RadialGradient struct {
	gradientBase
	cx, cy, r string // Center and radius as percentages or absolute units
	fx, fy    string // Focal point (optional)
}

// Radial creates a new RadialGradient with the specified ID, center, and radius.
func Radial(id, cx, cy, r string, stops ...Stop) *RadialGradient {
	return &RadialGradient{
		gradientBase: gradientBase{id: id, stops: stops},
		cx:           cx, cy: cy, r: r,
	}
}

// WithFocalPoint sets the focal point (fx, fy) for the radial gradient.
func (g *RadialGradient) WithFocalPoint(fx, fy string) *RadialGradient {
	g.fx = fx
	g.fy = fy
	return g
}

// Def returns the SVG <radialGradient> XML definition.
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

// Ref returns the SVG fill/stroke reference string, e.g., "url(#my-radial-grad)".
func (g *RadialGradient) Ref() string {
	return fmt.Sprintf("url(#%s)", g.id)
}
