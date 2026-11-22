package fleet_test

import (
	"errors"
	"math"
	"slices"
	"strings"
	"testing"

	"example.com/vyshkat/internal/domain"
	"example.com/vyshkat/internal/fleet"
)

// fakeServiceCenter — подмена техосмотра через DI: возвращает заданное решение
// и запоминает номера транспорта, который ему передали на проверку.
type fakeServiceCenter struct {
	verdict   fleet.Verdict
	inspected []string
}

func (f *fakeServiceCenter) Inspect(v domain.Vehicle) fleet.Verdict {
	f.inspected = append(f.inspected, v.ID())
	return f.verdict
}

// allowByName — подмена правила для новичков: разрешает транспорт по названию.
type allowByName map[string]bool

func (a allowByName) Allows(v domain.Vehicle) bool { return a[v.Name()] }

var (
	accept = fleet.Verdict{Accepted: true}
	reject = fleet.Verdict{Reason: "тормоза неисправны"}
)

// defaultPolicy — правило для новичков по умолчанию, если тест его не проверяет.
var defaultPolicy = fleet.MinSimplicityPolicy{Min: fleet.DefaultMinSimplicity}

// must останавливает тест, если конструктор вернул ошибку: must(domain.NewHelmet(...))(t).
// t передаётся вторым вызовом: Go не разрешает передать f() с двумя результатами вместе с другим аргументом.
func must[T any](v T, err error) func(*testing.T) T {
	return func(t *testing.T) T {
		t.Helper()
		if err != nil {
			t.Fatalf("подготовка данных: %v", err)
		}
		return v
	}
}

// newScooter создаёт электросамокат с номером, который парк выдал бы сейчас.
func newScooter(t *testing.T, f *fleet.Fleet, name string, simplicity int, kwh float64) *domain.ElectricScooter {
	t.Helper()
	s, err := domain.NewElectricScooter(f.NextID(domain.KindElectricScooter), name, simplicity, kwh)
	if err != nil {
		t.Fatalf("NewElectricScooter() error = %v", err)
	}
	return s
}

// TestAddVehicleUsesServiceCenter — главный тест подмены техосмотра через DI:
// парк получает фейк в конструкторе и выполняет его решение — и приём, и отказ.
func TestAddVehicleUsesServiceCenter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		verdict     fleet.Verdict
		wantErr     error
		wantReason  string // причина отказа должна дойти до вызывающего кода
		wantOnBoard int
		wantNextID  string
	}{
		{name: "техосмотр принял", verdict: accept, wantOnBoard: 1, wantNextID: "VK-SC-0002"},
		{
			name: "техосмотр отклонил", verdict: reject,
			wantErr: fleet.ErrInspectionRejected, wantReason: reject.Reason, wantNextID: "VK-SC-0001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := &fakeServiceCenter{verdict: tt.verdict}
			f := fleet.New(sc, defaultPolicy)
			scooter := newScooter(t, f, "Ninebot", 7, 0.4)

			err := f.AddVehicle(scooter)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddVehicle() error = %v, ожидалась %v", err, tt.wantErr)
			}
			if tt.wantReason != "" && !strings.Contains(err.Error(), tt.wantReason) {
				t.Errorf("AddVehicle() error = %v, ожидалась причина %q", err, tt.wantReason)
			}
			if !slices.Equal(sc.inspected, []string{"VK-SC-0001"}) {
				t.Errorf("техосмотр проверил %v, ожидался [VK-SC-0001]", sc.inspected)
			}
			if got := len(f.Inventory()); got != tt.wantOnBoard {
				t.Errorf("на балансе %d объектов, ожидалось %d", got, tt.wantOnBoard)
			}
			// Отказ не расходует номер: следующий самокат получит тот же (ADR-0004).
			if got := f.NextID(domain.KindElectricScooter); got != tt.wantNextID {
				t.Errorf("NextID() = %q, ожидался %q", got, tt.wantNextID)
			}
		})
	}
}

// TestAddVehicleChecksIDBeforeInspection: объект с чужим номером отклоняется до техосмотра.
func TestAddVehicleChecksIDBeforeInspection(t *testing.T) {
	t.Parallel()

	sc := &fakeServiceCenter{verdict: accept}
	f := fleet.New(sc, defaultPolicy)
	scooter, err := domain.NewElectricScooter("VK-SC-0042", "Ninebot", 7, 0.4)
	if err != nil {
		t.Fatalf("NewElectricScooter() error = %v", err)
	}

	err = f.AddVehicle(scooter)

	if !errors.Is(err, fleet.ErrUnexpectedID) {
		t.Fatalf("AddVehicle() error = %v, ожидалась %v", err, fleet.ErrUnexpectedID)
	}
	if len(sc.inspected) != 0 {
		t.Errorf("техосмотр вызван для %v, ожидалось без вызовов", sc.inspected)
	}
}

