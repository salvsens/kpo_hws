package main

import (
	"testing"

	"example.com/vyshkat/internal/cli"
	"example.com/vyshkat/internal/fleet"
)

// TestBuildContainer проверяет граф зависимостей dig (ADR-0003): ошибки регистрации
// видны только при запуске, поэтому тест запрашивает каждую зависимость.
func TestBuildContainer(t *testing.T) {
	t.Parallel()

	c, err := buildContainer()
	if err != nil {
		t.Fatalf("buildContainer() error = %v", err)
	}

	var (
		repl *cli.REPL
		sc   fleet.ServiceCenter
	)
	err = c.Invoke(func(r *cli.REPL, s fleet.ServiceCenter, _ fleet.BeginnerPolicy, _ *fleet.Fleet) {
		repl, sc = r, s
	})
	if err != nil {
		t.Fatalf("Invoke() error = %v", err)
	}
	if repl == nil {
		t.Error("контейнер вернул nil вместо *cli.REPL")
	}
	if _, ok := sc.(fleet.EnergyLimitInspection); !ok {
		t.Errorf("ServiceCenter = %T, ожидался fleet.EnergyLimitInspection", sc)
	}
}
