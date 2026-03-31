package cli

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/DraftOps1/tf-fastpath/internal/config"
)

func TestRunWithoutArgsPrintsUsage(t *testing.T) {
	t.Parallel()

	var stdout strings.Builder
	app := New(&stdout, &strings.Builder{}, config.Config{})

	err := app.Run(context.Background(), nil)
	if !errors.Is(err, ErrHelpRequested) {
		t.Fatalf("error = %v, want ErrHelpRequested", err)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunCommandHelp(t *testing.T) {
	t.Parallel()

	var stdout strings.Builder
	app := New(&stdout, &strings.Builder{}, config.Config{})

	err := app.Run(context.Background(), []string{"preview", "--help"})
	if !errors.Is(err, ErrHelpRequested) {
		t.Fatalf("error = %v, want ErrHelpRequested", err)
	}
	if !strings.Contains(stdout.String(), "usually sufficient for merge decisions") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunCommandScaffold(t *testing.T) {
	t.Parallel()

	var stdout strings.Builder
	app := New(&stdout, &strings.Builder{}, config.Config{})

	err := app.Run(context.Background(), []string{"gate", "--print-config"})
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("error = %v, want ErrNotImplemented", err)
	}
	for _, fragment := range []string{"WorkingDir:", "Command: gate", "Status: scaffold only"} {
		if !strings.Contains(stdout.String(), fragment) {
			t.Fatalf("stdout missing %q in %q", fragment, stdout.String())
		}
	}
}
