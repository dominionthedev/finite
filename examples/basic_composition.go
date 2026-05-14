// Basic composition example — Finite's native draw API.
//
// Shows shapes, gradients, filters, text, groups, decode, and validate
// working together in a single document.
//
// Run: go run examples/basic_composition.go
package main

import (
	"fmt"
	"os"

	"github.com/dominionthedev/finite/decode"
	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/validate"
)

func main() {
	canvas := doc.NewDocument(800, 600).WithBackground("#0a0a14")

	// ── Background glow ────────────────────────────────────────────────────
	canvas.Add(draw.NewEllipse(400, 280, 320, 240,
		draw.FillGradient(draw.Radial("orb-grad", "40%", "35%", "65%",
			draw.S("0%", "#ffffff", 0.08),
			draw.S("45%", "#a78bfa", 0.35),
			draw.S("100%", "#7c3aed", 0.0),
		)),
		draw.WithFilter(draw.Blur("orb-blur", 28)),
	))

	// ── Accent line ────────────────────────────────────────────────────────
	canvas.Add(draw.NewRect(100, 480, 600, 1,
		draw.FillGradient(draw.Linear("bar-grad", "0%", "0%", "100%", "0%",
			draw.S("0%", "#6366f1", 0.0),
			draw.S("40%", "#a78bfa", 1.0),
			draw.S("100%", "#ec4899", 0.0),
		)),
		draw.Opacity(0.6),
	))

	// ── Title ──────────────────────────────────────────────────────────────
	canvas.Add(draw.NewText("finite", 400, 290,
		draw.FontSize(72),
		draw.FontWeight("700"),
		draw.LetterSpacing(-3),
		draw.TextFillGradient(draw.Linear("text-grad", "0%", "0%", "100%", "0%",
			draw.S("0%", "#c4b5fd", 1.0),
			draw.S("100%", "#f0abfc", 1.0),
		)),
		draw.TextFilter(draw.Glow("text-glow", 6)),
		draw.Centered(),
	))

	canvas.Add(draw.NewText("programmable svg", 400, 350,
		draw.FontSize(16),
		draw.TextFill("#4b5563"),
		draw.Centered(),
	))

	// ── Decode an existing SVG and composite on top ────────────────────────
	if existing, err := decode.File("assets/example_input.svg"); err == nil {
		// merge all its elements into our canvas
		svg, _ := existing.Render()
		canvas.Add(&draw.RawElement{Content: svg})
	}

	// ── Validate before export ─────────────────────────────────────────────
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

	// ── Export ─────────────────────────────────────────────────────────────
	if err := canvas.Export("output.svg"); err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(1)
	}
	fmt.Println("✦ exported output.svg")
}
