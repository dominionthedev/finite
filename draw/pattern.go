package draw

import "fmt"

// Pattern is an SVG <pattern> paint server. Use it as a fill via FillPattern
// (it implements the same Ref/Def contract as gradients).
type Pattern struct {
	id      string
	width   float64
	height  float64
	content string
	units   string // patternUnits: "userSpaceOnUse" (default) or "objectBoundingBox"
}

// NewPattern creates a pattern tile of the given size with raw SVG content.
func NewPattern(id string, width, height float64, content string) *Pattern {
	return &Pattern{
		id:      id,
		width:   width,
		height:  height,
		content: content,
		units:   "userSpaceOnUse",
	}
}

// PatternDots creates a simple dotted pattern (circle at each tile origin).
func PatternDots(id string, spacing, radius float64, color string) *Pattern {
	if color == "" {
		color = "#000"
	}
	return &Pattern{
		id:      id,
		width:   spacing,
		height:  spacing,
		content: fmt.Sprintf(`<circle cx="%.4f" cy="%.4f" r="%.4f" fill="%s"/>`, spacing/2, spacing/2, radius, color),
		units:   "userSpaceOnUse",
	}
}

// PatternStripes creates diagonal-ish horizontal stripe tiles.
func PatternStripes(id string, spacing, thickness float64, color string) *Pattern {
	if color == "" {
		color = "#000"
	}
	return &Pattern{
		id:      id,
		width:   spacing,
		height:  spacing,
		content: fmt.Sprintf(`<rect x="0" y="0" width="%.4f" height="%.4f" fill="%s"/>`, spacing, thickness, color),
		units:   "userSpaceOnUse",
	}
}

// ObjectBoundingBox switches the pattern to objectBoundingBox units.
func (p *Pattern) ObjectBoundingBox() *Pattern {
	p.units = "objectBoundingBox"
	return p
}

// Def returns the <pattern> element for a <defs> block.
func (p *Pattern) Def() string {
	units := ""
	if p.units != "" && p.units != "userSpaceOnUse" {
		units = fmt.Sprintf(` patternUnits="%s"`, p.units)
	}
	return fmt.Sprintf(
		`<pattern id="%s" width="%.4f" height="%.4f"%s>%s</pattern>`,
		p.id, p.width, p.height, units, p.content,
	)
}

// Ref returns the fill/stroke reference, e.g. url(#my-pattern).
func (p *Pattern) Ref() string {
	return fmt.Sprintf("url(#%s)", p.id)
}
