package draw_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/draw"
)

func TestShapeTitleDesc(t *testing.T) {
	c := draw.NewCircle(10, 10, 5,
		draw.Fill("#fff"),
		draw.Title("Dot"),
		draw.Desc("A small white circle"),
		draw.Role("img"),
	)
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<title>Dot</title>") {
		t.Error("title missing:", svg)
	}
	if !strings.Contains(svg, "<desc>A small white circle</desc>") {
		t.Error("desc missing")
	}
	if !strings.Contains(svg, `role="img"`) {
		t.Error("role missing")
	}
	// Should not be self-closing when title present
	if strings.Contains(svg, "<circle ") && strings.Contains(svg, "/>") && !strings.Contains(svg, "</circle>") {
		// might still have other self-closing; ensure </circle> exists
	}
	if !strings.Contains(svg, "</circle>") {
		t.Error("expected open/close circle for a11y children")
	}
}

func TestAriaLabel(t *testing.T) {
	r := draw.NewRect(0, 0, 10, 10, draw.Fill("#000"), draw.AriaLabel("block"))
	svg, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `aria-label="block"`) {
		t.Fatal(svg)
	}
}

func TestGroupA11y(t *testing.T) {
	g := draw.NewGroup(draw.GroupTitle("Cluster"), draw.GroupRole("group"))
	g.Add(draw.NewCircle(0, 0, 3, draw.Fill("#f00")))
	svg, err := g.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<title>Cluster</title>") {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, `role="group"`) {
		t.Fatal(svg)
	}
}

func TestXMLEscapeInTitle(t *testing.T) {
	c := draw.NewCircle(0, 0, 1, draw.Title(`A & B <C>`))
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "A &amp; B &lt;C&gt;") {
		t.Fatal(svg)
	}
}
