// Data visualization example — generating a bar chart programmatically.
package main

import (
	"fmt"
	"os"

	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
)

func main() {
	canvas := doc.NewDocument(800, 400).WithBackground("#f9fafb")

	data := []float64{45, 80, 55, 90, 70, 100, 65}
	colors := []string{"#f87171", "#fb923c", "#fbbf24", "#34d399", "#22d3ee", "#818cf8", "#c084fc"}

	margin := 60.0
	chartWidth := 800.0 - (margin * 2)
	chartHeight := 400.0 - (margin * 2)
	barWidth := (chartWidth / float64(len(data))) - 10

	for i, val := range data {
		h := (val / 100.0) * chartHeight
		x := margin + float64(i)*(barWidth+10)
		y := 400.0 - margin - h

		canvas.Add(draw.NewRect(x, y, barWidth, h,
			draw.Fill(colors[i%len(colors)]),
			draw.Opacity(0.9),
		))

		// Value labels
		canvas.Add(draw.NewText(fmt.Sprintf("%.0f", val), x+(barWidth/2), y-10,
			draw.FontSize(12),
			draw.TextFill("#4b5563"),
			draw.Centered(),
		))
	}

	// X-Axis line
	canvas.Add(draw.NewLine(margin, 400-margin, 800-margin, 400-margin,
		draw.Stroke("#9ca3af", 2),
	))

	if err := canvas.Export("chart.svg"); err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(1)
	}
	fmt.Println("✦ exported chart.svg")
}
