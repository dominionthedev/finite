# Module: doc

The `doc` module provides the primary entry point for Finite. It defines the `Document` struct, which acts as the SVG canvas.

## API Reference

### `type Renderable interface`
Anything that can produce an SVG fragment.
- `Render() (string, error)`

### `type Document struct`
The main SVG canvas that composes `Renderable` elements.
- `Width int`
- `Height int`
- `Background string` (optional hex color)

### `func NewDocument(width, height int) *Document`
Creates a new document with the specified dimensions.

### `func (d *Document) WithBackground(color string) *Document`
Sets the background color and returns the document for chaining.

### `func (d *Document) Add(r Renderable)`
Adds a `Renderable` element to the document. Elements are drawn in the order they are added.

### `func (d *Document) AddDef(def string)`
Injects a raw SVG `<defs>` fragment.

### `func (d *Document) Render() (string, error)`
Composes all elements into a single SVG string.

### `func (d *Document) Export(filePath string) error`
Renders the document and writes it to a file.
