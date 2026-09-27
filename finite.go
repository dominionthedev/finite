// Package finite is the entry point for building visuals programmatically.
//
// Design: you compose visuals (scenes, layers, shapes, widgets). SVG is the
// medium and system underneath — how those visuals are represented and
// exchanged — not the thing you have to think about first.
//
//	scene := finite.New(800, 600).WithBackground("#0a0a14")
//	scene.Layer("content").Add(draw.NewCircle(400, 300, 80, draw.Fill("#a78bfa")))
//	_ = scene.Export("out.svg")
//
// Core building blocks live in subpackages:
//   - doc: scene/document, layers
//   - draw: geometry, paint, effects
//   - geom: transforms
//   - widget / instance: reusable components
//   - decode / validate: import and audit
package finite

import "github.com/dominionthedev/finite/doc"

// Scene is a visual canvas. It is the primary type you build against.
// Rendering materializes the scene as SVG.
type Scene = doc.Document

// Layer is a named, ordered group within a scene.
type Layer = doc.Layer

// Renderable is anything that can be placed in a scene or layer.
type Renderable = doc.Renderable

// New creates a scene with the given viewport size in user units.
func New(width, height int) *Scene {
	return doc.NewDocument(width, height)
}
