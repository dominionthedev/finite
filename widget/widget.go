// Package widget defines the contract and helpers for reusable visual components.
//
// A Widget is a self-contained visual unit. It owns its internal geometry,
// styles, and definition IDs. Instances of widgets are placed on a document
// via the instance package.
//
// Authoring rules (see CONTRIBUTING.md):
//   - Namespace all internal SVG ids (use Scope)
//   - Inline defs (gradients, filters) inside the fragment
//   - Return clean SVG with no template syntax
//   - Wrap errors: fmt.Errorf("pkg.Type: %w", err)
package widget

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// Widget is a reusable visual component.
//
// Name identifies the widget type (stable, not instance-specific).
// Render returns a valid SVG fragment — typically a <g> containing the
// widget's content and any inline <defs>. It must not include an <svg> root.
type Widget interface {
	Name() string
	Render() (string, error)
}

// scopeCounter ensures Scope prefixes stay unique across a process when
// the caller does not supply an explicit instance id.
var scopeCounter atomic.Uint64

// Scope helps widgets namespace internal SVG ids so multiple instances
// of the same widget do not collide (gradients, filters, clipPaths, etc.).
//
// Typical use inside a widget:
//
//	func (w *Card) Render() (string, error) {
//	    s := widget.NewScope(w.Name())
//	    gradID := s.ID("bg")
//	    // ... use gradID in fill="url(#...)" and in the <linearGradient id="...">
//	}
type Scope struct {
	prefix string
}

// NewScope creates a Scope from a widget name. A unique suffix is appended
// so concurrent instances do not share ids.
func NewScope(widgetName string) *Scope {
	n := scopeCounter.Add(1)
	safe := sanitize(widgetName)
	return &Scope{prefix: fmt.Sprintf("%s-%d", safe, n)}
}

// NewScopeID creates a Scope with an explicit instance id (no auto suffix).
// Use this when the caller controls identity and wants stable ids.
func NewScopeID(widgetName, instanceID string) *Scope {
	return &Scope{prefix: sanitize(widgetName) + "-" + sanitize(instanceID)}
}

// ID returns a namespaced id: prefix + "-" + suffix.
func (s *Scope) ID(suffix string) string {
	return s.prefix + "-" + sanitize(suffix)
}

// Prefix returns the scope prefix (useful for debugging).
func (s *Scope) Prefix() string {
	return s.prefix
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "w"
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	return b.String()
}

// Func adapts a name + render function into a Widget.
// Useful for small one-off components without a dedicated type.
type Func struct {
	N string
	F func() (string, error)
}

// Name returns the widget name.
func (f Func) Name() string {
	if f.N == "" {
		return "func"
	}
	return f.N
}

// Render calls the underlying function.
func (f Func) Render() (string, error) {
	if f.F == nil {
		return "", fmt.Errorf("widget.Func %q: nil render func", f.Name())
	}
	svg, err := f.F()
	if err != nil {
		return "", fmt.Errorf("widget.Func %q: %w", f.Name(), err)
	}
	return svg, nil
}

// Static is a Widget that always returns a fixed SVG fragment.
type Static struct {
	N       string
	Content string
}

// Name returns the widget name.
func (s Static) Name() string {
	if s.N == "" {
		return "static"
	}
	return s.N
}

// Render returns the fixed content.
func (s Static) Render() (string, error) {
	return s.Content, nil
}
