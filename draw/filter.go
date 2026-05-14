// Package draw provides SVG primitive shapes, text elements, and styling options.
package draw

import "fmt"

// Filter represents an SVG filter with a unique ID and one or more filter primitives.
// It produces a <filter> element within a <defs> block and can be referenced by other elements.
type Filter struct {
	id      string
	content string // inner filter primitives
}

// Def returns the complete <filter> XML element for use in a <defs> block.
func (f *Filter) Def() string {
	return fmt.Sprintf(`<filter id="%s">%s</filter>`, f.id, f.content)
}

// Ref returns the SVG filter attribute value, e.g., "url(#my-filter)".
func (f *Filter) Ref() string {
	return fmt.Sprintf("url(#%s)", f.id)
}

// Blur creates a Gaussian blur filter with the specified standard deviation.
func Blur(id string, stdDev float64) *Filter {
	return &Filter{
		id:      id,
		content: fmt.Sprintf(`<feGaussianBlur stdDeviation="%.4f"/>`, stdDev),
	}
}

// Glow creates a soft outer glow filter by layering multiple blurred copies of the source graphic.
func Glow(id string, radius float64) *Filter {
	return &Filter{
		id: id,
		content: fmt.Sprintf(`
<feGaussianBlur in="SourceGraphic" stdDeviation="%.4f" result="blur"/>
<feMerge>
  <feMergeNode in="blur"/>
  <feMergeNode in="blur"/>
  <feMergeNode in="SourceGraphic"/>
</feMerge>`, radius),
	}
}

// Noise creates a fractal noise texture using feTurbulence.
// blendMode controls how the noise is composited with the source graphic (e.g., "overlay", "multiply").
func Noise(id string, frequency float64, octaves int, blendMode string) *Filter {
	if blendMode == "" {
		blendMode = "overlay"
	}
	return &Filter{
		id: id,
		content: fmt.Sprintf(`
<feTurbulence type="fractalNoise" baseFrequency="%.4f" numOctaves="%d" stitchTiles="stitch" result="noise"/>
<feColorMatrix type="saturate" values="0" in="noise" result="grey"/>
<feBlend in="SourceGraphic" in2="grey" mode="%s"/>`, frequency, octaves, blendMode),
	}
}

// Shadow creates a simple drop shadow filter with the specified offset, blur, and color.
func Shadow(id string, dx, dy, blur float64, color string, opacity float64) *Filter {
	return &Filter{
		id: id,
		content: fmt.Sprintf(`
<feDropShadow dx="%.4f" dy="%.4f" stdDeviation="%.4f" flood-color="%s" flood-opacity="%.3f"/>`,
			dx, dy, blur, color, opacity),
	}
}
