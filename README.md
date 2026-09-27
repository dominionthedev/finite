# Finite ✦ Programmable visuals (SVG medium)

Finite is a Go library for **programmatically creating visuals**. SVG is the
medium and system behind those visuals — the representation layer — not the
primary mental model.

You compose scenes, layers, shapes, and widgets in Go. Finite turns that
composition into clean, scalable SVG.

## Features

- **Scene model** — `finite.New` → layers → groups → shapes / widgets
- **Geometry** — affine transforms, isometric / oblique helpers
- **Drawing** — shapes, paths (incl. arcs), polygons, polylines, text
- **Paint & effects** — fills, strokes (dash / cap / join), gradients, filters, clip, mask, pattern
- **Widgets** — reusable components with ID namespacing
- **Import / validate** — decode existing SVG; audit IDs and references
- **Accessibility** — title, desc, role, aria
- **Zero external deps** in core packages

## Install

```bash
go get github.com/dominionthedev/finite
```

## Quick start

```go
package main

import (
	"github.com/dominionthedev/finite"
	"github.com/dominionthedev/finite/draw"
)

func main() {
	scene := finite.New(800, 600).
		WithBackground("#0a0a14").
		WithTitle("Orb").
		WithRole("img")

	scene.Layer("content").Add(
		draw.NewCircle(400, 300, 100,
			draw.Fill("#a78bfa"),
			draw.WithFilter(draw.Glow("g", 8)),
		),
		draw.NewPolyline([]draw.Point{
			draw.Pt(100, 500), draw.Pt(400, 480), draw.Pt(700, 500),
		}, draw.Stroke("#6366f1", 2), draw.StrokeDash("6 4")),
	)

	_ = scene.Export("output.svg")
}
```

## Packages

| Package | Role |
|---------|------|
| `finite` | Entry: `New` → `Scene` |
| [`doc`](./docs/doc.md) | Scene/document, layers, viewBox, a11y |
| [`draw`](./docs/draw.md) | Shapes, text, paint, effects, groups |
| [`geom`](./docs/geom.md) | Affine matrices and transforms |
| [`widget`](./docs/widget.md) | Reusable components |
| [`instance`](./docs/instance.md) | Place widgets with transforms |
| [`decode`](./docs/decode.md) | Import existing SVG |
| [`validate`](./docs/validate.md) | Audit IDs and references |

## Examples

```bash
go run ./examples/basic
go run ./examples/logo
go run ./examples/chart
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md). Core packages stay dependency-free.

## License

MIT — see [LICENSE](./LICENSE).
