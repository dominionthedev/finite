# Contributing to finite

## Widget authoring rules

1. Implement `widget.Widget` — `Name() string` and `Render() (string, error)`
2. Always produce clean, valid SVG with no template syntax
3. Namespace all internal SVG `id` attributes with `widget.NewScope` / `NewScopeID`
   so multiple instances do not collide
4. Inline all `<defs>` (gradients, filters) within the rendered fragment —
   don't depend on document-level defs
5. Prefer `widget.Composite` when assembling larger components from smaller ones
6. Return `fmt.Errorf("pkg.Type: %w", err)` — always wrap errors with context

## Error format

```go
return "", fmt.Errorf("mywidget.Render: %w", err)
```

## No external dependencies in core packages

`doc`, `draw`, `widget`, `instance`, `decode`, `validate` — zero external deps.
Standard library only.

## Testing

Every exported type needs at least one test. Run:

```
go test ./...
```
