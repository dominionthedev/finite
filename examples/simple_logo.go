// Simple logo design example — using paths and gradients.
package main

import (
	"fmt"
	"os"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
)

func main() {
	canvas := doc.NewDocument(500, 500).WithBackground("#ffffff")

	// Logo Mark (Abstract Geometric Shape)
	canvas.Add(draw.NewPath(
		draw.FillGradient(draw.Linear("logo-grad", "0%", "0%", "0%", "100%",
			draw.S("0%", "#3b82f6", 1.0),
			draw.S("100%", "#8b5cf6", 1.0),
		)),
	).M(250, 100).L(400, 250).L(250, 400).L(100, 250).Z())

	// Inner cut-out
	canvas.Add(draw.NewCircle(250, 250, 40, draw.Fill("#ffffff")))

	// Brand Name
	canvas.Add(draw.NewText("FINITE", 250, 460,
		draw.FontSize(32),
		draw.FontWeight("900"),
		draw.LetterSpacing(8),
		draw.TextFill("#1f2937"),
		draw.Centered(),
	))

	if err := canvas.Export("logo.svg"); err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(1)
	}
	fmt.Println("✦ exported logo.svg")
}
