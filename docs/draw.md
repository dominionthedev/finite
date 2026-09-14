# Module: draw

The `draw` module contains the primitive SVG shapes, text, and styling options.

## Shapes

- `Circle`: `<circle>` element.
- `Ellipse`: `<ellipse>` element.
- `Rect`: `<rect>` element (supports rounded corners).
- `Line`: `<line>` element.
- `Path`: `<path>` element with a fluent builder API (`M`, `L`, `C`, `Q`, `Z`).

## Text

- `Text`: Styled `<text>` element with support for alignment, font properties, and filters.

## Styling (Options)

Finite uses a functional options pattern for styling:
- `Fill(color string)`
- `FillGradient(g fillProvider)`
- `Stroke(color string, width float64)`
- `Opacity(v float64)`
- `WithFilter(f *Filter)`
- `Transform(m geom.Matrix)` / `TransformString(t string)`

## Gradients

- `Linear(id, x1, y1, x2, y2 string, stops ...Stop)`
- `Radial(id, cx, cy, r string, stops ...Stop)`

## Filters

- `Blur(id string, stdDev float64)`
- `Glow(id string, radius float64)`
- `Noise(id string, frequency float64, octaves int, blendMode string)`
- `Shadow(id string, dx, dy, blur float64, color string, opacity float64)`


## Groups

`Group` is the primary composition primitive inside a document or layer.

```go
g := draw.NewGroup(draw.GroupID("cluster"), draw.GroupOpacity(0.9))
g.Translate(100, 80).Rotate(15)
g.Add(
    draw.NewCircle(0, 0, 20, draw.Fill("#a78bfa")),
    draw.NewCircle(30, 0, 12, draw.Fill("#c4b5fd")),
)
canvas.Add(g)
```

Options: `GroupID`, `GroupTransform`, `GroupTransformString`, `GroupOpacity`, `GroupFilter`.
Chainable methods: `Add`, `Translate`, `Rotate`, `RotateAround`, `Scale`, `ScaleUniform`, `SetTransform`.
