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
// and manages document-wide properties like dimensions, background, and shared definitions.
type Document struct {
	Width      int    // Width of the SVG in pixels
	Height     int    // Height of the SVG in pixels
	Background string // Optional background color (hex code, e.g., "#ffffff")

	renderables []Renderable
	defs        []string
}

// NewDocument initializes a new SVG canvas with the specified width and height.
func NewDocument(width, height int) *Document {
	return &Document{
		Width:  width,
		Height: height,
	}
}

// WithBackground sets a solid background color for the entire document.
// It returns the Document pointer to enable fluent chaining.
func (d *Document) WithBackground(color string) *Document {
	d.Background = color
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

// Render generates the complete SVG XML string for the entire document.
func (d *Document) Render() (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"`+
			` width="%d" height="%d" viewBox="0 0 %d %d">`,
		d.Width, d.Height, d.Width, d.Height,
	))
	sb.WriteByte('\n')

	// Shared defs block (document-level)
	if len(d.defs) > 0 {
		sb.WriteString("<defs>\n")
		for _, def := range d.defs {
			sb.WriteString(def)
			sb.WriteByte('\n')
		}
		sb.WriteString("</defs>\n")
	}

	// Background rectangle
	if d.Background != "" {
		sb.WriteString(fmt.Sprintf(
			`<rect width="%d" height="%d" fill="%s"/>`,
			d.Width, d.Height, d.Background,
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
