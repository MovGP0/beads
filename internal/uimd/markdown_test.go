package uimd

import (
	"os"
	"strings"
	"testing"
)

func TestRenderMarkdownStripsEscapesWhenANSIUnsupported(t *testing.T) {
	withMarkdownEnv(t, map[string]string{
		"NO_COLOR":        "1",
		"TERM":            "dumb",
		"CLICOLOR_FORCE":  "",
		"FORCE_HYPERLINK": "",
		"BD_AGENT_MODE":   "",
		"CLAUDE_CODE":     "",
	})

	out := RenderMarkdown("# Title\n\n[link](https://example.com)\n\n| A | B |\n| - | - |\n| 1 | 2 |\n")
	if strings.Contains(out, "\x1b[") || strings.Contains(out, "\x1b]8;") {
		t.Fatalf("expected no terminal escapes when ANSI is unsupported, got %q", out)
	}
	if !strings.Contains(out, "Title") || !strings.Contains(out, "example.com") {
		t.Fatalf("expected rendered markdown content, got %q", out)
	}
}

func TestRenderMarkdownStripsOSC8WhenHyperlinksUnsupported(t *testing.T) {
	withMarkdownEnv(t, map[string]string{
		"NO_COLOR":        "",
		"TERM":            "xterm-256color",
		"CLICOLOR_FORCE":  "1",
		"FORCE_HYPERLINK": "",
		"BD_AGENT_MODE":   "",
		"CLAUDE_CODE":     "",
	})

	out := RenderMarkdown("[link](https://example.com)")
	if strings.Contains(out, "\x1b]8;") {
		t.Fatalf("expected OSC 8 hyperlinks to be stripped, got %q", out)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI styling when color is forced, got %q", out)
	}
}

func TestRenderMarkdownKeepsOSC8WhenHyperlinksSupported(t *testing.T) {
	withMarkdownEnv(t, map[string]string{
		"NO_COLOR":        "",
		"TERM":            "xterm-256color",
		"CLICOLOR_FORCE":  "1",
		"FORCE_HYPERLINK": "1",
		"BD_AGENT_MODE":   "",
		"CLAUDE_CODE":     "",
	})

	out := RenderMarkdown("[link](https://example.com)")
	if !strings.Contains(out, "\x1b]8;") {
		t.Fatalf("expected OSC 8 hyperlink escapes, got %q", out)
	}
}

func TestRenderMarkdownCanKeepOSC8WithoutANSIColor(t *testing.T) {
	withMarkdownEnv(t, map[string]string{
		"NO_COLOR":        "",
		"TERM":            "xterm-256color",
		"CLICOLOR_FORCE":  "",
		"FORCE_HYPERLINK": "1",
		"BD_AGENT_MODE":   "",
		"CLAUDE_CODE":     "",
	})

	out := RenderMarkdown("[link](https://example.com)")
	if !strings.Contains(out, "\x1b]8;") {
		t.Fatalf("expected OSC 8 hyperlink escapes, got %q", out)
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("expected no ANSI SGR styling, got %q", out)
	}
}

func TestRenderMarkdownReturnsRawMarkdownInAgentMode(t *testing.T) {
	withMarkdownEnv(t, map[string]string{
		"NO_COLOR":        "",
		"TERM":            "xterm-256color",
		"CLICOLOR_FORCE":  "1",
		"FORCE_HYPERLINK": "1",
		"BD_AGENT_MODE":   "1",
		"CLAUDE_CODE":     "",
	})

	input := "# Title\n\n[link](https://example.com)"
	if out := RenderMarkdown(input); out != input {
		t.Fatalf("agent mode should return raw markdown, got %q", out)
	}
}

func withMarkdownEnv(t *testing.T, values map[string]string) {
	t.Helper()

	keys := []string{
		"BD_GIT_HOOK",
		"NO_COLOR",
		"CLICOLOR",
		"CLICOLOR_FORCE",
		"FORCE_HYPERLINK",
		"TERM",
		"TERM_PROGRAM",
		"WT_SESSION",
		"KITTY_WINDOW_ID",
		"WEZTERM_EXECUTABLE",
		"KONSOLE_VERSION",
		"DOMTERM",
		"GHOSTTY_RESOURCES_DIR",
		"VTE_VERSION",
		"BD_AGENT_MODE",
		"CLAUDE_CODE",
	}
	orig := make(map[string]string, len(keys))
	for _, key := range keys {
		orig[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	t.Cleanup(func() {
		for _, key := range keys {
			if orig[key] == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, orig[key])
			}
		}
	})

	for key, value := range values {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}
