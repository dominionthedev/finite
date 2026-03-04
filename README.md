# ✦ Alolyte

> Programmable SVG design library for Go.

Alolyte lets you compose high-quality, reusable, and animated SVG graphics using code — no GUI, no raw XML editing, just Go.

It sits between "hand-edit SVG by hand" and "use Figma" — a code-first design runtime for developers who want programmable, composable visuals.

---

## Status

> In Development. The core runtime is in place. Widgets, CLI, and animations are coming in v0.1.

---

## Concept

```
Widget   →   Instance   →   Document   →   SVG file
(SVG file)   (placed on     (the canvas)
              canvas with
              transforms
              + params)
```

**Widget** — a reusable `.svg` file with Go `text/template` variables for parameterization.

**Instance** — a widget placed at a position on the canvas, with scale, rotation, and param overrides.

**Document** — the canvas. Holds all instances and exports the final composed SVG.

---

## API Preview

```go
canvas := doc.NewDocument(800, 600).WithBackground("#0f0f1a")

bg, _  := widget.Load("assets/widgets/aurora_background.svg")
orb, _ := widget.Load("assets/widgets/gradient_orb.svg")

canvas.Add(instance.At(bg, 0, 0,
    instance.WithParams(instance.Params{
        "ID": "bg", "ColorA": "#6366f1", "ColorB": "#ec4899",
    }),
))

canvas.Add(instance.At(orb, 200, 100,
    instance.WithScale(1.4),
    instance.WithParams(instance.Params{
        "ID": "orb1", "PrimaryColor": "#ffffff", "SecondaryColor": "#a78bfa",
    }),
))

canvas.Export("output.svg")
```

Widget SVGs use Go template syntax for parameters:

```xml
<stop offset="0%"
  stop-color="{{if .PrimaryColor}}{{.PrimaryColor}}{{else}}#a78bfa{{end}}"/>
```

---

## Project Structure

```
alolyte/
├── doc/              # Document — canvas, composition, export
├── widget/           # Widget — SVG loader, template renderer, registry
├── instance/         # Instance — placement, transforms, param overrides
├── assets/
│   └── widgets/      # Built-in SVG widgets (coming in v0.1)
└── examples/         # Usage examples
```

---

## Roadmap

### v0.1

- [x] `widget` — SVG file loader with Go template param support
- [x] `widget` — Registry to scan and index a widgets folder
- [x] `instance` — Placement with position, scale, rotation, params
- [x] `doc` — Document canvas with `Add`, `Render`, `Export`
- [ ] Tests for all three packages
- [ ] Built-in widget collection (`gradient_orb`, `aurora_background`, `geometric_ring`, `glow_text`, `noise_texture`)
- [ ] Example usage in `examples/` folder
- [ ] CLI: `alolyte list`, `alolyte preview`, `alolyte render`, `alolyte new`
- [ ] Terminal image preview via Charm Mosaics
- [ ] First real example output — animated logo mark

### v0.2+

- [ ] Widget manifest (`widget.yaml`) — metadata, param schema, preview thumbnail
- [ ] `alolyte pull <pack>` — widget packs from GitHub
- [ ] Animation props (delay, duration, easing as params)
- [ ] DSL syntax sugar for common composition patterns

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

MIT © [leraniode](https://github.com/leraniode)

---

Part of [Leraniode](https://github.com/leraniode) – Building Tools that feel alive 🌱.

<p align="left">
    <a href="https://github.com/leraniode">
       <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/footer/leraniodeproductbrandimage.png" alt="Leraniode" width="600">
    </a>
</p>