# Module: decode

The `decode` module allows importing existing SVG files or raw SVG bytes into Finite's native document format.

## API Reference

### `func File(path string) (*doc.Document, error)`
Reads an SVG file from disk and returns a `doc.Document`.

### `func Bytes(data []byte) (*doc.Document, error)`
Parses raw SVG bytes into a `doc.Document`.

## How it works
The decoder parses the XML structure and maps standard SVG elements (`circle`, `rect`, `path`, etc.) to Finite's `draw` types. Unknown elements are preserved as `RawElement` to ensure fidelity.
