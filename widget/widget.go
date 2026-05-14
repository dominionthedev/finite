// Package widget defines the Widget interface — the contract every
// reusable Finite component must satisfy.
//
// A Widget is a self-contained visual component that knows how to render
// itself to a valid SVG fragment. It has a stable name for identification
// and produces clean, standard SVG with no template syntax visible to callers.
//
// Writing a custom widget:
//
//	type Hexagon struct {
//	    CX, CY float64
//	    Size    float64
//	    Color   string
//	}
//
//	func (h *Hexagon) Name() string { return "hexagon" }
//
//	func (h *Hexagon) Render() (string, error) {
//	    // compute points, return clean <polygon> SVG
//	}
//
//	canvas.Add(&Hexagon{CX: 300, CY: 300, Size: 80, Color: "#a78bfa"})
package widget

// Widget is anything that has a stable name and can render itself
// to a clean SVG fragment.
//
// Both builtin components and user-defined types satisfy this interface.
// A Widget can be added directly to a doc.Document or placed via instance.At().
type Widget interface {
	// Name returns a stable identifier for this widget type.
	// Used for debugging, logging, and the instance package.
	Name() string

	// Render produces a clean, valid SVG fragment (no <svg> wrapper).
	// The fragment is embedded directly inside the parent document's <svg>.
	// All internal defs (gradients, filters) must be inlined within the fragment.
	Render() (string, error)
}
