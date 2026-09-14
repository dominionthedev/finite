# Module: widget

Reusable visual components for Finite.

## `Widget` interface

```go
type Widget interface {
    Name() string
    Render() (string, error)
}
```

- `Name` — stable type identifier (not instance-specific)
- `Render` — valid SVG fragment (no `<svg>` root). Prefer a wrapping `<g>`.

## ID namespacing (`Scope`)

Multiple instances of the same widget must not collide on gradient/filter ids.

```go
func (w *Card) Render() (string, error) {
    s := widget.NewScope(w.Name())
    grad := s.ID("bg") // e.g. "card-3-bg"
    // use grad in url(#...) and in the <linearGradient id="...">
}
```

- `NewScope(name)` — unique prefix per call
- `NewScopeID(name, instanceID)` — stable, caller-controlled prefix

## Built-in helpers

| Type | Purpose |
|------|---------|
| `Func` | Adapt `func() (string, error)` as a Widget |
| `Static` | Fixed SVG fragment |
| `Composite` | Compose child widgets with optional transform/opacity |

```go
card := widget.NewComposite("card", []widget.Widget{
    widget.Static{Content: `<rect width="120" height="60" rx="8" fill="#1e1e2e"/>`},
    widget.Func{N: "label", F: func() (string, error) {
        return `<text y="35" fill="#cdd6f4">Hello</text>`, nil
    }},
}, widget.WithOpacity(0.95))
```

## Placement

Use `instance.At` to place a widget on a document:

```go
canvas.Add(instance.At(card, 100, 80, instance.WithScale(1.2)))
```

## Authoring rules

See CONTRIBUTING.md. Summary:

1. Implement `Widget`
2. Namespace internal ids with `Scope`
3. Inline all defs inside the fragment
4. Wrap errors with context
