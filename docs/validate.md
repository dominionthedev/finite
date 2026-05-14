# Module: validate

The `validate` module provides tools to audit an SVG document for common issues before export.

## Checks Performed
1. **Well-formed XML**: Ensures the SVG is valid XML.
2. **Duplicate IDs**: Detects if the same `id` attribute is used more than once.
3. **Unresolved URL References**: Checks if `url(#id)` references point to existing elements.
4. **Zero Dimensions**: Warns about elements with zero width/height/radius that will be invisible.

## API Reference

### `func Document(d *doc.Document) (*Report, error)`
Validates a Finite document.

### `func SVG(svg string) (*Report, error)`
Validates a raw SVG string.
