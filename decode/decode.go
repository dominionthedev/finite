// Package decode reads standard SVG files and converts them into
// Finite's native types — shapes, text, groups, gradients.
//
// This is Finite's import layer. Any existing SVG can be loaded,
// decoded into composable elements, and then extended or re-exported
// through the rest of the Finite pipeline.
//
//	canvas, err := decode.File("my-logo.svg")
//	canvas.Add(draw.NewText("v2", 400, 50, draw.TextFill("#fff"), draw.Centered()))
//	canvas.Export("output.svg")
package decode

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
)

// File reads an SVG file from disk and returns a populated doc.Document.
func File(path string) (*doc.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("decode.File: cannot read %q: %w", path, err)
	}
	return Bytes(data)
}

// Bytes parses raw SVG bytes into a doc.Document.
func Bytes(data []byte) (*doc.Document, error) {
	root, err := parseXML(data)
	if err != nil {
		return nil, fmt.Errorf("decode.Bytes: %w", err)
	}

	w := attrFloat(root.Attrs, "width", 0)
	h := attrFloat(root.Attrs, "height", 0)
	if vb, ok := root.Attrs["viewBox"]; ok {
		parts := parseViewBox(vb)
		if w == 0 { w = parts[2] }
		if h == 0 { h = parts[3] }
	}

	d := doc.NewDocument(int(w), int(h))

	for _, child := range root.Children {
		switch child.Tag {
		case "defs":
			// Re-inject defs as raw fragments
			for _, def := range child.Children {
				d.AddDef(renderRaw(def))
			}
		default:
			r := toRenderable(child)
			if r != nil {
				d.Add(r)
			}
		}
	}

	return d, nil
}

// ── XML node ─────────────────────────────────────────────────────────────────

type xmlNode struct {
	Tag      string
	Attrs    map[string]string
	Children []*xmlNode
	Text     string
}

func parseXML(data []byte) (*xmlNode, error) {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	var stack []*xmlNode
	var root *xmlNode

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			node := &xmlNode{
				Tag:   t.Name.Local,
				Attrs: make(map[string]string),
			}
			for _, a := range t.Attr {
				node.Attrs[a.Name.Local] = a.Value
			}
			// Merge inline style into attrs (lower priority than direct attrs)
			if style, ok := node.Attrs["style"]; ok {
				for k, v := range parseStyle(style) {
					if _, exists := node.Attrs[k]; !exists {
						node.Attrs[k] = v
					}
				}
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)
			if root == nil && t.Name.Local == "svg" {
				root = node
			}

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += strings.TrimSpace(string(t))
			}
		}
	}

	if root == nil {
		return nil, fmt.Errorf("no <svg> root element found")
	}
	return root, nil
}

// ── SVG element → draw.Renderable ────────────────────────────────────────────

