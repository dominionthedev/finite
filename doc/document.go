// Package doc provides the core SVG document model and rendering engine for Finite.
// It defines the Document canvas, the Renderable interface, and a lightweight
// layer system for ordered composition.
package doc

import (
	"fmt"
	"os"
	"strings"

	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/geom"
)

// Renderable is anything that can produce an SVG fragment.
// Shapes, text, groups, layers, and widget instances all implement this.
type Renderable interface {
	Render() (string, error)
}

// Layer is a named, ordered group of Renderables.
// Layers are the primary way to structure a Document: they preserve draw order,
// carry a stable id, and can hold a shared transform or opacity.
//
// A Layer is itself a Renderable, so layers can be nested if needed.
type Layer struct {
	name      string
	group     *draw.Group
}

// NewLayer creates a named layer.
func NewLayer(name string) *Layer {
	return &Layer{
		name:  name,
		group: draw.NewGroup(draw.GroupID(name)),
	}
}

// Name returns the layer name (also used as the SVG id).
func (l *Layer) Name() string {
	return l.name
}

// Add appends children to the layer and returns the layer for chaining.
func (l *Layer) Add(rs ...Renderable) *Layer {
	for _, r := range rs {
		l.group.Add(r)
	}
	return l
}

// Group returns the underlying draw.Group for advanced configuration
// (transforms, opacity, filters).
func (l *Layer) Group() *draw.Group {
	return l.group
}

// Translate applies a translation to the entire layer.
func (l *Layer) Translate(tx, ty float64) *Layer {
	l.group.Translate(tx, ty)
	return l
}

// SetTransform replaces the layer's transform.
func (l *Layer) SetTransform(m geom.Matrix) *Layer {
	l.group.SetTransform(m)
	return l
}

// Render emits the layer as an SVG group.
func (l *Layer) Render() (string, error) {
	return l.group.Render()
}

// Document is the SVG canvas. It holds document-level properties, an ordered
// list of root Renderables, and an optional set of named layers.
type Document struct {
	Width  int
	Height int

	ViewBoxMinX   float64
	ViewBoxMinY   float64
	ViewBoxWidth  float64
	ViewBoxHeight float64

	PreserveAspectRatio string
	Background          string

	// root holds elements added via Add (flat list, drawn in order).
	root []Renderable

	// layers holds named layers added via Layer / AddLayer.
	// Layers are drawn after root elements, in the order they were created.
	layers []*Layer

	defs []string
}

// NewDocument creates a canvas with the given viewport size.
// The viewBox is initialized to match (0 0 width height).
func NewDocument(width, height int) *Document {
	return &Document{
		Width:               width,
		Height:              height,
		ViewBoxWidth:        float64(width),
		ViewBoxHeight:       float64(height),
		PreserveAspectRatio: "xMidYMid meet",
	}
}

// WithBackground sets a solid background color.
func (d *Document) WithBackground(color string) *Document {
	d.Background = color
	return d
}

// WithViewBox sets the user coordinate system.
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

// Add appends one or more Renderables to the root drawing list.
// Prefer Layer when you want named, structured composition.
func (d *Document) Add(rs ...Renderable) {
	d.root = append(d.root, rs...)
}

// AddDef registers a raw SVG defs fragment (gradients, filters, clipPaths, etc.).
func (d *Document) AddDef(def string) {
	d.defs = append(d.defs, def)
}

// Layer returns an existing layer by name, or creates it if it does not exist.
// Layers are drawn in creation order after any root elements added via Add.
func (d *Document) Layer(name string) *Layer {
	for _, l := range d.layers {
		if l.name == name {
			return l
		}
	}
	l := NewLayer(name)
	d.layers = append(d.layers, l)
	return l
}

// AddLayer appends a pre-built layer.
func (d *Document) AddLayer(l *Layer) {
	d.layers = append(d.layers, l)
}

// Layers returns the document's layers in draw order.
func (d *Document) Layers() []*Layer {
	out := make([]*Layer, len(d.layers))
	copy(out, d.layers)
	return out
}

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

// Render generates the complete SVG document.
func (d *Document) Render() (string, error) {
	var sb strings.Builder

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

	if len(d.defs) > 0 {
		sb.WriteString("<defs>\n")
		for _, def := range d.defs {
			sb.WriteString(def)
			sb.WriteByte('\n')
		}
		sb.WriteString("</defs>\n")
	}

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

	// Root elements first
	for idx, r := range d.root {
		svg, err := r.Render()
		if err != nil {
			return "", fmt.Errorf("doc.Render: root[%d]: %w", idx, err)
		}
		sb.WriteString(svg)
		sb.WriteByte('\n')
	}

	// Then named layers in creation order
	for idx, l := range d.layers {
		svg, err := l.Render()
		if err != nil {
			return "", fmt.Errorf("doc.Render: layer[%d] %q: %w", idx, l.name, err)
		}
		sb.WriteString(svg)
		sb.WriteByte('\n')
	}

	sb.WriteString("</svg>")
	return sb.String(), nil
}

// Export renders the document and writes it to filePath.
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
