# Module: validate

Audit SVG output for structural problems before export.

## Checks

| Severity | Check |
|----------|--------|
| **error** | Malformed XML |
| **error** | Duplicate `id` values (includes guidance when layer names collide with child `GroupID`s) |
| **error** | `url(#…)` references to missing ids |
| **error** | Empty input |
| **warning** | Missing `xmlns` / `viewBox` on root |
| **warning** | Zero width/height/radius (invisible geometry) |
| **info** | Unused paint servers / effects (`linearGradient`, `filter`, `clipPath`, …) |
| **info** | Empty `<g>` groups |
| **info** | Document summary (element / group / id counts) |

## API

```go
report, err := validate.Document(canvas)
if err != nil { /* render failed */ }
if !report.Valid() {
    fmt.Fprintln(os.Stderr, report) // prints all severities
    os.Exit(1)
}
for _, w := range report.Warnings() { fmt.Println("warn:", w) }
for _, i := range report.Infos()    { fmt.Println("info:", i) }
```

### Report helpers

- `Valid()` — no errors
- `Empty()` — no issues at any severity
- `Errors()` / `Warnings()` / `Infos()` / `Issues()`
