// Basic composition example — layers, geom, clip, filters, widgets, validate.
//
// Run: go run ./examples/basic
package main

import (
	"fmt"
	"os"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/geom"
	"github.com/dominionthedev/finite/instance"
	"github.com/dominionthedev/finite/validate"
	"github.com/dominionthedev/finite/widget"
)

// badge is a small reusable widget with namespaced ids.
type badge struct {
	label string
}

func (b badge) Name() string { return "badge" }

func (b badge) Render() (string, error) {
	s := widget.NewScope(b.Name())
	gradID := s.ID("bg")
	bg, err := draw.NewRoundedRect(0, 0, 88, 28, 8,
		draw.FillGradient(draw.Linear(gradID, "0%", "0%", "100%", "0%",
			draw.S("0%", "#6366f1", 1),
			draw.S("100%", "#a78bfa", 1),
		)),
	).Render()
	if err != nil {
		return "", fmt.Errorf("badge.Render: %w", err)
	}
	label, err := draw.NewText(b.label, 44, 19,
		draw.FontSize(12),
		draw.FontWeight("600"),
		draw.TextFill("#f8fafc"),
		draw.Centered(),
	).Render()
	if err != nil {
		return "", fmt.Errorf("badge.Render: %w", err)
	}
	return "<g>" + bg + label + "</g>", nil
}

func main() {
	canvas := doc.NewDocument(800, 600).WithBackground("#0a0a14")

	// ── Background glow (layer) ───────────────────────────────────────────
	canvas.Layer("atmosphere").Add(
		draw.NewEllipse(400, 280, 320, 240,
			draw.FillGradient(draw.Radial("orb-grad", "40%", "35%", "65%",
				draw.S("0%", "#ffffff", 0.08),
				draw.S("45%", "#a78bfa", 0.35),
				draw.S("100%", "#7c3aed", 0.0),
			)),
			draw.WithFilter(draw.Blur("orb-blur", 28)),
		),
	)

	// ── Accent line ───────────────────────────────────────────────────────
	canvas.Layer("chrome").Add(
		draw.NewRect(100, 480, 600, 1,
			draw.FillGradient(draw.Linear("bar-grad", "0%", "0%", "100%", "0%",
				draw.S("0%", "#6366f1", 0.0),
				draw.S("40%", "#a78bfa", 1.0),
				draw.S("100%", "#ec4899", 0.0),
			)),
			draw.Opacity(0.6),
		),
	)

	// ── Title with clip + glow ────────────────────────────────────────────
	title := draw.NewGroup(draw.GroupID("title"))
	title.Add(
		draw.NewText("finite", 400, 290,
			draw.FontSize(72),
			draw.FontWeight("700"),
			draw.LetterSpacing(-3),
			draw.TextFillGradient(draw.Linear("text-grad", "0%", "0%", "100%", "0%",
				draw.S("0%", "#c4b5fd", 1.0),
				draw.S("100%", "#f0abfc", 1.0),
			)),
			draw.TextFilter(draw.Glow("text-glow", 6)),
			draw.Centered(),
		),
		draw.NewText("programmable svg", 400, 350,
			draw.FontSize(16),
			draw.TextFill("#4b5563"),
			draw.Centered(),
		),
	)
	canvas.Layer("content").Add(title)

	// ── Isometric sample (geom) ───────────────────────────────────────────
	iso := draw.NewGroup(
		draw.GroupID("basiciso"),
		draw.GroupTransform(geom.Compose(geom.Translate(120, 440), geom.Isometric())),
	)
	iso.Add(
		draw.NewRect(0, 0, 36, 36, draw.Fill("#34d399")),
		draw.NewRect(0, -12, 36, 12, draw.Fill("#6ee7b7")),
	)
	canvas.Layer("iso").Add(iso)

	// ── Pattern fill + clip ───────────────────────────────────────────────
	canvas.Layer("texture").Add(
		draw.NewRect(620, 80, 120, 120,
			draw.FillPattern(draw.PatternDots("dots", 10, 1.5, "#818cf8")),
			draw.WithClip(draw.ClipCircle("tex-clip", 680, 140, 55)),
			draw.Opacity(0.5),
		),
	)

	// ── Widget instance ───────────────────────────────────────────────────
	canvas.Layer("ui").Add(
		instance.At(badge{label: "v0.1"}, 56, 40, instance.WithScale(1.0)),
	)

	// ── Validate before export ────────────────────────────────────────────
	report, err := validate.Document(canvas)
	if err != nil {
		fmt.Fprintln(os.Stderr, "validate error:", err)
		os.Exit(1)
	}
	if !report.Valid() {
		fmt.Fprintln(os.Stderr, report)
		os.Exit(1)
	}
	for _, w := range report.Warnings() {
		fmt.Println("warn:", w)
	}
	for _, i := range report.Infos() {
		fmt.Println("info:", i)
	}

	if err := canvas.Export("output.svg"); err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(1)
	}
	fmt.Println("✦ exported output.svg")
}
