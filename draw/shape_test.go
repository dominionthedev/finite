package draw_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/draw"
)

func TestPathSetD(t *testing.T) {
	p := draw.NewPath()
	p.SetD("M100,200 L300,200 Z")
	svg, err := p.Render()
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(svg, `d="M100,200 L300,200 Z"`) {
		t.Errorf("Expected path data not found in SVG: %s", svg)
	}
}

func TestPathFluentAndSetD(t *testing.T) {
	p := draw.NewPath().M(10, 10).L(20, 20)
	p.SetD("M0,0 L10,10 Z")
	svg, err := p.Render()
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(svg, `d="M0,0 L10,10 Z"`) {
		t.Errorf("Expected path data not found in SVG: %s", svg)
	}
	if strings.Contains(svg, "M10.0000,10.0000") {
		t.Errorf("Old path data still present in SVG: %s", svg)
	}
}
