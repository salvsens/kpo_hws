package cli

import (
	"fmt"
	"io"
	"strconv"
)

// Число аргументов команд без кВт·ч; электрическим видам нужен ещё один.
const (
	vehicleArgs = 3 // <вид> <название> <простота>
	itemArgs    = 2 // <вид> <название>
)

// addVehicle разбирает аргументы, создаёт транспорт и передаёт его парку на техосмотр.
func (r *REPL) addVehicle(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: add-vehicle <вид> <название> <простота> [кВт·ч]", errUsage)
	}
	spec, ok := vehicleSpecs[args[0]]
	if !ok {
		return fmt.Errorf("%w %q, доступны: %s", errUnknownKind, args[0], kindNames(vehicleSpecs))
	}
	want := vehicleArgs
	if spec.electric {
		want++
	}
	if len(args) != want {
		return fmt.Errorf("%w: %s", errUsage, vehicleUsage(args[0], spec.electric))
	}

	simplicity, err := parseSimplicity(args[2])
	if err != nil {
		return err
	}
	var kwh float64
	if spec.electric {
		if kwh, err = parseKWh(args[3]); err != nil {
			return err
		}
	}

	v, err := spec.build(r.fleet.NextID(spec.kind), args[1], simplicity, kwh)
	if err != nil {
		return err
	}
	if err = r.fleet.AddVehicle(v); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "Техосмотр пройден. Принят: %s, инв. № %s\n", v.Name(), v.ID())
	return nil
}

// addItem разбирает аргументы, создаёт вещь и ставит её на баланс.
func (r *REPL) addItem(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: add-item <вид> <название> [кВт·ч]", errUsage)
	}
	spec, ok := itemSpecs[args[0]]
	if !ok {
		return fmt.Errorf("%w %q, доступны: %s", errUnknownKind, args[0], kindNames(itemSpecs))
	}
	want := itemArgs
	if spec.electric {
		want++
	}
	if len(args) != want {
		return fmt.Errorf("%w: %s", errUsage, itemUsage(args[0], spec.electric))
	}

	var kwh float64
	if spec.electric {
		var err error
		if kwh, err = parseKWh(args[2]); err != nil {
			return err
		}
	}

	a, err := spec.build(r.fleet.NextID(spec.kind), args[1], kwh)
	if err != nil {
		return err
	}
	if err = r.fleet.AddItem(a); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "На баланс поставлен: %s, инв. № %s\n", a.Name(), a.ID())
	return nil
}

// vehicleUsage возвращает подсказку по аргументам для конкретного вида транспорта.
func vehicleUsage(kind string, electric bool) string {
	if electric {
		return "add-vehicle " + kind + " <название> <простота> <кВт·ч>"
	}
	return "add-vehicle " + kind + " <название> <простота>"
}

// itemUsage возвращает подсказку по аргументам для конкретного вида вещи.
func itemUsage(kind string, electric bool) string {
	if electric {
		return "add-item " + kind + " <название> <кВт·ч>"
	}
	return "add-item " + kind + " <название>"
}

// parseSimplicity читает простоту; диапазон проверяет домен.
func parseSimplicity(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: простота %q — не целое число", errBadNumber, s)
	}
	return n, nil
}

// parseKWh читает суточный расход; знак проверяет домен.
func parseKWh(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: расход %q — не число", errBadNumber, s)
	}
	return f, nil
}
