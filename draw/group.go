// Package draw provides SVG primitive shapes, text elements, and styling options.
package draw

import (
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/geom"
)

// Group represents an SVG <g> element. It is used to group multiple Renderables
// and apply shared transformations, IDs, or opacity to the entire group.
type Group struct {
	id           string
	transform    geom.Matrix
	rawTransform string // set when a raw SVG transform string is provided
	opacity      float64
	children     []interface{ Render() (string, error) }
}

// GroupOption defines a functional option for configuring a Group.
type GroupOption func(*Group)

// GroupID sets the SVG ID for the group element.
func GroupID(id string) GroupOption {
	return func(g *Group) { g.id = id }
}

// GroupTransform sets a full transform matrix on the group.
func GroupTransform(m geom.Matrix) GroupOption {
	return func(g *Group) { g.transform = m }
}

// GroupTransformString sets a raw SVG transform string on the group.
// Prefer GroupTransform(geom.Matrix) when constructing transforms in code.
func GroupTransformString(t string) GroupOption {
	return func(g *Group) {
		// Store as a pure translation of the string by using a matrix that
		// String() cannot simplify; we keep the raw string path via a
		// side channel by overwriting after matrix assignment.
		// Simpler approach: parse is not implemented, so we keep the raw
		// string in a dedicated field path by using TransformString-style.
		// For now, store identity and rely on a raw override.
		g.rawTransform = t
	}
}

// GroupOpacity sets the overall opacity (0.0 to 1.0) for the group.
func GroupOpacity(o float64) GroupOption {
	return func(g *Group) { g.opacity = o }
}

// NewGroup initializes a new empty Group.
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

// Add appends a Renderable element to the group.
func (g *Group) Add(r interface{ Render() (string, error) }) *Group {
	g.children = append(g.children, r)
	return g
}

// SetTransform replaces the group's transform.
func (g *Group) SetTransform(m geom.Matrix) *Group {
	g.transform = m
	return g
}

// Translate applies an additional translation to the group.
func (g *Group) Translate(tx, ty float64) *Group {
	g.transform = g.transform.Mul(geom.Translate(tx, ty))
	return g
}

// Rotate applies an additional rotation (degrees) around the origin.
func (g *Group) Rotate(degrees float64) *Group {
	g.transform = g.transform.Mul(geom.Rotate(degrees))
	return g
}

// RotateAround applies an additional rotation around (cx, cy).
func (g *Group) RotateAround(degrees, cx, cy float64) *Group {
	g.transform = g.transform.Mul(geom.RotateAround(degrees, cx, cy))
	return g
}

// Scale applies an additional non-uniform scale.
func (g *Group) Scale(sx, sy float64) *Group {
	g.transform = g.transform.Mul(geom.Scale(sx, sy))
	return g
}

// ScaleUniform applies an additional uniform scale.
func (g *Group) ScaleUniform(s float64) *Group {
	g.transform = g.transform.Mul(geom.ScaleUniform(s))
	return g
}

// Render generates the SVG XML for the group and all its children.
func (g *Group) Render() (string, error) {
	var sb strings.Builder

	sb.WriteString("<g")
	if g.id != "" {
		sb.WriteString(fmt.Sprintf(` id="%s"`, g.id))
	}
	tf := g.rawTransform
	if tf == "" {
		tf = g.transform.String()
	}
	if tf != "" {
		sb.WriteString(fmt.Sprintf(` transform="%s"`, tf))
	}
	if g.opacity != 1.0 && g.opacity > 0 {
		sb.WriteString(fmt.Sprintf(` opacity="%.4f"`, g.opacity))
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

// RawElement allows embedding raw SVG strings directly into a document.
// This is useful for importing legacy SVG fragments or manual optimizations.
type RawElement struct {
	Content string // Raw SVG XML content
}

// Render returns the raw content string as-is.
func (r *RawElement) Render() (string, error) {
	return r.Content, nil
}
