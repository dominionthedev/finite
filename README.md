# Finite ✦ Programmable SVG

Finite is a Go framework for constructing visuals with SVG as the medium.
Build diagrams, UI marks, data graphics, and illustrations through a composable
scene model — transforms, layers, widgets, and a full styling stack.

## Features

- **Scene composition** — documents, named layers, and nestable groups
- **Geometry** — affine matrices, fluent transforms, isometric / oblique helpers
- **Drawing** — shapes, paths, text, gradients, filters, clip paths, masks, patterns
- **Widgets** — reusable components with ID namespacing and composition
- **Import / validate** — decode existing SVG; audit IDs and references
- **Zero external deps** in core packages (standard library only)

## Install

```bash
go get github.com/dominionthedev/finite
```

## Quick start

```go
package main

import (
	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
	"github.com/dominionthedev/finite/geom"
)

func main() {
	canvas := doc.NewDocument(800, 600).WithBackground("#0a0a14")

	// Named layers for structure
	canvas.Layer("glow").Add(
		draw.NewEllipse(400, 280, 280, 200,
			draw.FillGradient(draw.Radial("orb", "40%", "35%", "65%",
				draw.S("0%", "#ffffff", 0.08),
				draw.S("100%", "#7c3aed", 0.0),
			)),
			draw.WithFilter(draw.Blur("orb-blur", 24)),
		),
	)

	// Group with transform
	mark := draw.NewGroup(draw.GroupID("mark")).
		Translate(400, 300).
		Rotate(-8)
	mark.Add(
		draw.NewCircle(0, 0, 64,
			draw.Fill("#a78bfa"),
			draw.WithClip(draw.ClipCircle("mark-clip", 0, 0, 64)),
		),
	)
	canvas.Layer("content").Add(mark)

	// Isometric tile via geom
	tile := draw.NewGroup(draw.GroupTransform(geom.Isometric())).
		Translate(120, 420)
	tile.Add(draw.NewRect(0, 0, 40, 40, draw.Fill("#34d399")))
	canvas.Layer("iso").Add(tile)

	_ = canvas.Export("output.svg")
}
```

## Packages

| Package | Role |
|---------|------|
| [`doc`](./docs/doc.md) | Document canvas, layers, viewBox |
| [`draw`](./docs/draw.md) | Shapes, text, gradients, filters, clip, mask, pattern, groups |
| [`geom`](./docs/geom.md) | Affine matrices and transforms |
| [`widget`](./docs/widget.md) | Reusable components (`Scope`, `Composite`, …) |
| [`instance`](./docs/instance.md) | Place widgets with transforms |
| [`decode`](./docs/decode.md) | Import existing SVG |
| [`validate`](./docs/validate.md) | Audit IDs and references |

## Examples

```bash
go run ./examples/basic   # layers, filters, validate
go run ./examples/logo    # path mark + type
go run ./examples/chart   # simple bar chart
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md). Core packages stay dependency-free.

## License

MIT — see [LICENSE](./LICENSE).

<p align="center">DominionDev</p>
