package widget

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Widget is a reusable SVG design component loaded from a .svg file.
// The SVG content can use Go text/template syntax for parameterization.
//
// Example inside an SVG widget file:
//
//	<stop stop-color="{{if .PrimaryColor}}{{.PrimaryColor}}{{else}}#a78bfa{{end}}"/>
type Widget struct {
	Name    string
	Path    string
	Content string // raw SVG template source
}

// Load reads an SVG file from disk and returns a Widget.
func Load(path string) (Widget, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Widget{}, fmt.Errorf("widget.Load: cannot read %q: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".svg")

	return Widget{
		Name:    name,
		Path:    path,
		Content: string(data),
	}, nil
}

// Render executes the widget's SVG template with the given params,
// returning the rendered SVG string ready for embedding into a Document.
func (w Widget) Render(params map[string]string) (string, error) {
	tmpl, err := template.New(w.Name).Option("missingkey=zero").Parse(w.Content)
	if err != nil {
		return "", fmt.Errorf("widget.Render: bad template in %q: %w", w.Name, err)
	}

	// Convert map[string]string → map[string]any for template execution
	data := make(map[string]any, len(params))
	for k, v := range params {
		data[k] = v
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("widget.Render: execute failed for %q: %w", w.Name, err)
	}

	return buf.String(), nil
}
