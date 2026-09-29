package finite_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite"
	"github.com/dominionthedev/finite/draw"
)

func TestNewVisual(t *testing.T) {
	v := finite.New(400, 300).WithBackground("#111")
	v.Layer("main").Add(draw.NewCircle(200, 150, 40, draw.Fill("#a78bfa")))
	svg, err := v.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<circle") {
		t.Fatal("missing circle")
	}
	b, err := v.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 || !strings.Contains(string(b), "</svg>") {
		t.Fatal("Bytes failed")
	}
}
