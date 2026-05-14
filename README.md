# Finite ✦ Programmable SVG

Finite is a lightweight, composable SVG engine for Go. It allows you to build complex, styled, and validated SVG documents using a clean, fluent API.

## Features

- **Fluent API**: Build shapes, paths, and text with intuitive method chaining.
- **Styling**: Support for solid fills, linear/radial gradients, and SVG filters (blur, glow, shadow).
- **Import/Export**: Decode existing SVGs, modify them, and export back to clean XML.
- **Validation**: Built-in SVG auditor for duplicate IDs, unresolved references, and XML syntax.
- **Widgets**: Reusable visual components that can be transformed and composited.

## Installation

```bash
go get github.com/dominionthedev/finite
```

## Quick Start

```go
package main

import (
	"github.com/dominionthedev/finite/doc"
	"github.com/dominionthedev/finite/draw"
)

func main() {
	// Create a canvas
	canvas := doc.NewDocument(800, 600).WithBackground("#0a0a14")

	// Add a styled circle
	canvas.Add(draw.NewCircle(400, 300, 100,
		draw.Fill("#a78bfa"),
		draw.WithFilter(draw.Glow("main-glow", 8)),
	))

	// Export to file
	canvas.Export("output.svg")
}
```

## Documentation

Detailed module documentation can be found in the `docs/` directory:

- [doc](./docs/doc.md): The main SVG canvas and composition engine.
- [draw](./docs/draw.md): Shapes, text, gradients, and filters.
- [decode](./docs/decode.md): Importing existing SVGs.
- [instance](./docs/instance.md): Positioning and transforming components.
- [validate](./docs/validate.md): Auditing SVG documents.
- [widget](./docs/widget.md): Reusable component contract.

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

<p align="center">DominionDev</p>
<p align="center">
<a href="https://github.com/dominionthedev">GitHub</a> • <a href="https://dominionthedev.github.io">Website</a>
</p>

<p align="center">
    <a href="">
        <img src="https://raw.githubusercontent.com/dominionthedev/dominionthedev/main/assets/watermark-animated.svg"/>
    </a>
</p>

