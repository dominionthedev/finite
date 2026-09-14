package widget_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/geom"
	"github.com/dominionthedev/finite/widget"
)

func TestScopeNamespaces(t *testing.T) {
	s1 := widget.NewScope("card")
	s2 := widget.NewScope("card")
	if s1.Prefix() == s2.Prefix() {
		t.Fatal("expected unique prefixes for separate scopes")
	}
	id := s1.ID("grad")
	if !strings.HasPrefix(id, s1.Prefix()+"-") {
		t.Errorf("ID not namespaced: %q", id)
	}
	if strings.Contains(id, " ") {
		t.Errorf("ID should be sanitized: %q", id)
	}
}

func TestScopeIDStable(t *testing.T) {
	s := widget.NewScopeID("icon", "primary")
	if s.ID("fill") != "icon-primary-fill" {
		t.Errorf("got %q", s.ID("fill"))
	}
}

func TestFuncWidget(t *testing.T) {
	w := widget.Func{
		N: "dot",
		F: func() (string, error) {
			return `<circle cx="0" cy="0" r="4" fill="#fff"/>`, nil
		},
	}
	if w.Name() != "dot" {
		t.Fatal(w.Name())
	}
	svg, err := w.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<circle") {
		t.Fatal(svg)
	}
}

func TestFuncNil(t *testing.T) {
	w := widget.Func{N: "broken"}
	if _, err := w.Render(); err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestStatic(t *testing.T) {
	w := widget.Static{N: "mark", Content: `<rect width="10" height="10"/>`}
	svg, err := w.Render()
	if err != nil {
		t.Fatal(err)
	}
	if svg != `<rect width="10" height="10"/>` {
		t.Fatal(svg)
	}
}

func TestComposite(t *testing.T) {
	a := widget.Static{N: "a", Content: `<circle r="5"/>`}
	b := widget.Static{N: "b", Content: `<rect width="1" height="1"/>`}
	c := widget.NewComposite("pair", []widget.Widget{a, b},
		widget.WithOpacity(0.5),
		widget.WithLocalTransform(geom.Translate(10, 20)),
	)
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `opacity="0.5000"`) {
		t.Error("opacity missing")
	}
	if !strings.Contains(svg, "translate") {
		t.Error("transform missing")
	}
	if !strings.Contains(svg, "<circle") || !strings.Contains(svg, "<rect") {
		t.Error("children missing")
	}
}

func TestCompositeAdd(t *testing.T) {
	c := widget.NewComposite("stack", nil)
	c.Add(widget.Static{Content: `<path d="M0,0"/>`})
	if c.Name() != "stack" {
		t.Fatal(c.Name())
	}
	svg, err := c.Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<path") {
		t.Fatal(svg)
	}
}

func TestSanitize(t *testing.T) {
	s := widget.NewScopeID("My Widget!", "x y")
	id := s.ID("a/b")
	for _, bad := range []string{" ", "!", "/"} {
		if strings.Contains(id, bad) {
			t.Errorf("unsanitized %q in %q", bad, id)
		}
	}
}
