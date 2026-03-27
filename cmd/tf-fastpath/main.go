package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/DraftOps1/tf-fastpath/internal/cli"
	"github.com/DraftOps1/tf-fastpath/internal/config"
)

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr, os.Getenv, os.Getwd))
}

func run(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	lookupEnv func(string) string,
	getwd func() (string, error),
) int {
	cfg, err := config.Load(lookupEnv, getwd)
	if err != nil {
		fmt.Fprintf(stderr, "tf-fastpath: load config: %v\n", err)
		return 1
	}

	app := cli.New(stdout, stderr, cfg)
	if err := app.Run(ctx, args); err != nil {
		switch {
		case errors.Is(err, cli.ErrHelpRequested):
			return 0
		case errors.Is(err, cli.ErrNotImplemented):
			fmt.Fprintf(stderr, "tf-fastpath: %v\n", err)
			return 2
		default:
			fmt.Fprintf(stderr, "tf-fastpath: %v\n", err)
			return 1
		}
	}

	return 0
}
