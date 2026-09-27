package draw_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/draw"
)

func TestPolygon(t *testing.T) {
	p := draw.NewPolygon([]draw.Point{
		draw.Pt(0, 0), draw.Pt(40, 0), draw.Pt(20, 30),
	}, draw.Fill("#a78bfa"))
	svg, err := p.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<polygon") {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, "0.0000,0.0000") {
		t.Fatal(svg)
	}
}

func TestPolylineStroke(t *testing.T) {
	p := draw.NewPolyline([]draw.Point{
		draw.Pt(0, 0), draw.Pt(10, 20), draw.Pt(30, 15),
	}, draw.Stroke("#fff", 2), draw.StrokeDash("4 2"), draw.StrokeLinecap("round"))
	svg, err := p.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<polyline") {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, `stroke-dasharray="4 2"`) {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, `stroke-linecap="round"`) {
		t.Fatal(svg)
	}
	if !strings.Contains(svg, `fill="none"`) {
		t.Fatal("polyline should default fill none")
	}
}

func TestNoFill(t *testing.T) {
	c := draw.NewCircle(0, 0, 10, draw.NoFill(), draw.Stroke("#000", 1))
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `fill="none"`) {
		t.Fatal(svg)
	}
}

func TestPathArc(t *testing.T) {
	p := draw.NewPath(draw.NoFill(), draw.Stroke("#fff", 1.5)).
		M(10, 50).A(40, 40, 0, 1, 1, 90, 50)
	svg, err := p.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "A40.0000,40.0000") {
		t.Fatal(svg)
	}
}
