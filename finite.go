// Package finite is the entry point for building visuals programmatically.
//
// Design: you compose visuals (layers, shapes, widgets). SVG is the medium
// and system underneath — how those visuals are represented and exchanged —
// not the thing you have to think about first.
//
//	visual := finite.New(800, 600).WithBackground("#0a0a14")
//	visual.Layer("content").Add(draw.NewCircle(400, 300, 80, draw.Fill("#a78bfa")))
//	_ = visual.Export("out.svg")
//
// Core building blocks live in subpackages:
//   - doc: visual/document, layers
//   - draw: geometry, paint, effects
//   - geom: transforms
//   - widget / instance: reusable components
//   - decode / validate: import and audit
//
// Naming note: Finite deals in visuals. Higher-level composition of mixed
// media fragments into editorial "scenes" is a separate concern (e.g. Alolyte).
package finite

import "github.com/dominionthedev/finite/doc"

// Visual is a canvas for programmatic graphics. It is the primary type you
// build against. Rendering materializes the visual as SVG.
type Visual = doc.Document

// Layer is a named, ordered group within a visual.
type Layer = doc.Layer

// Renderable is anything that can be placed in a visual or layer.
type Renderable = doc.Renderable

// New creates a visual with the given viewport size in user units.
func New(width, height int) *Visual {
	return doc.NewDocument(width, height)
}
