package main

import (
	"context"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Parallel()

	var stdout strings.Builder
	var stderr strings.Builder

	exitCode := run(
		context.Background(),
		[]string{"--help"},
		&stdout,
		&stderr,
		func(string) string { return "" },
		func() (string, error) { return "/repo", nil },
	)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}

	got := stdout.String()
	for _, fragment := range []string{"Usage:", "index", "preview", "verify", "gate"} {
		if !strings.Contains(got, fragment) {
			t.Fatalf("help output missing %q in %q", fragment, got)
		}
	}
}

func TestRunCommandScaffold(t *testing.T) {
	t.Parallel()

	var stdout strings.Builder
	var stderr strings.Builder

	exitCode := run(
		context.Background(),
		[]string{"preview"},
		&stdout,
		&stderr,
		func(string) string { return "" },
		func() (string, error) { return "/repo", nil },
	)

	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}

	if !strings.Contains(stdout.String(), "Command: preview") {
		t.Fatalf("stdout = %q, want scaffold summary", stdout.String())
	}

	if !strings.Contains(stderr.String(), "preview: scaffold only") {
		t.Fatalf("stderr = %q, want scaffold error", stderr.String())
	}
}
