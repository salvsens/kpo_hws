// Точка входа приложения «ВышКат»: сборка зависимостей в контейнере dig и запуск REPL.
package main

import (
	"fmt"
	"os"

	"go.uber.org/dig"

	"example.com/vyshkat/internal/cli"
	"example.com/vyshkat/internal/fleet"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

// run собирает контейнер и запускает REPL на стандартных потоках.
func run() error {
	c, err := buildContainer()
	if err != nil {
		return fmt.Errorf("сборка зависимостей: %w", err)
	}
	return c.Invoke(func(r *cli.REPL) error {
		return r.Run(os.Stdin, os.Stdout)
	})
}

// buildContainer регистрирует конструкторы всех зависимостей приложения.
// Вынесена отдельно, чтобы тест проверил, что граф зависимостей собирается.
func buildContainer() (*dig.Container, error) {
	c := dig.New()
	providers := []any{
		newServiceCenter,
		newBeginnerPolicy,
		fleet.New,
		cli.New,
	}
	for _, p := range providers {
		if err := c.Provide(p); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// newServiceCenter задаёт техосмотр приложения.
// Возвращает интерфейс: dig связывает зависимости по типу, а fleet.New ждёт fleet.ServiceCenter.
func newServiceCenter() fleet.ServiceCenter {
	return fleet.EnergyLimitInspection{MaxDailyKWh: fleet.DefaultMaxDailyKWh}
}

// newBeginnerPolicy задаёт правило «для новичков» приложения.
func newBeginnerPolicy() fleet.BeginnerPolicy {
	return fleet.MinSimplicityPolicy{Min: fleet.DefaultMinSimplicity}
}
