package instance

import (
	"fmt"
	"strings"

	"github.com/leraniode/alolyte/widget"
)

// Params holds key-value overrides passed into a widget's SVG template.
// Keys map directly to template variables, e.g. "PrimaryColor" → "#FF00FF".
type Params map[string]string

// Instance is a widget placed on the document canvas.
// It carries its own position, transform, and parameter overrides —
// the same widget can be placed many times with different settings.
type Instance struct {
	Widget   widget.Widget
	X, Y     float64
	Scale    float64
	Rotation float64
	Params   Params
}

// Option is a functional option for configuring an Instance at creation time.
type Option func(*Instance)

// At creates an Instance placing widget w at (x, y) on the canvas.
// Use the Option helpers to set scale, rotation, and params.
//
//	instance.At(orb, 200, 200,
//	    instance.WithScale(1.5),
//	    instance.WithParams(instance.Params{"PrimaryColor": "#FF6B6B"}),
//	)
func At(w widget.Widget, x, y float64, opts ...Option) Instance {
	inst := Instance{
		Widget:   w,
		X:        x,
		Y:        y,
		Scale:    1.0,
		Rotation: 0,
		Params:   Params{},
	}
	for _, opt := range opts {
		opt(&inst)
	}
	return inst
}

// WithScale sets the uniform scale multiplier for the instance.
func WithScale(s float64) Option {
	return func(i *Instance) { i.Scale = s }
}

// WithRotation sets the rotation in degrees around the instance's origin.
func WithRotation(deg float64) Option {
	return func(i *Instance) { i.Rotation = deg }
}

// WithParams merges the given Params into the instance,
// overriding any keys already set.
func WithParams(p Params) Option {
	return func(i *Instance) {
		for k, v := range p {
			i.Params[k] = v
		}
	}
}

// Render returns the final SVG snippet for this instance:
// the widget's rendered SVG content wrapped in a <g> transform group.
func (inst Instance) Render() (string, error) {
	svg, err := inst.Widget.Render(map[string]string(inst.Params))
	if err != nil {
		return "", fmt.Errorf("instance.Render (%s): %w", inst.Widget.Name, err)
	}

	inner := stripSVGWrapper(svg)
	transform := buildTransform(inst.X, inst.Y, inst.Scale, inst.Rotation)

	return fmt.Sprintf(`<g transform="%s">%s</g>`, transform, inner), nil
}

// buildTransform produces an SVG transform attribute value.
func buildTransform(x, y, scale, rotation float64) string {
	parts := []string{
		fmt.Sprintf("translate(%.4f,%.4f)", x, y),
	}
	if scale != 1.0 {
		parts = append(parts, fmt.Sprintf("scale(%.4f)", scale))
	}
	if rotation != 0 {
		parts = append(parts, fmt.Sprintf("rotate(%.4f)", rotation))
	}
	return strings.Join(parts, " ")
}

// stripSVGWrapper removes the outer <svg ...>...</svg> tags from widget content
// so it embeds cleanly inside a parent SVG document.
func stripSVGWrapper(content string) string {
	content = strings.TrimSpace(content)

	start := strings.Index(content, ">")
	if start == -1 {
		return content
	}

	end := strings.LastIndex(content, "</svg>")
	if end == -1 {
		return strings.TrimSpace(content[start+1:])
	}

	return strings.TrimSpace(content[start+1 : end])
}
