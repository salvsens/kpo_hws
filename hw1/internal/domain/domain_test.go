package domain_test

import (
	"errors"
	"testing"

	"example.com/vyshkat/internal/domain"
)

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

// TestConstructors — table-driven тест всех конструкторов:
// корректные данные дают объект нужного вида, нарушение инварианта — sentinel-ошибку.
func TestConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		build    func() (domain.Asset, error)
		wantKind domain.Kind
		wantErr  error
	}{
		{
			name:     "электросамокат",
			build:    func() (domain.Asset, error) { return domain.NewElectricScooter("VK-SC-0001", "Ninebot", 7, 0.4) },
			wantKind: domain.KindElectricScooter,
		},
		{
			name:     "велосипед",
			build:    func() (domain.Asset, error) { return domain.NewBicycle("VK-BK-0001", "Stels", 9) },
			wantKind: domain.KindBicycle,
		},
		{
			name:     "электровелосипед",
			build:    func() (domain.Asset, error) { return domain.NewElectricBike("VK-EB-0001", "Eltreco", 6, 1.2) },
			wantKind: domain.KindElectricBike,
		},
		{
			name:     "шлем",
			build:    func() (domain.Asset, error) { return domain.NewHelmet("VK-HL-0001", "Шлем-M") },
			wantKind: domain.KindHelmet,
		},
		{
			name:     "док-станция",
			build:    func() (domain.Asset, error) { return domain.NewDockStation("VK-DS-0001", "Док-1", 0.8) },
			wantKind: domain.KindDockStation,
		},
		{
			name:     "зарядный шкаф",
			build:    func() (domain.Asset, error) { return domain.NewChargingCabinet("VK-CC-0001", "Шкаф", 1.5) },
			wantKind: domain.KindChargingCabinet,
		},
		{
			name:    "пустой номер",
			build:   func() (domain.Asset, error) { return domain.NewHelmet("", "Шлем-M") },
			wantErr: domain.ErrInvalidAsset,
		},
		{
			name:    "название из пробелов",
			build:   func() (domain.Asset, error) { return domain.NewBicycle("VK-BK-0001", "   ", 5) },
			wantErr: domain.ErrInvalidAsset,
		},
		{
			name:    "простота 0 — ниже шкалы",
			build:   func() (domain.Asset, error) { return domain.NewBicycle("VK-BK-0001", "Stels", 0) },
			wantErr: domain.ErrInvalidSimplicity,
		},
		{
			name:    "простота 11 — выше шкалы",
			build:   func() (domain.Asset, error) { return domain.NewElectricScooter("VK-SC-0001", "Ninebot", 11, 0.4) },
			wantErr: domain.ErrInvalidSimplicity,
		},
		{
			name:    "нулевой расход у электровелосипеда",
			build:   func() (domain.Asset, error) { return domain.NewElectricBike("VK-EB-0001", "Eltreco", 6, 0) },
			wantErr: domain.ErrInvalidEnergy,
		},
		{
			name:    "отрицательный расход у шкафа",
			build:   func() (domain.Asset, error) { return domain.NewChargingCabinet("VK-CC-0001", "Шкаф", -1) },
			wantErr: domain.ErrInvalidEnergy,
		},
		{
			name:    "пустое название у док-станции",
			build:   func() (domain.Asset, error) { return domain.NewDockStation("VK-DS-0001", "", 0.8) },
			wantErr: domain.ErrInvalidAsset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.build()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка = %v, ожидалась %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, ожидался %v", got.Kind(), tt.wantKind)
			}
		})
	}
}

// TestFieldsAreTrimmed проверяет, что пробелы по краям номера и названия отбрасываются.
func TestFieldsAreTrimmed(t *testing.T) {
	t.Parallel()

	b := must(domain.NewBicycle("  VK-BK-0001 ", " Stels ", 9))(t)

	if b.ID() != "VK-BK-0001" || b.Name() != "Stels" {
		t.Errorf("ID, Name = %q, %q; ожидались без пробелов", b.ID(), b.Name())
	}
	if b.BeginnerSimplicity() != 9 {
		t.Errorf("BeginnerSimplicity() = %d, ожидалось 9", b.BeginnerSimplicity())
	}
}

// TestCapabilities фиксирует таблицу свойств из спецификации:
// кто транспорт и кто потребляет энергию — через проверку интерфейсов.
func TestCapabilities(t *testing.T) {
	t.Parallel()

	scooter := must(domain.NewElectricScooter("VK-SC-0001", "Ninebot", 7, 0.4))(t)
	bike := must(domain.NewBicycle("VK-BK-0001", "Stels", 9))(t)
	ebike := must(domain.NewElectricBike("VK-EB-0001", "Eltreco", 6, 1.2))(t)
	helmet := must(domain.NewHelmet("VK-HL-0001", "Шлем-M"))(t)
	dock := must(domain.NewDockStation("VK-DS-0001", "Док-1", 0.8))(t)
	cabinet := must(domain.NewChargingCabinet("VK-CC-0001", "Шкаф", 1.5))(t)

	tests := []struct {
		asset        domain.Asset
		wantVehicle  bool
		wantConsumer bool
		wantKWh      float64
	}{
		{asset: scooter, wantVehicle: true, wantConsumer: true, wantKWh: 0.4},
		{asset: bike, wantVehicle: true},
		{asset: ebike, wantVehicle: true, wantConsumer: true, wantKWh: 1.2},
		{asset: helmet},
		{asset: dock, wantConsumer: true, wantKWh: 0.8},
		{asset: cabinet, wantConsumer: true, wantKWh: 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.asset.Kind().Title, func(t *testing.T) {
			t.Parallel()

			_, isVehicle := tt.asset.(domain.Vehicle)
			consumer, isConsumer := tt.asset.(domain.EnergyConsumer)

			if isVehicle != tt.wantVehicle {
				t.Errorf("транспорт = %v, ожидалось %v", isVehicle, tt.wantVehicle)
			}
			if isConsumer != tt.wantConsumer {
				t.Fatalf("потребитель энергии = %v, ожидалось %v", isConsumer, tt.wantConsumer)
			}
			if isConsumer && consumer.DailyEnergyKWh() != tt.wantKWh {
				t.Errorf("DailyEnergyKWh() = %v, ожидалось %v", consumer.DailyEnergyKWh(), tt.wantKWh)
			}
		})
	}
}
