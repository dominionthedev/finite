// Package doc provides the core SVG document model and rendering engine for Finite.
// It defines the Document canvas and the Renderable interface for composable elements.
package doc

import (
	"fmt"
	"os"
	"strings"
)

// Renderable is the core interface for anything that can be drawn onto a Finite document.
// Standard shapes, text, widget instances, and groups all implement this interface.
type Renderable interface {
	// Render returns a string containing a valid SVG fragment.
	Render() (string, error)
}

// Document represents an SVG canvas. It acts as a container for Renderable elements
// and manages document-wide properties like dimensions, coordinate system, background,
// and shared definitions.
type Document struct {
	// Width and Height control the SVG viewport size (presentation attributes).
	// They may be 0 when the document is intended to be purely viewBox-driven.
	Width  int
	Height int

	// ViewBox defines the user coordinate system: minX, minY, width, height.
	// When zero values are present, Render falls back to "0 0 Width Height".
	ViewBoxMinX   float64
	ViewBoxMinY   float64
	ViewBoxWidth  float64
	ViewBoxHeight float64

	// PreserveAspectRatio controls how the viewBox is mapped to the viewport.
	// Common values: "xMidYMid meet" (default), "none", "xMinYMin slice", etc.
	PreserveAspectRatio string

	Background string // Optional background color (hex, e.g. "#0a0a14")

	renderables []Renderable
	defs        []string
}

// NewDocument initializes a new SVG canvas with the specified viewport width and height.
// The viewBox is initialized to match (0 0 width height).
func NewDocument(width, height int) *Document {
	return &Document{
		Width:             width,
		Height:            height,
		ViewBoxWidth:      float64(width),
		ViewBoxHeight:     float64(height),
		PreserveAspectRatio: "xMidYMid meet",
	}
}

// WithBackground sets a solid background color for the entire document.
// It returns the Document pointer to enable fluent chaining.
func (d *Document) WithBackground(color string) *Document {
	d.Background = color
	return d
}

// WithViewBox sets the user coordinate system.
// This is the primary way to control the coordinate space independently of
// the presentation width/height (important for responsive and nested SVGs).
func (d *Document) WithViewBox(minX, minY, width, height float64) *Document {
	d.ViewBoxMinX = minX
	d.ViewBoxMinY = minY
	d.ViewBoxWidth = width
	d.ViewBoxHeight = height
	return d
}

// WithPreserveAspectRatio sets the preserveAspectRatio attribute.
func (d *Document) WithPreserveAspectRatio(value string) *Document {
	d.PreserveAspectRatio = value
	return d
}

// Add appends a Renderable element to the document's drawing list.
// Elements are rendered in the order they are added (Painter's Algorithm).
func (d *Document) Add(r Renderable) {
	d.renderables = append(d.renderables, r)
}

// AddDef registers a raw SVG <defs> fragment to be included in the document's header.
// This is typically used for shared gradients, filters, or clip paths.
func (d *Document) AddDef(def string) {
	d.defs = append(d.defs, def)
}

// viewBoxString returns the viewBox attribute value.
func (d *Document) viewBoxString() string {
	w, h := d.ViewBoxWidth, d.ViewBoxHeight
	if w == 0 {
		w = float64(d.Width)
	}
	if h == 0 {
		h = float64(d.Height)
	}
	return fmt.Sprintf("%.4g %.4g %.4g %.4g", d.ViewBoxMinX, d.ViewBoxMinY, w, h)
}

// Render generates the complete SVG XML string for the entire document.
func (d *Document) Render() (string, error) {
	var sb strings.Builder

	// Opening <svg> tag
	sb.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"`)

	if d.Width > 0 {
		sb.WriteString(fmt.Sprintf(` width="%d"`, d.Width))
	}
	if d.Height > 0 {
		sb.WriteString(fmt.Sprintf(` height="%d"`, d.Height))
	}

	sb.WriteString(fmt.Sprintf(` viewBox="%s"`, d.viewBoxString()))

	if d.PreserveAspectRatio != "" && d.PreserveAspectRatio != "xMidYMid meet" {
		sb.WriteString(fmt.Sprintf(` preserveAspectRatio="%s"`, d.PreserveAspectRatio))
	}
	sb.WriteString(">\n")

	// Shared defs block (document-level)
	if len(d.defs) > 0 {
		sb.WriteString("<defs>\n")
		for _, def := range d.defs {
			sb.WriteString(def)
			sb.WriteByte('\n')
		}
		sb.WriteString("</defs>\n")
	}

	// Background rectangle (covers the full viewBox)
	if d.Background != "" {
		w, h := d.ViewBoxWidth, d.ViewBoxHeight
		if w == 0 {
			w = float64(d.Width)
		}
		if h == 0 {
			h = float64(d.Height)
		}
		sb.WriteString(fmt.Sprintf(
			`<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="%s"/>`,
			d.ViewBoxMinX, d.ViewBoxMinY, w, h, d.Background,
		))
		sb.WriteByte('\n')
	}

	// Render all elements in painter's order
	for idx, r := range d.renderables {
		svg, err := r.Render()
		if err != nil {
			return "", fmt.Errorf("doc.Render: element[%d]: %w", idx, err)
		}
		sb.WriteString(svg)
		sb.WriteByte('\n')
	}

	sb.WriteString("</svg>")
	return sb.String(), nil
}

// Export renders the document and writes the SVG to filePath.
func (d *Document) Export(filePath string) error {
	svg, err := d.Render()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, []byte(svg), 0644); err != nil {
		return fmt.Errorf("doc.Export: write failed for %q: %w", filePath, err)
	}
	return nil
}
