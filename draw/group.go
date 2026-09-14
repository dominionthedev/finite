// Package draw provides SVG primitive shapes, text elements, and styling options.
package draw

import (
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/geom"
)

// Group represents an SVG <g> element. It is the primary composition primitive
// in Finite: a named, transformable container for other Renderables.
//
// Groups nest cleanly. Transforms, opacity, and filters applied to a group
// affect the entire subtree.
type Group struct {
	id           string
	transform    geom.Matrix
	rawTransform string // used when a raw SVG transform string is provided (e.g. from decode)
	opacity      float64
	filter       *Filter
	children     []interface{ Render() (string, error) }
}

// GroupOption configures a Group.
type GroupOption func(*Group)

// GroupID sets the SVG id attribute.
func GroupID(id string) GroupOption {
	return func(g *Group) { g.id = id }
}

// GroupTransform sets the group's transform matrix.
func GroupTransform(m geom.Matrix) GroupOption {
	return func(g *Group) {
		g.transform = m
		g.rawTransform = ""
	}
}

// GroupTransformString sets a raw SVG transform string.
// Prefer GroupTransform(geom.Matrix) when building transforms in code.
func GroupTransformString(t string) GroupOption {
	return func(g *Group) {
		g.rawTransform = t
		g.transform = geom.Identity()
	}
}

// GroupOpacity sets the group's opacity (0–1).
func GroupOpacity(o float64) GroupOption {
	return func(g *Group) { g.opacity = o }
}

// GroupFilter attaches an SVG filter to the entire group.
func GroupFilter(f *Filter) GroupOption {
	return func(g *Group) { g.filter = f }
}

// NewGroup creates an empty Group.
func NewGroup(opts ...GroupOption) *Group {
	g := &Group{
		opacity:   1.0,
		transform: geom.Identity(),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Add appends one or more children and returns the group for chaining.
func (g *Group) Add(rs ...interface{ Render() (string, error) }) *Group {
	g.children = append(g.children, rs...)
	return g
}

// Children returns a copy of the group's child list.
func (g *Group) Children() []interface{ Render() (string, error) } {
	out := make([]interface{ Render() (string, error) }, len(g.children))
	copy(out, g.children)
	return out
}

// Len returns the number of direct children.
func (g *Group) Len() int {
	return len(g.children)
}

// ID returns the group's SVG id, if any.
func (g *Group) ID() string {
	return g.id
}

// Transform returns the group's current transform matrix.
// If a raw transform string was set, this returns Identity; use the rendered
// attribute for the raw value.
func (g *Group) Transform() geom.Matrix {
	return g.transform
}

// SetTransform replaces the group's transform matrix.
func (g *Group) SetTransform(m geom.Matrix) *Group {
	g.transform = m
	g.rawTransform = ""
	return g
}

// Translate applies an additional translation.
func (g *Group) Translate(tx, ty float64) *Group {
	g.transform = g.transform.Mul(geom.Translate(tx, ty))
	g.rawTransform = ""
	return g
}

// Rotate applies an additional rotation in degrees around the origin.
func (g *Group) Rotate(degrees float64) *Group {
	g.transform = g.transform.Mul(geom.Rotate(degrees))
	g.rawTransform = ""
	return g
}

// RotateAround applies an additional rotation around (cx, cy).
func (g *Group) RotateAround(degrees, cx, cy float64) *Group {
	g.transform = g.transform.Mul(geom.RotateAround(degrees, cx, cy))
	g.rawTransform = ""
	return g
}

// Scale applies an additional non-uniform scale.
func (g *Group) Scale(sx, sy float64) *Group {
	g.transform = g.transform.Mul(geom.Scale(sx, sy))
	g.rawTransform = ""
	return g
}

// ScaleUniform applies an additional uniform scale.
func (g *Group) ScaleUniform(s float64) *Group {
	g.transform = g.transform.Mul(geom.ScaleUniform(s))
	g.rawTransform = ""
	return g
}

// transformAttr returns the transform attribute value to emit, if any.
func (g *Group) transformAttr() string {
	if g.rawTransform != "" {
		return g.rawTransform
	}
	return g.transform.String()
}

// Render generates the SVG for the group and its children.
func (g *Group) Render() (string, error) {
	var sb strings.Builder

	// Inline filter def if present
	if g.filter != nil {
		sb.WriteString("<defs>")
		sb.WriteString(g.filter.Def())
		sb.WriteString("</defs>")
	}

	sb.WriteString("<g")
	if g.id != "" {
		sb.WriteString(fmt.Sprintf(` id="%s"`, g.id))
	}
	if tf := g.transformAttr(); tf != "" {
		sb.WriteString(fmt.Sprintf(` transform="%s"`, tf))
	}
	if g.opacity != 1.0 && g.opacity > 0 {
		sb.WriteString(fmt.Sprintf(` opacity="%.4f"`, g.opacity))
	}
	if g.filter != nil {
		sb.WriteString(fmt.Sprintf(` filter="%s"`, g.filter.Ref()))
	}
	sb.WriteString(">")

	for i, child := range g.children {
		svg, err := child.Render()
		if err != nil {
			return "", fmt.Errorf("draw.Group: child[%d]: %w", i, err)
		}
		sb.WriteString(svg)
	}

	sb.WriteString("</g>")
	return sb.String(), nil
}

// RawElement embeds a raw SVG fragment. Useful for decode results or
// hand-tuned markup.
type RawElement struct {
	Content string
}

// Render returns the raw content unchanged.
func (r *RawElement) Render() (string, error) {
	return r.Content, nil
}
