package draw

import (
	"fmt"
	"strings"
)

// xmlEscape escapes text for use inside SVG text nodes and attributes.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// a11yContent returns <title> / <desc> child markup when set.
func a11yContent(title, desc string) string {
	var b strings.Builder
	if title != "" {
		b.WriteString("<title>")
		b.WriteString(xmlEscape(title))
		b.WriteString("</title>")
	}
	if desc != "" {
		b.WriteString("<desc>")
		b.WriteString(xmlEscape(desc))
		b.WriteString("</desc>")
	}
	return b.String()
}

// a11yAttrs returns role and aria-label attribute fragments.
func a11yAttrs(role, ariaLabel string) string {
	var b strings.Builder
	if role != "" {
		b.WriteString(fmt.Sprintf(` role="%s"`, xmlEscape(role)))
	}
	if ariaLabel != "" {
		b.WriteString(fmt.Sprintf(` aria-label="%s"`, xmlEscape(ariaLabel)))
	}
	return b.String()
}
