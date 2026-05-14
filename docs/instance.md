# Module: instance

The `instance` module provides a way to place `Widget` components on a document with geometric transforms.

## API Reference

### `func At(w widget.Widget, x, y float64, opts ...Option) *Instance`
Places a widget at a specific `(x, y)` coordinate.

### Options
- `WithScale(s float64)`: Scales the widget.
- `WithRotation(deg float64)`: Rotates the widget in degrees.

## Usage
`Instance` implements the `doc.Renderable` interface, so it can be added directly to a `doc.Document`.
