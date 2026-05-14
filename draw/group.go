// Package draw provides SVG primitive shapes, text elements, and styling options.
package draw

import (
	"fmt"
	"strings"
)

// Group represents an SVG <g> element. It is used to group multiple Renderables
// and apply shared transformations, IDs, or opacity to the entire group.
type Group struct {
	id        string
	transform string
	opacity   float64
	children  []interface{ Render() (string, error) }
}

// GroupOption defines a functional option for configuring a Group.
type GroupOption func(*Group)

// GroupID sets the SVG ID for the group element.
func GroupID(id string) GroupOption {
	return func(g *Group) { g.id = id }
}

// GroupTransform sets the SVG transform attribute for the group.
func GroupTransform(t string) GroupOption {
	return func(g *Group) { g.transform = t }
}

// GroupOpacity sets the overall opacity (0.0 to 1.0) for the group.
func GroupOpacity(o float64) GroupOption {
	return func(g *Group) { g.opacity = o }
}

// NewGroup initializes a new empty Group.
func NewGroup(opts ...GroupOption) *Group {
	g := &Group{opacity: 1.0}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Add appends a Renderable element to the group.
func (g *Group) Add(r interface{ Render() (string, error) }) {
	g.children = append(g.children, r)
}

// Render generates the SVG XML for the group and all its children.
func (g *Group) Render() (string, error) {
	var sb strings.Builder

	sb.WriteString("<g")
	if g.id != "" {
		sb.WriteString(fmt.Sprintf(` id="%s"`, g.id))
	}
	if g.transform != "" {
		sb.WriteString(fmt.Sprintf(` transform="%s"`, g.transform))
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