// TestAddItem проверяет постановку вещей: без техосмотра, транспорт в обход не проходит,
// повторный номер не принимается.
func TestAddItem(t *testing.T) {
	t.Parallel()

	t.Run("вещь ставится без техосмотра", func(t *testing.T) {
		t.Parallel()

		sc := &fakeServiceCenter{verdict: reject}
		f := fleet.New(sc, defaultPolicy)
		helmet := must(domain.NewHelmet(f.NextID(domain.KindHelmet), "Шлем-M"))(t)

		err := f.AddItem(helmet)
		if err != nil {
			t.Fatalf("AddItem() error = %v", err)
		}
		if len(sc.inspected) != 0 {
			t.Errorf("техосмотр вызван для %v, ожидалось без вызовов", sc.inspected)
		}
	})

	t.Run("транспорт через AddItem не проходит", func(t *testing.T) {
		t.Parallel()

		f := fleet.New(&fakeServiceCenter{verdict: accept}, defaultPolicy)
		scooter := newScooter(t, f, "Ninebot", 7, 0.4)

		err := f.AddItem(scooter)

		if !errors.Is(err, fleet.ErrVehicleAsItem) {
			t.Fatalf("AddItem() error = %v, ожидалась %v", err, fleet.ErrVehicleAsItem)
		}
		if len(f.Inventory()) != 0 {
			t.Errorf("на балансе %d объектов, ожидалось 0", len(f.Inventory()))
		}
	})

	t.Run("повторная постановка того же объекта", func(t *testing.T) {
		t.Parallel()

		f := fleet.New(&fakeServiceCenter{verdict: accept}, defaultPolicy)
		cabinet := must(domain.NewChargingCabinet(f.NextID(domain.KindChargingCabinet), "Шкаф", 1.5))(t)
		if err := f.AddItem(cabinet); err != nil {
			t.Fatalf("AddItem() error = %v", err)
		}

		err := f.AddItem(cabinet)

		if !errors.Is(err, fleet.ErrUnexpectedID) {
			t.Errorf("AddItem() error = %v, ожидалась %v", err, fleet.ErrUnexpectedID)
		}
	})
}

// TestNextID проверяет формат номера и отдельную нумерацию для каждого вида (FR-5).
func TestNextID(t *testing.T) {
	t.Parallel()

	f := fleet.New(&fakeServiceCenter{verdict: accept}, defaultPolicy)
	if err := f.AddVehicle(newScooter(t, f, "Ninebot", 7, 0.4)); err != nil {
		t.Fatalf("AddVehicle() error = %v", err)
	}

	tests := []struct {
		kind domain.Kind
		want string
	}{
		{kind: domain.KindElectricScooter, want: "VK-SC-0002"},
		{kind: domain.KindBicycle, want: "VK-BK-0001"},
		{kind: domain.KindChargingCabinet, want: "VK-CC-0001"},
	}

	for _, tt := range tests {
		t.Run(tt.kind.Title, func(t *testing.T) {
			t.Parallel()

			if got := f.NextID(tt.kind); got != tt.want {
				t.Errorf("NextID(%s) = %q, ожидался %q", tt.kind.Code, got, tt.want)
			}
		})
	}
}

// TestReportsAndLists проверяет отчёт, список для новичков и инвентарь на одном парке.
func TestReportsAndLists(t *testing.T) {
	t.Parallel()

	// Правило для новичков подменено: разрешает только Stels, хотя у Ninebot простота тоже высокая.
	f := fleet.New(&fakeServiceCenter{verdict: accept}, allowByName{"Stels": true})
	scooter := newScooter(t, f, "Ninebot", 7, 0.4)
	bike := must(domain.NewBicycle(f.NextID(domain.KindBicycle), "Stels", 9))(t)
	helmet := must(domain.NewHelmet(f.NextID(domain.KindHelmet), "Шлем-M"))(t)
	dock := must(domain.NewDockStation(f.NextID(domain.KindDockStation), "Док-1", 0.8))(t)
	for _, err := range []error{f.AddVehicle(scooter), f.AddVehicle(bike), f.AddItem(helmet), f.AddItem(dock)} {
		if err != nil {
			t.Fatalf("подготовка парка: %v", err)
		}
	}

	t.Run("отчёт", func(t *testing.T) {
		t.Parallel()

		got := f.Report()

		if got.Total != 4 || got.Vehicles != 2 || got.Items != 2 {
			t.Errorf("Report() = %+v, ожидалось 4 = 2 транспорта + 2 вещи", got)
		}
		if math.Abs(got.DailyKWh-1.2) > 1e-9 {
			t.Errorf("DailyKWh = %v, ожидалось 1.2", got.DailyKWh)
		}
	})

	t.Run("для новичков решает подменённое правило", func(t *testing.T) {
		t.Parallel()

		got := f.Beginners()

		if len(got) != 1 || got[0].Name() != "Stels" {
			t.Errorf("Beginners() = %v, ожидался только Stels", got)
		}
	})

	t.Run("инвентарь в порядке постановки", func(t *testing.T) {
		t.Parallel()

		got := f.Inventory()

		var ids []string
		for _, a := range got {
			ids = append(ids, a.ID())
		}
		want := []string{"VK-SC-0001", "VK-BK-0001", "VK-HL-0001", "VK-DS-0001"}
		if !slices.Equal(ids, want) {
			t.Errorf("Inventory() = %v, ожидалось %v", ids, want)
		}
	})
}

// TestInventoryReturnsCopy: изменение возвращённого среза не меняет парк.
func TestInventoryReturnsCopy(t *testing.T) {
	t.Parallel()

	f := fleet.New(&fakeServiceCenter{verdict: accept}, defaultPolicy)
	if err := f.AddVehicle(newScooter(t, f, "Ninebot", 7, 0.4)); err != nil {
		t.Fatalf("AddVehicle() error = %v", err)
	}

	got := f.Inventory()
	got[0] = nil

	if f.Inventory()[0] == nil {
		t.Error("Inventory() вернул внутренний срез, а не копию")
	}
}
