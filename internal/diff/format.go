package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Render writes a human-readable diff report to w.
func Render(w io.Writer, r Result, envA, envB string, color bool) {
	fmt.Fprintf(w, "Comparing secrets: %s vs %s\n", envA, envB)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, k := range r.OnlyInA {
		line := fmt.Sprintf("  only in %s: %s", envA, k)
		if color {
			line = colorGreen + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	for _, k := range r.OnlyInB {
		line := fmt.Sprintf("  only in %s: %s", envB, k)
		if color {
			line = colorRed + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	for _, d := range r.Differing {
		line := fmt.Sprintf("  ~ %s: %q (%s) vs %q (%s)", d.Key, d.ValueA, envA, d.ValueB, envB)
		if color {
			line = colorYellow + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	fmt.Fprintf(w, "\nSummary: %d matching, %d differing, %d only-in-%s, %d only-in-%s\n",
		len(r.Matching), len(r.Differing), len(r.OnlyInA), envA, len(r.OnlyInB), envB)
}
