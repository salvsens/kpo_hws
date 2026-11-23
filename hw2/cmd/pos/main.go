// Package main connects the gradebook and starts the console.
package main

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/dig"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/cli"
	"github.com/salvsens/kpo_hws/hw2/internal/memory"
)

// main starts the gradebook on the standard input and output streams.
func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		// The program exits if it cannot report the error.
		_, _ = fmt.Fprintln(os.Stderr, "Ошибка запуска или ввода-вывода:", err)
		os.Exit(1)
	}
}

// run resolves the console and connects it to the streams.
func run(in io.Reader, out io.Writer) error {
	container, err := buildContainer()
	if err != nil {
		return err
	}
	return container.Invoke(func(repl *cli.REPL) error {
		return repl.Run(in, out)
	})
}

// buildContainer registers the constructors for the console.
func buildContainer() (*dig.Container, error) {
	container := dig.New()
	if err := container.Provide(memory.NewRepository, dig.As((*app.Repository)(nil))); err != nil {
		return nil, fmt.Errorf("register repository: %w", err)
	}
	if err := container.Provide(app.NewGradebook); err != nil {
		return nil, fmt.Errorf("register gradebook: %w", err)
	}
	if err := container.Provide(cli.New); err != nil {
		return nil, fmt.Errorf("register console: %w", err)
	}
	return container, nil
}
