package widget

import (
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/geom"
)

// Composite builds a widget from other widgets (and optional raw fragments).
// It is the standard way to assemble larger components from smaller ones.
//
// The composite owns a local transform and opacity applied to the whole tree.
type Composite struct {
	name      string
	transform geom.Matrix
	opacity   float64
	children  []Widget
}

// CompositeOption configures a Composite.
type CompositeOption func(*Composite)

// WithOpacity sets opacity for the entire composite (0–1).
func WithOpacity(o float64) CompositeOption {
	return func(c *Composite) { c.opacity = o }
}

// WithLocalTransform sets a local transform applied around the composite origin.
func WithLocalTransform(m geom.Matrix) CompositeOption {
	return func(c *Composite) { c.transform = m }
}

// NewComposite creates a composite widget with the given name and children.
func NewComposite(name string, children []Widget, opts ...CompositeOption) *Composite {
	c := &Composite{
		name:      name,
		transform: geom.Identity(),
		opacity:   1.0,
		children:  children,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Name returns the composite name.
func (c *Composite) Name() string {
	return c.name
}

// Add appends children and returns the composite for chaining.
func (c *Composite) Add(ws ...Widget) *Composite {
	c.children = append(c.children, ws...)
	return c
}

// Render emits a <g> containing all children, with optional transform/opacity.
func (c *Composite) Render() (string, error) {
	var sb strings.Builder
	sb.WriteString("<g")
	if tf := c.transform.String(); tf != "" {
		sb.WriteString(fmt.Sprintf(` transform="%s"`, tf))
	}
	if c.opacity != 1.0 && c.opacity > 0 {
		sb.WriteString(fmt.Sprintf(` opacity="%.4f"`, c.opacity))
	}
	sb.WriteString(">")

	for i, child := range c.children {
		if child == nil {
			continue
		}
		svg, err := child.Render()
		if err != nil {
			return "", fmt.Errorf("widget.Composite %q: child[%d] %s: %w", c.name, i, child.Name(), err)
		}
		sb.WriteString(svg)
	}
	sb.WriteString("</g>")
	return sb.String(), nil
}
