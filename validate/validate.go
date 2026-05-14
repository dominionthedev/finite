// Package validate checks a Finite document or raw SVG string for
// structural problems before export.
//
//	report, err := validate.Document(canvas)
//	if !report.Valid() {
//	    fmt.Println(report)
//	    os.Exit(1)
//	}
package validate

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/dominionthedev/finite/doc"
)

// Severity classifies how serious an issue is.
type Severity int

const (
	SeverityError   Severity = iota // blocks export
	SeverityWarning                 // worth knowing, won't block
	SeverityInfo                    // informational
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	default:
		return "info"
	}
}

// Issue is a single finding from the validator.
type Issue struct {
	Severity Severity
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Severity, i.Message)
}

// Report holds all issues found during validation.
type Report struct {
	issues []Issue
}

// Valid returns true if there are no error-severity issues.
func (r *Report) Valid() bool {
	for _, i := range r.issues {
		if i.Severity == SeverityError {
			return false
		}
	}
	return true
}

// Errors returns only error-severity issues.
func (r *Report) Errors() []Issue {
	var out []Issue
	for _, i := range r.issues {
		if i.Severity == SeverityError {
			out = append(out, i)
		}
	}
	return out
}

// Warnings returns only warning-severity issues.
func (r *Report) Warnings() []Issue {
	var out []Issue
	for _, i := range r.issues {
		if i.Severity == SeverityWarning {
			out = append(out, i)
		}
	}
	return out
}

// String returns a human-readable summary of all issues.
func (r *Report) String() string {
	if len(r.issues) == 0 {
		return "validate: no issues found"
	}
	var b strings.Builder
	for _, i := range r.issues {
		b.WriteString(i.String())
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func (r *Report) add(sev Severity, format string, args ...any) {
	r.issues = append(r.issues, Issue{
		Severity: sev,
		Message:  fmt.Sprintf(format, args...),
	})
}

// Document renders a doc.Document and validates the resulting SVG.
func Document(d *doc.Document) (*Report, error) {
	svg, err := d.Render()
	if err != nil {
		return nil, fmt.Errorf("validate.Document: render failed: %w", err)
	}
	return SVG(svg)
}

// SVG validates a raw SVG string.
func SVG(svg string) (*Report, error) {
	r := &Report{}

	// ── 1. Well-formed XML ────────────────────────────────────────────────
	dec := xml.NewDecoder(strings.NewReader(svg))
	ids := map[string]int{}        // id → count
	urlRefs := []string{}          // url(#...) references found
	definedIDs := map[string]bool{} // all defined ids

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			// Check SVG namespace
			if t.Name.Local == "svg" {
				hasNS := false
				hasViewBox := false
				for _, a := range t.Attr {
					if a.Name.Local == "xmlns" {
						hasNS = true
					}
					if a.Name.Local == "viewBox" {
						hasViewBox = true
					}
				}
				if !hasNS {
					r.add(SeverityWarning, "root <svg> is missing xmlns attribute")
				}
				if !hasViewBox {
					r.add(SeverityWarning, "root <svg> is missing viewBox attribute")
				}
			}

			for _, a := range t.Attr {
				// Collect all id definitions
				if a.Name.Local == "id" {
					ids[a.Value]++
					definedIDs[a.Value] = true
				}
				// Collect url(#...) references
				if strings.Contains(a.Value, "url(#") {
					ref := extractURLRef(a.Value)
					if ref != "" {
						urlRefs = append(urlRefs, ref)
					}
				}
				// Zero-dimension check
				if t.Name.Local == "rect" || t.Name.Local == "ellipse" {
					if (a.Name.Local == "width" || a.Name.Local == "height" ||
						a.Name.Local == "r" || a.Name.Local == "rx" || a.Name.Local == "ry") &&
						a.Value == "0" {
						r.add(SeverityWarning, "<%s> has zero %s — element will be invisible",
							t.Name.Local, a.Name.Local)
					}
				}
			}
		}
	}

	// ── 2. Well-formed XML second pass (syntax check) ─────────────────────
	dec2 := xml.NewDecoder(strings.NewReader(svg))
	for {
		_, err := dec2.Token()
		if err != nil {
			if err.Error() != "EOF" {
				r.add(SeverityError, "SVG is not well-formed XML: %s", err)
			}
			break
		}
	}

	// ── 3. Duplicate IDs ──────────────────────────────────────────────────
	for id, count := range ids {
		if count > 1 {
			r.add(SeverityError, "duplicate id %q appears %d times", id, count)
		}
	}

	// ── 4. Unresolved url(#...) references ───────────────────────────────
	for _, ref := range urlRefs {
		if !definedIDs[ref] {
			r.add(SeverityError, "url(#%s) references an id that does not exist", ref)
		}
	}

	return r, nil
}

// extractURLRef pulls the id out of a url(#id) value.
func extractURLRef(s string) string {
	start := strings.Index(s, "url(#")
	if start == -1 {
		return ""
	}
	start += 5
	end := strings.Index(s[start:], ")")
	if end == -1 {
		return ""
	}
	return s[start : start+end]
}
