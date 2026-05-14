# Module: widget

The `widget` module defines the contract for reusable visual components.

## `type Widget interface`
```go
type Widget interface {
    Name() string
    Render() (string, error)
}
```

## Creating a Custom Widget
Any type that implements the `Name()` and `Render()` methods can be used as a widget. Widgets are self-contained and should return a valid SVG fragment (typically wrapped in a `<g>` tag if it contains multiple elements).