func toRenderable(n *xmlNode) doc.Renderable {
	opts := shapeOpts(n)
	switch n.Tag {
	case "circle":
		cx := attrFloat(n.Attrs, "cx", 0)
		cy := attrFloat(n.Attrs, "cy", 0)
		r  := attrFloat(n.Attrs, "r", 0)
		return draw.NewCircle(cx, cy, r, opts...)

	case "ellipse":
		cx := attrFloat(n.Attrs, "cx", 0)
		cy := attrFloat(n.Attrs, "cy", 0)
		rx := attrFloat(n.Attrs, "rx", 0)
		ry := attrFloat(n.Attrs, "ry", 0)
		return draw.NewEllipse(cx, cy, rx, ry, opts...)

	case "rect":
		x := attrFloat(n.Attrs, "x", 0)
		y := attrFloat(n.Attrs, "y", 0)
		w := attrFloat(n.Attrs, "width", 0)
		h := attrFloat(n.Attrs, "height", 0)
		rx := attrFloat(n.Attrs, "rx", 0)
		if rx > 0 {
			return draw.NewRoundedRect(x, y, w, h, rx, opts...)
		}
		return draw.NewRect(x, y, w, h, opts...)

	case "line":
		x1 := attrFloat(n.Attrs, "x1", 0)
		y1 := attrFloat(n.Attrs, "y1", 0)
		x2 := attrFloat(n.Attrs, "x2", 0)
		y2 := attrFloat(n.Attrs, "y2", 0)
		return draw.NewLine(x1, y1, x2, y2, opts...)

	case "path":
		p := draw.NewPath(opts...)
		if d, ok := n.Attrs["d"]; ok {
			p.SetD(d)
		}
		return p

	case "text":
		x := attrFloat(n.Attrs, "x", 0)
		y := attrFloat(n.Attrs, "y", 0)
		content := n.Text
		return draw.NewText(content, x, y, textOpts(n)...)

	case "g":
		g := draw.NewGroup(groupOpts(n)...)
		for _, child := range n.Children {
			r := toRenderable(child)
			if r != nil {
				g.Add(r)
			}
		}
		return g

	default:
		// Unknown element: passthrough as raw SVG
		raw := renderRaw(n)
		if raw == "" {
			return nil
		}
		return &draw.RawElement{Content: raw}
	}
}

// ── Option builders ───────────────────────────────────────────────────────────

func shapeOpts(n *xmlNode) []draw.ShapeOption {
	var opts []draw.ShapeOption

	if fill, ok := n.Attrs["fill"]; ok && fill != "none" && fill != "" {
		opts = append(opts, draw.Fill(fill))
	}
	if stroke, ok := n.Attrs["stroke"]; ok && stroke != "none" {
		sw := attrFloat(n.Attrs, "stroke-width", 1.0)
		opts = append(opts, draw.Stroke(stroke, sw))
	}
	if op, ok := n.Attrs["opacity"]; ok {
		if v, err := strconv.ParseFloat(op, 64); err == nil {
			opts = append(opts, draw.Opacity(v))
		}
	}
	if tfm, ok := n.Attrs["transform"]; ok {
		opts = append(opts, draw.Transform(tfm))
	}
	return opts
}

func textOpts(n *xmlNode) []draw.TextOption {
	var opts []draw.TextOption
	if fill, ok := n.Attrs["fill"]; ok {
		opts = append(opts, draw.TextFill(fill))
	}
	if fs, ok := n.Attrs["font-size"]; ok {
		if v := parseFloat(fs); v > 0 {
			opts = append(opts, draw.FontSize(v))
		}
	}
	if fw, ok := n.Attrs["font-weight"]; ok {
		opts = append(opts, draw.FontWeight(fw))
	}
	if ff, ok := n.Attrs["font-family"]; ok {
		opts = append(opts, draw.FontFamily(ff))
	}
	if ta, ok := n.Attrs["text-anchor"]; ok && ta == "middle" {
		opts = append(opts, draw.Centered())
	}
	return opts
}

func groupOpts(n *xmlNode) []draw.GroupOption {
	var opts []draw.GroupOption
	if id, ok := n.Attrs["id"]; ok {
		opts = append(opts, draw.GroupID(id))
	}
	if tfm, ok := n.Attrs["transform"]; ok {
		opts = append(opts, draw.GroupTransform(tfm))
	}
	if op, ok := n.Attrs["opacity"]; ok {
		if v, err := strconv.ParseFloat(op, 64); err == nil {
			opts = append(opts, draw.GroupOpacity(v))
		}
	}
	return opts
}

// ── Transform math ────────────────────────────────────────────────────────────

// Transform holds a decomposed SVG transform.
type Transform struct {
	TranslateX, TranslateY float64
	ScaleX, ScaleY         float64
	RotateDeg              float64
}

