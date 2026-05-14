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
- `Transform(t string)`

## Gradients

- `Linear(id, x1, y1, x2, y2 string, stops ...Stop)`
- `Radial(id, cx, cy, r string, stops ...Stop)`

## Filters

- `Blur(id string, stdDev float64)`
- `Glow(id string, radius float64)`
- `Noise(id string, frequency float64, octaves int, blendMode string)`
- `Shadow(id string, dx, dy, blur float64, color string, opacity float64)`
