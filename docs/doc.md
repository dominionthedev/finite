# Module: doc

The `doc` module is the primary entry point for Finite. It defines the `Document`
canvas, the `Renderable` interface, and a lightweight **layer** system for
structured composition.

## Core types

### `Renderable`

Anything that can produce an SVG fragment:

```go
type Renderable interface {
    Render() (string, error)
}
```

Shapes, text, groups, layers, and widget instances all implement this.

### `Document`

The SVG canvas. It manages:

- Viewport size (`Width`, `Height`)
- User coordinate system (`viewBox`, `preserveAspectRatio`)
- Background color
- Root renderables (via `Add`)
- Named layers (via `Layer`)
- Shared definitions (`AddDef`)

```go
canvas := doc.NewDocument(800, 600).
    WithBackground("#0a0a14").
    WithViewBox(0, 0, 800, 600)

canvas.Layer("bg").Add(draw.NewRect(0, 0, 800, 600, draw.Fill("#111")))
canvas.Layer("content").Add(
    draw.NewCircle(400, 300, 80, draw.Fill("#a78bfa")),
)
canvas.Export("out.svg")
```

### `Layer`

A named, ordered group of renderables. Layers are the preferred way to structure
a document:

- Stable SVG `id` (the layer name)
- Creation order = draw order
- Shared transform / opacity via the underlying `draw.Group`
- Retrievable by name: `doc.Layer("content")` returns the existing layer

```go
layer := canvas.Layer("iso")
layer.SetTransform(geom.Isometric())
layer.Add(draw.NewRect(0, 0, 40, 40, draw.Fill("#f59e0b")))
```

## Composition model

1. **Root elements** — added with `Document.Add`, drawn first
2. **Layers** — named groups, drawn after root elements in creation order
3. **Groups** (`draw.Group`) — nestable containers used inside layers or at root

Prefer layers for document structure. Use groups for local nesting and
shared transforms inside a layer.
