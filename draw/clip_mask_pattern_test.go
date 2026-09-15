package draw_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/draw"
)

func TestClipCircleOnShape(t *testing.T) {
	c := draw.NewRect(0, 0, 100, 100,
		draw.Fill("#a78bfa"),
		draw.WithClip(draw.ClipCircle("c1", 50, 50, 40)),
	)
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `clip-path="url(#c1)"`) {
		t.Error("clip-path ref missing")
	}
	if !strings.Contains(svg, `<clipPath id="c1">`) {
		t.Error("clipPath def missing")
	}
	if !strings.Contains(svg, `<circle cx="50.0000"`) {
		t.Error("clip geometry missing")
	}
}

func TestClipPathData(t *testing.T) {
	c := draw.ClipPathData("p1", "M0,0 L100,0 L50,80 Z")
	if !strings.Contains(c.Def(), `d="M0,0 L100,0 L50,80 Z"`) {
		t.Fatal(c.Def())
	}
}

func TestMaskOnShape(t *testing.T) {
	r := draw.NewCircle(50, 50, 40,
		draw.Fill("#34d399"),
		draw.WithMask(draw.MaskRect("m1", 0, 0, 100, 100, "#fff")),
	)
	svg, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `mask="url(#m1)"`) {
		t.Error("mask ref missing")
	}
	if !strings.Contains(svg, `<mask id="m1">`) {
		t.Error("mask def missing")
	}
}

func TestMaskGradient(t *testing.T) {
	g := draw.Linear("fade", "0%", "0%", "100%", "0%",
		draw.S("0%", "#fff", 1),
		draw.S("100%", "#000", 1),
	)
	m := draw.MaskGradient("mg", 0, 0, 200, 100, g)
	def := m.Def()
	if !strings.Contains(def, "linearGradient") {
		t.Error("gradient should be inlined in mask")
	}
	if !strings.Contains(def, `fill="url(#fade)"`) {
		t.Error("rect should reference gradient")
	}
}

func TestPatternFill(t *testing.T) {
	r := draw.NewRect(0, 0, 200, 200,
		draw.FillPattern(draw.PatternDots("dots", 12, 2, "#6366f1")),
	)
	svg, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `fill="url(#dots)"`) {
		t.Error("pattern fill ref missing")
	}
	if !strings.Contains(svg, `<pattern id="dots"`) {
		t.Error("pattern def missing")
	}
}

func TestPatternStripes(t *testing.T) {
	p := draw.PatternStripes("stripes", 8, 3, "#ec4899")
	if !strings.Contains(p.Def(), `id="stripes"`) {
		t.Fatal(p.Def())
	}
}

func TestGroupClipAndMask(t *testing.T) {
	g := draw.NewGroup(
		draw.GroupID("clipped"),
		draw.GroupClip(draw.ClipRect("gc", 10, 10, 80, 80)),
		draw.GroupMask(draw.MaskRect("gm", 0, 0, 100, 100, "#ccc")),
	)
	g.Add(draw.NewCircle(50, 50, 40, draw.Fill("#f59e0b")))
	svg, err := g.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `clip-path="url(#gc)"`) {
		t.Error("group clip missing")
	}
	if !strings.Contains(svg, `mask="url(#gm)"`) {
		t.Error("group mask missing")
	}
	if !strings.Contains(svg, `<clipPath id="gc">`) || !strings.Contains(svg, `<mask id="gm">`) {
		t.Error("defs missing")
	}
}

func TestObjectBoundingBoxUnits(t *testing.T) {
	c := draw.ClipRect("obb", 0, 0, 1, 1).ObjectBoundingBox()
	if !strings.Contains(c.Def(), `clipPathUnits="objectBoundingBox"`) {
		t.Fatal(c.Def())
	}
	m := draw.MaskRect("mobb", 0, 0, 1, 1, "#fff").ObjectBoundingBox()
	if !strings.Contains(m.Def(), `maskUnits="objectBoundingBox"`) {
		t.Fatal(m.Def())
	}
	p := draw.NewPattern("pobb", 0.1, 0.1, `<circle r="0.05"/>`).ObjectBoundingBox()
	if !strings.Contains(p.Def(), `patternUnits="objectBoundingBox"`) {
		t.Fatal(p.Def())
	}
}
