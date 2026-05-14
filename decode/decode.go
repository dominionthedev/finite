// Package decode provides functionality to parse existing SVG files and bytes into
// Finite's native document model. It supports common SVG elements and transforms.
package decode

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
)

// File reads an SVG file from the specified path and decodes it into a doc.Document.
func File(path string) (*doc.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("decode.File: cannot read %q: %w", path, err)
	}
	return Bytes(data)
}

// Bytes parses raw SVG XML data and returns a populated doc.Document.
func Bytes(data []byte) (*doc.Document, error) {
	root, err := parseXML(data)
	if err != nil {
		return nil, fmt.Errorf("decode.Bytes: %w", err)
	}

	w := attrFloat(root.Attrs, "width", 0)
	h := attrFloat(root.Attrs, "height", 0)
	if vb, ok := root.Attrs["viewBox"]; ok {
		parts := parseViewBox(vb)
		if w == 0 {
			w = parts[2]
		}
		if h == 0 {
			h = parts[3]
		}
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

// xmlNode is an internal representation of an XML element for easier processing.
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

// toRenderable converts an internal xmlNode into a Finite doc.Renderable.
func toRenderable(n *xmlNode) doc.Renderable {
	opts := shapeOpts(n)
	switch n.Tag {
	case "circle":
		cx := attrFloat(n.Attrs, "cx", 0)
		cy := attrFloat(n.Attrs, "cy", 0)
		r := attrFloat(n.Attrs, "r", 0)
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
		if v, err := strconv.ParseFloat(strings.TrimSuffix(fs, "px"), 64); err == nil {
			opts = append(opts, draw.FontSize(v))
		}
	}
	if fw, ok := n.Attrs["font-weight"]; ok {
		opts = append(opts, draw.FontWeight(fw))
	}
	if ta, ok := n.Attrs["text-anchor"]; ok {
		opts = append(opts, draw.Anchor(draw.TextAnchor(ta)))
	}
	return opts
}

func groupOpts(n *xmlNode) []draw.GroupOption {
	var opts []draw.GroupOption
	if id, ok := n.Attrs["id"]; ok {
		opts = append(opts, draw.GroupID(id))
	}
	if op, ok := n.Attrs["opacity"]; ok {
		if v, err := strconv.ParseFloat(op, 64); err == nil {
			opts = append(opts, draw.GroupOpacity(v))
		}
	}
	if tfm, ok := n.Attrs["transform"]; ok {
		opts = append(opts, draw.GroupTransform(tfm))
	}
	return opts
}

func splitFunctions(s string) []string {
	var fns []string
	var current strings.Builder
	depth := 0
	for _, r := range s {
		if r == '(' {
			depth++
		}
		if r == ')' {
			depth--
		}
		current.WriteRune(r)
		if depth == 0 && r == ' ' {
			if current.Len() > 0 {
				fns = append(fns, strings.TrimSpace(current.String()))
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		fns = append(fns, strings.TrimSpace(current.String()))
	}
	return fns
}

func parseFn(s string) (string, []float64) {
	open := strings.Index(s, "(")
	close := strings.Index(s, ")")
	if open == -1 || close == -1 {
		return "", nil
	}
	name := strings.TrimSpace(s[:open])
	argsStr := s[open+1 : close]
	argsStr = strings.ReplaceAll(argsStr, ",", " ")
	parts := strings.Fields(argsStr)
	args := make([]float64, len(parts))
	for i, p := range parts {
		args[i] = parseFloat(p)
	}
	return name, args
}

func renderRaw(n *xmlNode) string {
	var sb strings.Builder
	sb.WriteString("<" + n.Tag)
	for k, v := range n.Attrs {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, k, v))
	}
	if len(n.Children) == 0 && n.Text == "" {
		sb.WriteString("/>")
	} else {
		sb.WriteString(">")
		sb.WriteString(n.Text)
		for _, child := range n.Children {
			sb.WriteString(renderRaw(child))
		}
		sb.WriteString("</" + n.Tag + ">")
	}
	return sb.String()
}

func parseViewBox(s string) [4]float64 {
	parts := strings.Fields(strings.ReplaceAll(s, ",", " "))
	var res [4]float64
	for i := 0; i < 4 && i < len(parts); i++ {
		res[i] = parseFloat(parts[i])
	}
	return res
}

func attrFloat(attrs map[string]string, key string, def float64) float64 {
	if v, ok := attrs[key]; ok {
		return parseFloat(v)
	}
	return def
}

func parseFloat(s string) float64 {
	s = strings.TrimSuffix(s, "px")
	s = strings.TrimSuffix(s, "pt")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseStyle(s string) map[string]string {
	res := make(map[string]string)
	parts := strings.Split(s, ";")
	for _, p := range parts {
		kv := strings.Split(p, ":")
		if len(kv) == 2 {
			res[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return res
}
