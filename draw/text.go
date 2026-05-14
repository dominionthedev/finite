package draw

import (
	"fmt"
	"strings"
)

// TextAnchor controls horizontal text alignment.
type TextAnchor string

const (
	AnchorStart  TextAnchor = "start"
	AnchorMiddle TextAnchor = "middle"
	AnchorEnd    TextAnchor = "end"
)

// Text is a styled SVG <text> element.
// Supports solid fills, gradient fills, and glow filters.
type Text struct {
	content     string
	x, y        float64
	fontFamily  string
	fontSize    float64
	fontWeight  string
	letterSpace float64
	anchor      TextAnchor
	style       shapeStyle
}

// TextOption configures a Text element.
type TextOption func(*Text)

// NewText creates a text element at (x, y).
//
//	draw.NewText("Finite", 400, 300,
//	    draw.FontSize(64),
//	    draw.FontWeight("700"),
//	    draw.TextFill("#ffffff"),
//	    draw.Centered(),
//	)
func NewText(content string, x, y float64, opts ...TextOption) *Text {
	t := &Text{
		content:    content,
		x:          x,
		y:          y,
		fontFamily: `ui-monospace, "SF Mono", monospace`,
		fontSize:   48,
		fontWeight: "400",
		anchor:     AnchorStart,
	}
	t.style.opacity = 1.0
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// FontSize sets the font size in pixels.
func FontSize(size float64) TextOption {
	return func(t *Text) { t.fontSize = size }
}

// FontWeight sets the font weight ("400", "700", "900", etc.)
func FontWeight(w string) TextOption {
	return func(t *Text) { t.fontWeight = w }
}

// FontFamily sets the font-family string.
func FontFamily(f string) TextOption {
	return func(t *Text) { t.fontFamily = f }
}

// LetterSpacing sets the letter-spacing in pixels.
func LetterSpacing(s float64) TextOption {
	return func(t *Text) { t.letterSpace = s }
}

// TextFill sets a solid hex color fill for the text.
func TextFill(color string) TextOption {
	return func(t *Text) {
		t.style.fill = color
		t.style.fillDef = ""
	}
}

// TextFillGradient sets a gradient fill for the text.
func TextFillGradient(g fillProvider) TextOption {
	return func(t *Text) {
		t.style.fill = g.Ref()
		t.style.fillDef = g.Def()
	}
}

// TextFilter attaches a filter (e.g. Glow) to the text element.
func TextFilter(f *Filter) TextOption {
	return func(t *Text) { t.style.filter = f }
}

// Centered sets text-anchor and dominant-baseline to center the text at (x, y).
func Centered() TextOption {
	return func(t *Text) { t.anchor = AnchorMiddle }
}

// Anchor sets the text-anchor alignment.
func Anchor(a TextAnchor) TextOption {
	return func(t *Text) { t.anchor = a }
}

// TextOpacity sets the text opacity.
func TextOpacity(v float64) TextOption {
	return func(t *Text) { t.style.opacity = v }
}

// Render implements doc.Renderable.
func (t *Text) Render() (string, error) {
	var attrs strings.Builder

	attrs.WriteString(fmt.Sprintf(` x="%.4f" y="%.4f"`, t.x, t.y))
	family := strings.ReplaceAll(t.fontFamily, `"`, `&quot;`)
	attrs.WriteString(fmt.Sprintf(` font-family="%s"`, family))
	attrs.WriteString(fmt.Sprintf(` font-size="%.4f"`, t.fontSize))
	attrs.WriteString(fmt.Sprintf(` font-weight="%s"`, t.fontWeight))

	if t.anchor == AnchorMiddle {
		attrs.WriteString(` text-anchor="middle" dominant-baseline="middle"`)
	} else {
		attrs.WriteString(fmt.Sprintf(` text-anchor="%s"`, t.anchor))
	}

	if t.letterSpace != 0 {
		attrs.WriteString(fmt.Sprintf(` letter-spacing="%.4f"`, t.letterSpace))
	}

	attrs.WriteString(t.style.attrString())

	el := fmt.Sprintf(`<text%s>%s</text>`, attrs.String(), t.content)
	return wrap(t.style.buildDefs(), el), nil
}