// ParseTransform parses an SVG transform attribute string.
// Handles translate(), scale(), rotate(), and matrix().
func ParseTransform(s string) Transform {
	t := Transform{ScaleX: 1, ScaleY: 1}
	for _, part := range splitFunctions(s) {
		name, args := parseFn(part)
		switch name {
		case "translate":
			if len(args) >= 1 { t.TranslateX = args[0] }
			if len(args) >= 2 { t.TranslateY = args[1] }
		case "scale":
			if len(args) >= 1 { t.ScaleX = args[0]; t.ScaleY = args[0] }
			if len(args) >= 2 { t.ScaleY = args[1] }
		case "rotate":
			if len(args) >= 1 { t.RotateDeg = args[0] }
		case "matrix":
			if len(args) == 6 {
				a, b, c, d, e, f := args[0], args[1], args[2], args[3], args[4], args[5]
				t.TranslateX = e
				t.TranslateY = f
				t.ScaleX = math.Sqrt(a*a + b*b)
				t.ScaleY = math.Sqrt(c*c + d*d)
				t.RotateDeg = math.Atan2(b, a) * (180 / math.Pi)
			}
		}
	}
	return t
}

func splitFunctions(s string) []string {
	var parts []string
	depth, start := 0, 0
	for i, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				parts = append(parts, s[start:i+1])
				start = i + 1
			}
		}
	}
	return parts
}

func parseFn(s string) (string, []float64) {
	lp := strings.Index(s, "(")
	rp := strings.LastIndex(s, ")")
	if lp == -1 || rp == -1 {
		return "", nil
	}
	name := strings.TrimSpace(s[:lp])
	var args []float64
	for _, tok := range strings.FieldsFunc(s[lp+1:rp], func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	}) {
		if v, err := strconv.ParseFloat(tok, 64); err == nil {
			args = append(args, v)
		}
	}
	return name, args
}

// ── Colour extraction ─────────────────────────────────────────────────────────

// ExtractColors walks a decoded document and returns every unique
// hex color found in fill, stroke, and stop-color attributes.
func ExtractColors(d *doc.Document) []string {
	// We re-render and scan the SVG string for hex colors.
	svg, err := d.Render()
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var colors []string
	i := 0
	for i < len(svg) {
		if svg[i] == '#' && i+7 <= len(svg) {
			candidate := svg[i : i+7]
			valid := true
			for _, ch := range candidate[1:] {
				if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
					valid = false
					break
				}
			}
			if valid && !seen[candidate] {
				seen[candidate] = true
				colors = append(colors, candidate)
			}
		}
		i++
	}
	return colors
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func renderRaw(n *xmlNode) string {
	var b strings.Builder
	b.WriteString("<")
	b.WriteString(n.Tag)
	for k, v := range n.Attrs {
		b.WriteString(fmt.Sprintf(` %s="%s"`, k, v))
	}
	if len(n.Children) == 0 && n.Text == "" {
		b.WriteString("/>")
		return b.String()
	}
	b.WriteString(">")
	b.WriteString(n.Text)
	for _, child := range n.Children {
		b.WriteString(renderRaw(child))
	}
	b.WriteString(fmt.Sprintf("</%s>", n.Tag))
	return b.String()
}

func parseViewBox(s string) [4]float64 {
	var out [4]float64
	parts := strings.Fields(strings.ReplaceAll(s, ",", " "))
	for i, p := range parts {
		if i >= 4 { break }
		out[i], _ = strconv.ParseFloat(p, 64)
	}
	return out
}

func attrFloat(attrs map[string]string, key string, def float64) float64 {
	v, ok := attrs[key]
	if !ok { return def }
	return parseFloat(v)
}

func parseFloat(s string) float64 {
	s = strings.TrimRight(strings.TrimSpace(s), "pxtemr%")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil { return 0 }
	return v
}

func parseStyle(s string) map[string]string {
	m := make(map[string]string)
	for _, decl := range strings.Split(s, ";") {
		parts := strings.SplitN(strings.TrimSpace(decl), ":", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m
}
