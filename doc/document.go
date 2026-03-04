package doc

import (
	"fmt"
	"os"
	"strings"

	"github.com/leraniode/alolyte/instance"
)

// Document is the SVG canvas.
// It holds all widget instances and composes them into a final SVG on export.
type Document struct {
	Width      int
	Height     int
	Background string // optional hex fill, e.g. "#0f0f1a". Empty = transparent.

	instances []instance.Instance
	defs      []string // raw <defs> fragments (shared gradients, filters, etc.)
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

// Add places a widget instance onto the canvas.
// Instances are rendered in the order they are added (painter's algorithm).
func (d *Document) Add(inst instance.Instance) {
	d.instances = append(d.instances, inst)
}

// AddDef injects a raw SVG <defs> fragment shared across the whole document.
// Useful for gradients or filters referenced by multiple widgets.
func (d *Document) AddDef(def string) {
	d.defs = append(d.defs, def)
}

// Render composes all instances and returns the complete SVG string.
func (d *Document) Render() (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"`+
			` width="%d" height="%d" viewBox="0 0 %d %d">`,
		d.Width, d.Height, d.Width, d.Height,
	))
	sb.WriteByte('\n')

	// Shared defs block
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

	// Render instances in painter's order
	for idx, inst := range d.instances {
		rendered, err := inst.Render()
		if err != nil {
			return "", fmt.Errorf("doc.Render: instance[%d] (%s): %w", idx, inst.Widget.Name, err)
		}
		sb.WriteString(rendered)
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
