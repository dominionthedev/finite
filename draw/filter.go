package draw

import "fmt"

// Filter represents an SVG filter with a unique ID.
// It produces a <filter> def and a url(#id) reference.
type Filter struct {
	id      string
	content string // inner filter primitives
}

// Def returns the complete <filter> element for a <defs> block.
func (f *Filter) Def() string {
	return fmt.Sprintf(`<filter id="%s">%s</filter>`, f.id, f.content)
}

// Ref returns the filter attribute value: "url(#id)".
func (f *Filter) Ref() string {
	return fmt.Sprintf("url(#%s)", f.id)
}

// Blur creates a Gaussian blur filter.
//
//	draw.Blur("my-blur", 12)
func Blur(id string, stdDev float64) *Filter {
	return &Filter{
		id:      id,
		content: fmt.Sprintf(`<feGaussianBlur stdDeviation="%.4f"/>`, stdDev),
	}
}

// Glow creates a soft glow filter by layering a blurred copy beneath the source.
//
//	draw.Glow("my-glow", 8)
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

// Noise creates an feTurbulence grain/texture overlay filter.
// blendMode controls how the noise blends: "overlay", "multiply", "screen", etc.
//
//	draw.Noise("my-noise", 0.65, 3, "overlay")
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

// Shadow creates a drop shadow filter.
//
//	draw.Shadow("my-shadow", 4, 4, 6, "#000000", 0.5)
func Shadow(id string, dx, dy, blur float64, color string, opacity float64) *Filter {
	return &Filter{
		id: id,
		content: fmt.Sprintf(`
<feDropShadow dx="%.4f" dy="%.4f" stdDeviation="%.4f" flood-color="%s" flood-opacity="%.3f"/>`,
			dx, dy, blur, color, opacity),
	}
}
