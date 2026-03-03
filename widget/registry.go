package widget

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Registry holds all widgets discovered from a directory.
// Drop any .svg file into the watched folder and it becomes available.
type Registry struct {
	Dir     string
	widgets map[string]Widget
}

// LoadRegistry scans dir for .svg files and returns a populated Registry.
func LoadRegistry(dir string) (*Registry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("widget.LoadRegistry: cannot read %q: %w", dir, err)
	}

	r := &Registry{
		Dir:     dir,
		widgets: make(map[string]Widget),
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".svg") {
			continue
		}
		w, err := Load(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		r.widgets[w.Name] = w
	}

	return r, nil
}

// Get returns a widget by name, or an error if not found.
func (r *Registry) Get(name string) (Widget, error) {
	w, ok := r.widgets[name]
	if !ok {
		return Widget{}, fmt.Errorf("widget.Registry: %q not found in %q", name, r.Dir)
	}
	return w, nil
}

// List returns all widget names in the registry.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.widgets))
	for name := range r.widgets {
		names = append(names, name)
	}
	return names
}

// MustGet is like Get but panics on error. Useful in design scripts.
func (r *Registry) MustGet(name string) Widget {
	w, err := r.Get(name)
	if err != nil {
		panic(err)
	}
	return w
}
