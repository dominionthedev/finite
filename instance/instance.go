// Package instance provides tools to place Widgets on a document canvas with geometric transforms.
package instance

import (
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/widget"
)

// Instance represents a widget placed at a specific location with optional scaling and rotation.
type Instance struct {
	w        widget.Widget
	x, y     float64
	scale    float64
	rotation float64 // rotation in degrees
}

// Option defines a functional option for configuring a widget Instance.
type Option func(*Instance)

// At creates a new Instance of the specified widget at the coordinates (x, y).
func At(w widget.Widget, x, y float64, opts ...Option) *Instance {
	inst := &Instance{w: w, x: x, y: y, scale: 1.0}
	for _, opt := range opts {
		opt(inst)
	}
	return inst
}

// WithScale returns an Option that sets the uniform scale for the instance.
func WithScale(s float64) Option { return func(i *Instance) { i.scale = s } }

// WithRotation returns an Option that sets the rotation (in degrees) for the instance.
func WithRotation(deg float64) Option { return func(i *Instance) { i.rotation = deg } }

// Render generates the SVG XML for the instance, applying the necessary transforms.
func (inst *Instance) Render() (string, error) {
	svg, err := inst.w.Render()
	if err != nil {
		return "", fmt.Errorf("instance.Render (%s): %w", inst.w.Name(), err)
	}
	tfm := buildTransform(inst.x, inst.y, inst.scale, inst.rotation)
	return fmt.Sprintf(`<g transform="%s">%s</g>`, tfm, svg), nil
}

func buildTransform(x, y, scale, rotation float64) string {
	parts := []string{fmt.Sprintf("translate(%.4f,%.4f)", x, y)}
	if scale != 1.0 {
		parts = append(parts, fmt.Sprintf("scale(%.4f)", scale))
	}
	if rotation != 0 {
		parts = append(parts, fmt.Sprintf("rotate(%.4f)", rotation))
	}
	return strings.Join(parts, " ")
}
