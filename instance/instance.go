// Package instance places any Widget onto a document canvas with
// optional transforms — translate, scale, rotation.
//
// The same widget can be placed multiple times at different positions,
// scales, and orientations. Each placement is an independent Instance.
package instance

import (
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/widget"
)

// Instance is a widget placed on the document canvas with a position
// and optional transform.
type Instance struct {
	w        widget.Widget
	x, y     float64
	scale    float64
	rotation float64 // degrees
}

// Option configures an Instance.
type Option func(*Instance)

// At creates an Instance placing w at (x, y).
//
//	instance.At(ring, 300, 300, instance.WithScale(1.5), instance.WithRotation(30))
func At(w widget.Widget, x, y float64, opts ...Option) *Instance {
	inst := &Instance{w: w, x: x, y: y, scale: 1.0}
	for _, opt := range opts {
		opt(inst)
	}
	return inst
}

// WithScale sets a uniform scale multiplier.
func WithScale(s float64) Option { return func(i *Instance) { i.scale = s } }

// WithRotation sets rotation in degrees.
func WithRotation(deg float64) Option { return func(i *Instance) { i.rotation = deg } }

// Render implements doc.Renderable.
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
