package draw

import "fmt"

// ClipPath is an SVG <clipPath> definition. Attach it to shapes or groups
// with WithClip / GroupClip to restrict painting to the clip region.
//
// Luminance is not used for clipPath — the geometry itself defines the region
// (filled areas of the clip content are visible).
type ClipPath struct {
	id      string
	content string
	units   string // "userSpaceOnUse" (default) or "objectBoundingBox"
}

// NewClipPath creates a clipPath from raw SVG content (e.g. a <path> or <circle>).
func NewClipPath(id, content string) *ClipPath {
	return &ClipPath{id: id, content: content, units: "userSpaceOnUse"}
}

// ClipCircle creates a circular clip region.
func ClipCircle(id string, cx, cy, r float64) *ClipPath {
	return &ClipPath{
		id:      id,
		content: fmt.Sprintf(`<circle cx="%.4f" cy="%.4f" r="%.4f"/>`, cx, cy, r),
		units:   "userSpaceOnUse",
	}
}

// ClipRect creates a rectangular clip region.
func ClipRect(id string, x, y, w, h float64) *ClipPath {
	return &ClipPath{
		id:      id,
		content: fmt.Sprintf(`<rect x="%.4f" y="%.4f" width="%.4f" height="%.4f"/>`, x, y, w, h),
		units:   "userSpaceOnUse",
	}
}

// ClipPathData creates a clip region from SVG path data (the d attribute).
func ClipPathData(id, d string) *ClipPath {
	return &ClipPath{
		id:      id,
		content: fmt.Sprintf(`<path d="%s"/>`, d),
		units:   "userSpaceOnUse",
	}
}

// ObjectBoundingBox switches the clipPath to objectBoundingBox units
// (coordinates are relative to the referencing element's bounding box).
func (c *ClipPath) ObjectBoundingBox() *ClipPath {
	c.units = "objectBoundingBox"
	return c
}

// Def returns the <clipPath> element for a <defs> block.
func (c *ClipPath) Def() string {
	units := ""
	if c.units != "" && c.units != "userSpaceOnUse" {
		units = fmt.Sprintf(` clipPathUnits="%s"`, c.units)
	}
	return fmt.Sprintf(`<clipPath id="%s"%s>%s</clipPath>`, c.id, units, c.content)
}

// Ref returns the clip-path attribute value, e.g. url(#my-clip).
func (c *ClipPath) Ref() string {
	return fmt.Sprintf("url(#%s)", c.id)
}
