// Package instance provides tools to place Widgets on a document canvas with geometric transforms.
package instance

import (
	"fmt"

	"github.com/dominionthedev/finite/geom"
	"github.com/dominionthedev/finite/widget"
)

// Instance represents a widget placed on the canvas with a full affine transform.
//
// Transforms are composed so that local operations (scale, rotation) are applied
// first (around the widget's local origin), then the instance is translated into
// place. This matches the common mental model: "place the widget, then scale/rotate it".
type Instance struct {
	w         widget.Widget
	transform geom.Matrix
}

// Option defines a functional option for configuring a widget Instance.
type Option func(*Instance)

// At creates a new Instance of the specified widget at (x, y).
// Additional options (scale, rotation, etc.) are applied in local space.
func At(w widget.Widget, x, y float64, opts ...Option) *Instance {
	// Collect local transforms from options first, then wrap with translation.
	local := geom.Identity()
	// Temporary instance used only to gather option effects on a local matrix.
	tmp := &Instance{w: w, transform: local}
	for _, opt := range opts {
		opt(tmp)
	}
	// Final transform: translate first in matrix multiplication order so that
	// local transforms run first on points, then translation.
	// point' = Translate * Local * point
	final := geom.Translate(x, y).Mul(tmp.transform)
	return &Instance{w: w, transform: final}
}

// WithTransform replaces any local transform with the given matrix.
// Note: when used with At, this becomes the local part before translation.
func WithTransform(m geom.Matrix) Option {
	return func(i *Instance) { i.transform = m }
}

// WithScale applies a uniform scale in local space.
func WithScale(s float64) Option {
	return func(i *Instance) {
		i.transform = i.transform.Mul(geom.ScaleUniform(s))
	}
}

// WithScaleXY applies a non-uniform scale in local space.
func WithScaleXY(sx, sy float64) Option {
	return func(i *Instance) {
		i.transform = i.transform.Mul(geom.Scale(sx, sy))
	}
}

// WithRotation applies a rotation (degrees) in local space around the origin.
func WithRotation(deg float64) Option {
	return func(i *Instance) {
		i.transform = i.transform.Mul(geom.Rotate(deg))
	}
}

// WithRotationAround applies a rotation (degrees) around (cx, cy) in local space.
func WithRotationAround(deg, cx, cy float64) Option {
	return func(i *Instance) {
		i.transform = i.transform.Mul(geom.RotateAround(deg, cx, cy))
	}
}

// Render generates the SVG XML for the instance, applying the transform.
func (inst *Instance) Render() (string, error) {
	svg, err := inst.w.Render()
	if err != nil {
		return "", fmt.Errorf("instance.Render (%s): %w", inst.w.Name(), err)
	}
	tf := inst.transform.String()
	if tf == "" {
		return svg, nil
	}
	return fmt.Sprintf(`<g transform="%s">%s</g>`, tf, svg), nil
}

// Transform returns the current transform matrix.
func (inst *Instance) Transform() geom.Matrix {
	return inst.transform
}
