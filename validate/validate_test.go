package validate_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/validate"
)

func TestValidDocument(t *testing.T) {
	d := doc.NewDocument(200, 200)
	d.Add(draw.NewCircle(100, 100, 40, draw.Fill("#a78bfa")))
	r, err := validate.Document(d)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Valid() {
		t.Fatalf("expected valid, got: %s", r)
	}
	if r.Empty() {
		// summary info should always be present for non-empty docs
		t.Fatal("expected at least summary info")
	}
	if len(r.Infos()) == 0 {
		t.Fatal("expected Infos() summary")
	}
}

func TestDuplicateIDMessage(t *testing.T) {
	// Layer name and child GroupID collide — the real bug from examples/basic
	d := doc.NewDocument(100, 100)
	g := draw.NewGroup(draw.GroupID("iso"))
	g.Add(draw.NewRect(0, 0, 10, 10, draw.Fill("#fff")))
	d.Layer("iso").Add(g)

	r, err := validate.Document(d)
	if err != nil {
		t.Fatal(err)
	}
	if r.Valid() {
		t.Fatal("expected invalid due to duplicate id")
	}
	errs := r.Errors()
	if len(errs) == 0 {
		t.Fatal("expected errors")
	}
	msg := errs[0].Message
	if !strings.Contains(msg, `duplicate id "iso"`) {
		t.Fatalf("message: %s", msg)
	}
	if !strings.Contains(msg, "layer names become group ids") {
		t.Fatalf("expected guidance about layers: %s", msg)
	}
}

func TestUnresolvedRef(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><rect fill="url(#missing)" width="10" height="10"/></svg>`
	r, err := validate.SVG(svg)
	if err != nil {
		t.Fatal(err)
	}
	if r.Valid() {
		t.Fatal("expected invalid")
	}
	found := false
	for _, e := range r.Errors() {
		if strings.Contains(e.Message, "missing") {
			found = true
		}
	}
	if !found {
		t.Fatalf("errors: %s", r)
	}
}

func TestUnusedPaintServerInfo(t *testing.T) {
	// Gradient defined but never referenced
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10">
<defs><linearGradient id="orphan"><stop offset="0%" stop-color="#fff"/></linearGradient></defs>
<rect width="10" height="10" fill="#000"/>
</svg>`
	r, err := validate.SVG(svg)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Valid() {
		t.Fatalf("should be valid: %s", r)
	}
	found := false
	for _, i := range r.Infos() {
		if strings.Contains(i.Message, "orphan") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unused paint server info, got: %s", r)
	}
}

func TestZeroDimensionWarning(t *testing.T) {
	d := doc.NewDocument(100, 100)
	d.Add(draw.NewRect(0, 0, 0, 50, draw.Fill("#f00")))
	r, err := validate.Document(d)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range r.Warnings() {
		if strings.Contains(w.Message, "zero width") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected zero width warning: %s", r)
	}
}

func TestEmptySVG(t *testing.T) {
	r, err := validate.SVG("")
	if err != nil {
		t.Fatal(err)
	}
	if r.Valid() {
		t.Fatal("empty should be invalid")
	}
}
