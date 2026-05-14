// Package widget defines the core interface for reusable visual components in Finite.
package widget

// Widget is the interface that wraps the basic Name and Render methods.
// Reusable components like icons, UI elements, or complex illustrations should implement this.
type Widget interface {
	// Name returns a unique or stable identifier for the widget type.
	Name() string

	// Render generates a valid SVG fragment (without the <svg> root element).
	Render() (string, error)
}
