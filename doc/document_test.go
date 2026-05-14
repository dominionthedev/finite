package doc_test

import (
	"strings"
	"testing"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/instance"
)

// minimal Widget implementation for testing
type testWidget struct{ name, color string }
func (w *testWidget) Name() string { return w.name }
func (w *testWidget) Render() (string, error) {
	return `<circle cx="50" cy="50" r="40" fill="` + w.color + `"/>`, nil
}

func TestNewDocument(t *testing.T) {
	d := doc.NewDocument(800, 600)
	if d.Width != 800 || d.Height != 600 {
		t.Errorf("expected 800×600, got %d×%d", d.Width, d.Height)
	}
}
func TestWithBackground(t *testing.T) {
	d := doc.NewDocument(800, 600).WithBackground("#0a0a14")
	if d.Background != "#0a0a14" {
		t.Errorf("got %q", d.Background)
	}
}
func TestRenderEmpty(t *testing.T) {
	svg, err := doc.NewDocument(400, 300).Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, `width="400"`) { t.Error("missing width") }
	if !strings.Contains(svg, `</svg>`)      { t.Error("missing closing tag") }
}
func TestRenderBackground(t *testing.T) {
	svg, err := doc.NewDocument(400, 300).WithBackground("#1a1a2e").Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, `fill="#1a1a2e"`) { t.Error("background fill missing") }
}
func TestRenderCircle(t *testing.T) {
	d := doc.NewDocument(400, 400)
	d.Add(draw.NewCircle(200, 200, 80, draw.Fill("#a78bfa")))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "<circle")      { t.Error("circle missing") }
	if !strings.Contains(svg, `fill="#a78bfa"`) { t.Error("fill missing") }
}
func TestRenderRect(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewRect(50, 50, 200, 100, draw.Fill("#ec4899")))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "<rect") { t.Error("rect missing") }
}
func TestRenderText(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewText("finite", 400, 300, draw.FontSize(64), draw.TextFill("#fff"), draw.Centered()))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "<text") { t.Error("text element missing") }
	if !strings.Contains(svg, "finite") { t.Error("text content missing") }
}
func TestRenderLinearGradient(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewRect(0, 0, 800, 600, draw.FillGradient(
		draw.Linear("g1", "0%", "0%", "100%", "0%",
			draw.S("0%", "#6366f1", 1.0), draw.S("100%", "#ec4899", 1.0)),
	)))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "linearGradient") { t.Error("linearGradient missing") }
}
func TestRenderRadialGradient(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewEllipse(400, 300, 200, 180, draw.FillGradient(
		draw.Radial("r1", "50%", "50%", "50%",
			draw.S("0%", "#fff", 1.0), draw.S("100%", "#7c3aed", 0.0)),
	)))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "radialGradient") { t.Error("radialGradient missing") }
}
func TestRenderFilter(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewCircle(400, 300, 150, draw.Fill("#a78bfa"), draw.WithFilter(draw.Blur("b1", 15))))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "feGaussianBlur") { t.Error("feGaussianBlur missing") }
}
func TestWidget(t *testing.T) {
	d := doc.NewDocument(400, 400)
	d.Add(&testWidget{"w", "#a78bfa"})
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "#a78bfa") { t.Error("widget color missing") }
}
func TestInstanceAt(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(instance.At(&testWidget{"w", "#ec4899"}, 100, 150))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "translate(100") { t.Error("translate missing") }
}
func TestInstanceScale(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(instance.At(&testWidget{"w", "#6366f1"}, 0, 0, instance.WithScale(2.0)))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "scale(") { t.Error("scale missing") }
}
func TestInstanceRotation(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(instance.At(&testWidget{"w", "#06b6d4"}, 200, 200, instance.WithRotation(45)))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "rotate(") { t.Error("rotation missing") }
}
func TestAddDef(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.AddDef(`<filter id="cf"><feGaussianBlur stdDeviation="5"/></filter>`)
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "cf") { t.Error("def id missing") }
}
func TestPaintersOrder(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewRect(0, 0, 800, 600, draw.Fill("#000")))
	d.Add(draw.NewCircle(400, 300, 100, draw.Fill("#fff")))
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if strings.Index(svg, "<rect") > strings.Index(svg, "<circle") {
		t.Error("painter order violated")
	}
}
func TestRawElement(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(&draw.RawElement{Content: `<polygon points="200,10 250,190 160,210" fill="#34d399"/>`})
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "<polygon") { t.Error("polygon missing") }
}
func TestGroup(t *testing.T) {
	d := doc.NewDocument(800, 600)
	g := draw.NewGroup(draw.GroupID("layer-1"), draw.GroupOpacity(0.5))
	g.Add(draw.NewCircle(100, 100, 50, draw.Fill("#f59e0b")))
	d.Add(g)
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, `id="layer-1"`) { t.Error("group id missing") }
}
func TestPath(t *testing.T) {
	d := doc.NewDocument(800, 600)
	d.Add(draw.NewPath(draw.Fill("#8b5cf6")).M(100, 200).L(300, 200).Z())
	svg, err := d.Render()
	if err != nil { t.Fatal(err) }
	if !strings.Contains(svg, "<path") { t.Error("path missing") }
}
func TestExport(t *testing.T) {
	d := doc.NewDocument(200, 200).WithBackground("#111")
	d.Add(draw.NewCircle(100, 100, 50, draw.Fill("#a78bfa")))
	if err := d.Export(t.TempDir() + "/out.svg"); err != nil {
		t.Fatal(err)
	}
}
func TestExportBadPath(t *testing.T) {
	if err := doc.NewDocument(100, 100).Export("/no/such/path/out.svg"); err == nil {
		t.Error("expected error for bad path")
	}
}
