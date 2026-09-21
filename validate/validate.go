// Package validate provides tools to audit SVG documents for structural issues,
// duplicate IDs, unresolved references, and other common problems.
package validate

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/dominionthedev/finite/doc"
)

// Severity indicates the importance of a validation issue.
type Severity int

const (
	SeverityError   Severity = iota // Blocks correct rendering or violates SVG rules
	SeverityWarning                 // Likely problem (visibility, missing attrs)
	SeverityInfo                    // Observation worth knowing, not necessarily wrong
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

// Issue is a single validation finding.
type Issue struct {
	Severity Severity
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Severity, i.Message)
}

// Report collects issues from a validation run.
type Report struct {
	issues []Issue
}

// Valid reports whether the document has no error-severity issues.
func (r *Report) Valid() bool {
	for _, i := range r.issues {
		if i.Severity == SeverityError {
			return false
		}
	}
	return true
}

// Empty reports whether any issues were recorded (any severity).
func (r *Report) Empty() bool {
	return len(r.issues) == 0
}

// Errors returns error-severity issues.
func (r *Report) Errors() []Issue {
	return r.filter(SeverityError)
}

// Warnings returns warning-severity issues.
func (r *Report) Warnings() []Issue {
	return r.filter(SeverityWarning)
}

// Infos returns info-severity issues.
func (r *Report) Infos() []Issue {
	return r.filter(SeverityInfo)
}

// Issues returns all findings in the order they were recorded.
func (r *Report) Issues() []Issue {
	out := make([]Issue, len(r.issues))
	copy(out, r.issues)
	return out
}

func (r *Report) filter(sev Severity) []Issue {
	var out []Issue
	for _, i := range r.issues {
		if i.Severity == sev {
			out = append(out, i)
		}
	}
	return out
}

// String returns a multi-line summary of all issues.
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

// Document validates a Finite document by rendering it and auditing the SVG.
func Document(d *doc.Document) (*Report, error) {
	svg, err := d.Render()
	if err != nil {
		return nil, fmt.Errorf("validate.Document: render failed: %w", err)
	}
	return SVG(svg)
}

// idHit tracks where an id was defined (element local name).
type idHit struct {
	tag string
}

// SVG validates a raw SVG XML string.
func SVG(svg string) (*Report, error) {
	r := &Report{}

	if strings.TrimSpace(svg) == "" {
		r.add(SeverityError, "SVG is empty")
		return r, nil
	}

	ids := map[string][]idHit{} // id → definition sites
	urlRefs := map[string]int{} // ref → use count
	elementCount := 0
	groupCount := 0
	emptyGroups := 0
	depth := 0
	groupDepthStart := map[int]bool{} // depth → group opened with no children yet

	dec := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			elementCount++
			depth++

			if t.Name.Local == "g" {
				groupCount++
				groupDepthStart[depth] = true
			} else {
				// Non-group child marks parent groups as non-empty
				for d := range groupDepthStart {
					if d < depth {
						groupDepthStart[d] = false
					}
				}
			}

			if t.Name.Local == "svg" && depth == 1 {
				hasNS, hasViewBox := false, false
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
				if a.Name.Local == "id" && a.Value != "" {
					ids[a.Value] = append(ids[a.Value], idHit{tag: t.Name.Local})
				}
				if strings.Contains(a.Value, "url(#") {
					if ref := extractURLRef(a.Value); ref != "" {
						urlRefs[ref]++
					}
				}
				if t.Name.Local == "rect" || t.Name.Local == "ellipse" || t.Name.Local == "circle" {
					if isZeroDimAttr(a.Name.Local) && isZeroValue(a.Value) {
						r.add(SeverityWarning, "<%s> has zero %s — element will be invisible",
							t.Name.Local, a.Name.Local)
					}
				}
			}

		case xml.EndElement:
			if t.Name.Local == "g" {
				if groupDepthStart[depth] {
					emptyGroups++
				}
				delete(groupDepthStart, depth)
			}
			depth--
		}
	}

	// Well-formed XML check
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

	// Duplicate IDs — include element names to make layer/child clashes obvious
	dupIDs := make([]string, 0)
	for id, hits := range ids {
		if len(hits) > 1 {
			dupIDs = append(dupIDs, id)
		}
	}
	sort.Strings(dupIDs)
	for _, id := range dupIDs {
		hits := ids[id]
		tags := make([]string, len(hits))
		for i, h := range hits {
			tags[i] = "<" + h.tag + ">"
		}
		r.add(SeverityError,
			"duplicate id %q appears %d times on %s — layer names become group ids; child GroupID values must be unique",
			id, len(hits), strings.Join(tags, ", "))
	}

	// Unresolved references
	refNames := make([]string, 0, len(urlRefs))
	for ref := range urlRefs {
		refNames = append(refNames, ref)
	}
	sort.Strings(refNames)
	for _, ref := range refNames {
		if _, ok := ids[ref]; !ok {
			r.add(SeverityError, "url(#%s) references an id that does not exist", ref)
		}
	}

	// Unused definitions (info): ids never referenced via url(#...)
	// Skip structural ids that are normal for layers/groups (only flag paint servers / filters-ish tags)
	unused := make([]string, 0)
	for id, hits := range ids {
		if urlRefs[id] > 0 {
			continue
		}
		for _, h := range hits {
			if isPaintServerTag(h.tag) {
				unused = append(unused, id)
				break
			}
		}
	}
	sort.Strings(unused)
	for _, id := range unused {
		r.add(SeverityInfo, "id %q defines a paint server or effect that is never referenced via url(#...)", id)
	}

	if emptyGroups > 0 {
		r.add(SeverityInfo, "%d empty <g> group(s) — no painted children", emptyGroups)
	}

	// Summary info when the document is non-trivial
	if elementCount > 0 {
		r.add(SeverityInfo, "document has %d element(s), %d group(s), %d unique id(s)",
			elementCount, groupCount, len(ids))
	}

	return r, nil
}

func isZeroDimAttr(name string) bool {
	switch name {
	case "width", "height", "r", "rx", "ry":
		return true
	default:
		return false
	}
}

func isZeroValue(v string) bool {
	v = strings.TrimSpace(v)
	return v == "0" || v == "0.0" || v == "0.00" || v == "0.000" || v == "0.0000"
}

func isPaintServerTag(tag string) bool {
	switch tag {
	case "linearGradient", "radialGradient", "pattern", "filter", "clipPath", "mask", "marker", "symbol":
		return true
	default:
		return false
	}
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
