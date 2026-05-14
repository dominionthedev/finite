package draw

import (
	"fmt"
	"strings"
)

// ── Group ─────────────────────────────────────────────────────────────────────

// Group is a collection of Renderables wrapped in an SVG <g> element.
// Use it to apply a shared transform or opacity to multiple elements,
// or to logically organise layers.
type Group struct {
	id         string
	transform  string
	opacity    float64
	children   []interface{ Render() (string, error) }
}

// GroupOption configures a Group.
type GroupOption func(*Group)

// GroupID sets the SVG id attribute on the <g> element.
func GroupID(id string) GroupOption {
	return func(g *Group) { g.id = id }
}

// GroupTransform sets the transform attribute on the <g> element.
func GroupTransform(t string) GroupOption {
	return func(g *Group) { g.transform = t }
}

// GroupOpacity sets the opacity attribute on the <g> element.
func GroupOpacity(o float64) GroupOption {
	return func(g *Group) { g.opacity = o }
}

// NewGroup creates a new empty Group.
func NewGroup(opts ...GroupOption) *Group {
	g := &Group{opacity: 1.0}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Add appends a Renderable to the group.
func (g *Group) Add(r interface{ Render() (string, error) }) {
	g.children = append(g.children, r)
}

// Render implements doc.Renderable.
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

// ── RawElement ────────────────────────────────────────────────────────────────

// RawElement is an escape hatch that passes an SVG string through
// directly without any transformation or validation.
//
// Use it for:
//   - SVG fragments produced by the decode package
//   - Externally generated SVG you want to composite
//   - Manually written SVG for edge cases the draw API doesn't cover
type RawElement struct {
	Content string
}

// Render implements doc.Renderable. Returns Content as-is.
func (r *RawElement) Render() (string, error) {
	return r.Content, nil
}
