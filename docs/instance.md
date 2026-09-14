# Module: instance

Place `widget.Widget` values on a document with geometric transforms.

## API

```go
instance.At(w, x, y, opts...) *Instance
```

Options (applied in local space, before translation):

- `WithScale(s)`
- `WithScaleXY(sx, sy)`
- `WithRotation(deg)`
- `WithRotationAround(deg, cx, cy)`
- `WithTransform(m geom.Matrix)` — replaces local transform

`Instance` implements `doc.Renderable`, so it can be added to a document or layer.

```go
canvas.Layer("ui").Add(
    instance.At(myWidget, 40, 40, instance.WithScale(0.8)),
)
```
