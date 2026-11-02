package utils

import (
	"strings"
	"testing"

	"charm.land/glamour/v2"
	"github.com/charmbracelet/x/ansi"
)

const linkURL = "https://github.com/o/r/blob/abc/x/seal.mjs#L147"

func renderLinks(t *testing.T, env []string, md string) string {
	t.Helper()
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
		GlamourHyperlinks(env),
	)
	if err != nil {
		t.Fatal(err)
	}
	out, err := r.Render(md)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestGlamourHyperlinks(t *testing.T) {
	md := "See [seal.mjs:147][c1].\n\n[c1]: " + linkURL + "\n"
	osc8 := ";" + linkURL + "\a" // tail of the OSC 8 open sequence

	t.Run("OSC 8 terminal prints only the text", func(t *testing.T) {
		out := renderLinks(t, []string{"TERM_PROGRAM=ghostty"}, md)
		if !strings.Contains(out, "\x1b]8;") || !strings.Contains(out, osc8) {
			t.Fatalf("text is not an OSC 8 link: %q", out)
		}
		if got := ansi.Strip(out); strings.Contains(got, linkURL) {
			t.Fatalf("URL printed after text: %q", got)
		}
		if got := ansi.Strip(out); !strings.Contains(got, "seal.mjs:147") {
			t.Fatalf("text missing: %q", got)
		}
	})

	t.Run("other terminal prints the text and the URL", func(t *testing.T) {
		out := ansi.Strip(renderLinks(t, nil, md))
		if !strings.Contains(out, "seal.mjs:147 "+linkURL) {
			t.Fatalf("URL missing: %q", out)
		}
	})

	t.Run("autolink still prints the URL", func(t *testing.T) {
		out := ansi.Strip(renderLinks(t, []string{"TERM_PROGRAM=ghostty"}, "<"+linkURL+">\n"))
		if !strings.Contains(out, linkURL) {
			t.Fatalf("autolink URL missing: %q", out)
		}
	})
}

func TestSupportsHyperlinks(t *testing.T) {
	for _, tc := range []struct {
		env  []string
		want bool
	}{
		{[]string{"TERM_PROGRAM=ghostty"}, true},
		{[]string{"TERM=xterm-ghostty"}, true},
		{[]string{"TERM=xterm-kitty"}, true},
		{[]string{"TERM_PROGRAM=WezTerm"}, true},
		{[]string{"TERM_PROGRAM=iTerm.app"}, true},
		{[]string{"VTE_VERSION=7600"}, true},
		{[]string{"TERM_PROGRAM=Apple_Terminal"}, false},
		{[]string{"TERM=xterm-256color"}, false},
		{nil, false},
	} {
		if got := SupportsHyperlinks(tc.env); got != tc.want {
			t.Errorf("SupportsHyperlinks(%v) = %v, want %v", tc.env, got, tc.want)
		}
	}
}
