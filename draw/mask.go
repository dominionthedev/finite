package draw

import "fmt"

// Mask is an SVG <mask> definition. Masks use luminance (and optionally alpha)
// of their content to control visibility of the referencing element.
//
// White = fully visible, black = fully hidden, gray = partial.
type Mask struct {
	id      string
	content string
	units   string // "userSpaceOnUse" (default) or "objectBoundingBox"
}

// NewMask creates a mask from raw SVG content.
func NewMask(id, content string) *Mask {
	return &Mask{id: id, content: content, units: "userSpaceOnUse"}
}

// MaskRect creates a rectangular mask filled with the given color.
// Use white for fully visible, black for hidden, or any gray for partial.
func MaskRect(id string, x, y, w, h float64, fill string) *Mask {
	if fill == "" {
		fill = "#fff"
	}
	return &Mask{
		id:      id,
		content: fmt.Sprintf(`<rect x="%.4f" y="%.4f" width="%.4f" height="%.4f" fill="%s"/>`, x, y, w, h, fill),
		units:   "userSpaceOnUse",
	}
}

// MaskGradient creates a mask whose content is a full-size rect filled with a gradient.
// The gradient definition is inlined inside the mask. Useful for soft falloff / vignette.
func MaskGradient(id string, x, y, w, h float64, g fillProvider) *Mask {
	return &Mask{
		id: id,
		content: g.Def() + fmt.Sprintf(
			`<rect x="%.4f" y="%.4f" width="%.4f" height="%.4f" fill="%s"/>`,
			x, y, w, h, g.Ref(),
		),
		units: "userSpaceOnUse",
	}
}

// ObjectBoundingBox switches the mask to objectBoundingBox units.
func (m *Mask) ObjectBoundingBox() *Mask {
	m.units = "objectBoundingBox"
	return m
}

// Def returns the <mask> element for a <defs> block.
func (m *Mask) Def() string {
	units := ""
	if m.units != "" && m.units != "userSpaceOnUse" {
		units = fmt.Sprintf(` maskUnits="%s"`, m.units)
	}
	return fmt.Sprintf(`<mask id="%s"%s>%s</mask>`, m.id, units, m.content)
}

// Ref returns the mask attribute value, e.g. url(#my-mask).
func (m *Mask) Ref() string {
	return fmt.Sprintf("url(#%s)", m.id)
}
