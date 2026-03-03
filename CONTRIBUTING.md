# Contributing to Alolyte

Thanks for your interest in contributing! Alolyte is early-stage, so the best contributions right now are:

- Bug reports and feedback on the core API
- New widgets (see below)
- Examples and use cases

---

## Getting Started

```bash
git clone https://github.com/leraniode/alolyte
cd alolyte
go test ./...
```

No external dependencies in the core runtime. Everything compiles with the Go standard library.

---

## Contributing a Widget

Widgets live in `assets/widgets/` as `.svg` files. A widget is a self-contained SVG that uses Go `text/template` syntax for its parameters.

### Rules for widgets

1. **Always provide a default for every parameter** — so the widget looks good with zero config.

```xml
<!-- Good -->
<stop stop-color="{{if .PrimaryColor}}{{.PrimaryColor}}{{else}}#a78bfa{{end}}"/>

<!-- Bad — will render empty if PrimaryColor is not passed -->
<stop stop-color="{{.PrimaryColor}}"/>
```

2. **Use `{{.ID}}` to namespace all SVG ids** — prevents collisions when a widget is placed multiple times.

```xml
<radialGradient id="orb-grad-{{.ID}}">
  ...
</radialGradient>
<ellipse fill="url(#orb-grad-{{.ID}})"/>
```

3. **Keep width/height in the `<svg>` tag** — the Instance renderer uses it for layout context.

4. **Animations are welcome** — use `<animate>` and `<animateTransform>`. Make speed/duration a param with a sensible default.

---

## Code Style

- Standard `gofmt` formatting — no exceptions.
- Errors are always wrapped with context: `fmt.Errorf("package.Func: %w", err)`.
- No external dependencies in `doc`, `widget`, or `instance` packages.
- CLI and tooling may use external dependencies (Cobra, Lipgloss, Mosaics).

---

## Commit Style

```
feat: add geometric_ring widget
fix: strip self-closing svg wrapper correctly
docs: add widget authoring guide
test: add registry scan tests
```

---

## Questions

Open an issue. This is a small project — we'll respond fast.
