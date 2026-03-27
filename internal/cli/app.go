package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"

	"github.com/DraftOps1/tf-fastpath/internal/config"
)

var (
	// ErrHelpRequested reports that help text was printed successfully.
	ErrHelpRequested = errors.New("help requested")
	// ErrNotImplemented reports that the command surface exists but the implementation is pending.
	ErrNotImplemented = errors.New("command scaffold only")
)

type App struct {
	stdout io.Writer
	stderr io.Writer
	config config.Config
}

type command struct {
	name        string
	summary     string
	description string
}

type notImplementedError struct {
	command string
}

func (e notImplementedError) Error() string {
	return fmt.Sprintf("%s: scaffold only; implementation not yet available", e.command)
}

func (e notImplementedError) Unwrap() error {
	return ErrNotImplemented
}

func New(stdout io.Writer, stderr io.Writer, cfg config.Config) *App {
	return &App{stdout: stdout, stderr: stderr, config: cfg}
}

func (a *App) Run(ctx context.Context, args []string) error {
	_ = ctx

	if len(args) == 0 {
		a.printUsage()
		return ErrHelpRequested
	}

	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		a.printUsage()
		return ErrHelpRequested
	}

	cmd, ok := a.lookupCommand(args[0])
	if !ok {
		a.printUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}

	return a.runCommand(cmd, args[1:])
}

func (a *App) runCommand(cmd command, args []string) error {
	fs := flag.NewFlagSet(cmd.name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var printConfig bool
	fs.BoolVar(&printConfig, "print-config", false, "print resolved shared configuration")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			a.printCommandUsage(cmd)
			return ErrHelpRequested
		}
		return fmt.Errorf("parse %s flags: %w", cmd.name, err)
	}

	if printConfig {
		a.printConfig()
	}

	fmt.Fprintf(a.stdout, "Command: %s\n", cmd.name)
	fmt.Fprintf(a.stdout, "Purpose: %s\n", cmd.summary)
	fmt.Fprintf(a.stdout, "Status: scaffold only\n")

	return notImplementedError{command: cmd.name}
}

func (a *App) printUsage() {
	fmt.Fprintln(a.stdout, "Usage:")
	fmt.Fprintln(a.stdout, "  tf-fastpath <command> [flags]")
	fmt.Fprintln(a.stdout)
	fmt.Fprintln(a.stdout, "Commands:")
	for _, cmd := range commands() {
		fmt.Fprintf(a.stdout, "  %-8s %s\n", cmd.name, cmd.summary)
	}
	fmt.Fprintln(a.stdout)
	fmt.Fprintln(a.stdout, "Shared configuration (environment variables):")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_WORKING_DIR      root directory for Terraform/OpenTofu execution")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_DATA_DIR         directory for derived planning data")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_SQLITE_PATH      SQLite file path for derived data")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_TERRAFORM_BIN    terraform executable path")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_OPENTOFU_BIN     tofu executable path")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_GIT_BIN          git executable path")
	fmt.Fprintln(a.stdout, "  TFFASTPATH_DEFAULT_RUNTIME  terraform or opentofu")
}

func (a *App) printCommandUsage(cmd command) {
	fmt.Fprintf(a.stdout, "Usage:\n  tf-fastpath %s [--print-config]\n\n", cmd.name)
	fmt.Fprintf(a.stdout, "%s\n\n", cmd.description)
	fmt.Fprintf(a.stdout, "Purpose: %s\n", cmd.summary)
}

func (a *App) printConfig() {
	fmt.Fprintf(a.stdout, "WorkingDir: %s\n", a.config.WorkingDir)
	fmt.Fprintf(a.stdout, "DataDir: %s\n", a.config.DataDir)
	fmt.Fprintf(a.stdout, "SQLitePath: %s\n", a.config.SQLitePath)
	fmt.Fprintf(a.stdout, "TerraformBin: %s\n", a.config.TerraformBin)
	fmt.Fprintf(a.stdout, "OpenTofuBin: %s\n", a.config.OpenTofuBin)
	fmt.Fprintf(a.stdout, "GitBin: %s\n", a.config.GitBin)
	fmt.Fprintf(a.stdout, "DefaultRuntime: %s\n", a.config.DefaultRuntime)
	fmt.Fprintln(a.stdout)
}

func (a *App) lookupCommand(name string) (command, bool) {
	for _, cmd := range commands() {
		if cmd.name == name {
			return cmd, true
		}
	}
	return command{}, false
}

func commands() []command {
	items := []command{
		{
			name:        "index",
			summary:     "build derived planning data from state, provider schema, and HCL",
			description: "Create the derived indexes and metadata that later preview and gate steps will reuse.",
		},
		{
			name:        "preview",
			summary:     "return a fast, non-authoritative preview from local changes",
			description: "Estimate affected addresses and render a fast preview without replacing the authoritative full plan.",
		},
		{
			name:        "verify",
			summary:     "record drift and preview confidence signals",
			description: "Run verification-oriented checks such as refresh-only planning and drift capture.",
		},
		{
			name:        "gate",
			summary:     "run the authoritative plan before merge or apply",
			description: "Execute the final gate that remains authoritative for merge and apply decisions.",
		},
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].name < items[j].name
	})

	return items
}
