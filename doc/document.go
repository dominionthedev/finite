package doc

import (
	"fmt"
	"os"
	"strings"
)

// Renderable is anything that can produce an SVG fragment.
// instance.Instance, draw shapes, text elements — all implement this.
// Widgets are just one kind of Renderable.
type Renderable interface {
	Render() (string, error)
}

// Document is the SVG canvas.
// It composes Renderables — shapes, text, gradients, widget instances —
// into a single exported SVG file.
type Document struct {
	Width      int
	Height     int
	Background string // optional hex fill, e.g. "#0f0f1a". Empty = transparent.

	renderables []Renderable
	defs        []string // shared <defs> fragments (gradients, filters, etc.)
}

// NewDocument creates a blank canvas with the given pixel dimensions.
func NewDocument(width, height int) *Document {
	return &Document{
		Width:  width,
		Height: height,
	}
}

// WithBackground sets a solid background color for the canvas.
// Returns the Document so calls can be chained.
func (d *Document) WithBackground(color string) *Document {
	d.Background = color
	return d
}

// Add places any Renderable onto the canvas.
// Renderables are drawn in the order they are added (painter's algorithm).
//
// Accepted types:
//   - instance.Instance       — a widget SVG placed with transforms
//   - draw.Circle, draw.Rect  — native shapes
//   - draw.Text               — styled text
//   - Any type with Render() (string, error)
func (d *Document) Add(r Renderable) {
	d.renderables = append(d.renderables, r)
}

// AddDef injects a raw SVG <defs> fragment shared across the whole document.
// Use this for gradients or filters referenced by multiple elements.
func (d *Document) AddDef(def string) {
	d.defs = append(d.defs, def)
}

// Render composes all renderables and returns the complete SVG string.
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
