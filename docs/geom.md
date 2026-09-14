# Module: geom

The `geom` package is the foundation for all positioning and transformation in Finite.

It provides a proper 2D affine matrix type and helpers so that transforms can be composed in code instead of by concatenating SVG transform strings.

## Core types

### `Matrix`

Represents the SVG/CSS affine matrix:

```
| a  c  e |
| b  d  f |
| 0  0  1 |
```

- `Identity()`, `Translate`, `Scale`, `ScaleUniform`, `Rotate`, `RotateAround`, `SkewX`, `SkewY`
- `Mul` — matrix multiplication
- `TransformPoint` — apply the matrix to a point
- `String()` — produces a compact, readable SVG `transform` attribute value when possible

### `Transform` builder

Fluent builder for accumulating operations:

```go
tf := geom.NewTransform().
    Translate(100, 80).
    Rotate(30).
    ScaleUniform(1.2)

m := tf.Matrix()
```

### Projection helpers

- `Isometric()` — common isometric projection (rotate 30° → skewX -30° → scale Y by √3/2)
- `Oblique(degrees)` — simple oblique projection via horizontal skew

## Integration

- `draw.Group` accepts `geom.Matrix` via `GroupTransform` and has chainable `Translate` / `Rotate` / `Scale` methods
- `draw.Transform(m geom.Matrix)` is the preferred way to attach a transform to a shape
- `instance.At` composes local scale/rotation first, then translation
- `doc.Document` now has explicit `WithViewBox` and `WithPreserveAspectRatio` for proper coordinate system control

## Design notes

- Zero external dependencies
- Prefer building matrices with the API over parsing SVG transform strings
- `String()` intentionally prefers human-readable forms (`translate`, `scale`, `rotate`) when the matrix matches those patterns; otherwise it emits `matrix(...)`
