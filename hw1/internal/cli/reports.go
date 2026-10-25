package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// columnPadding — отступ между колонками таблиц.
const columnPadding = 2

// printReport печатает число единиц и суточное энергопотребление (FR-6).
func (r *REPL) printReport(out io.Writer) {
	rep := r.fleet.Report()
	fmt.Fprintf(out, "Всего: %d (транспорт: %d, вещи: %d)\n", rep.Total, rep.Vehicles, rep.Items)
	fmt.Fprintf(out, "Суточное энергопотребление: %.2f кВт·ч\n", rep.DailyKWh)
}

// printBeginners печатает транспорт, подходящий для новичков (FR-7).
func (r *REPL) printBeginners(out io.Writer) error {
	vehicles := r.fleet.Beginners()
	if len(vehicles) == 0 {
		fmt.Fprintln(out, "Транспорта для новичков нет.")
		return nil
	}
	w := newTable(out)
	for _, v := range vehicles {
		fmt.Fprintf(w, "%s\t%s\t%s\tпростота %d\n", v.ID(), v.Kind().Title, v.Name(), v.BeginnerSimplicity())
	}
	return w.Flush()
}

// printInventory печатает все объекты на балансе в порядке постановки (FR-8).
func (r *REPL) printInventory(out io.Writer) error {
	assets := r.fleet.Inventory()
	if len(assets) == 0 {
		fmt.Fprintln(out, "На балансе пока ничего нет.")
		return nil
	}
	w := newTable(out)
	for _, a := range assets {
		fmt.Fprintf(w, "%s\t%s\t%s\n", a.ID(), a.Kind().Title, a.Name())
	}
	return w.Flush()
}

// newTable создаёт писатель, выравнивающий колонки, разделённые табуляцией.
func newTable(out io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(out, 0, 0, columnPadding, ' ', 0)
}
