package domain

import "fmt"

// Границы шкалы «простоты для новичка».
const (
	MinSimplicity = 1
	MaxSimplicity = 10
)

// Vehicle — прокатный транспорт. Любой транспорт стоит на балансе, поэтому встраивает Asset.
type Vehicle interface {
	Asset
	BeginnerSimplicity() int
}

// baseVehicle — общая часть всех видов транспорта.
type baseVehicle struct {
	baseAsset
	simplicity int
}

// newBaseVehicle проверяет номер, название и простоту транспорта.
func newBaseVehicle(kind Kind, id, name string, simplicity int) (baseVehicle, error) {
	asset, err := newBaseAsset(kind, id, name)
	if err != nil {
		return baseVehicle{}, err
	}
	if simplicity < MinSimplicity || simplicity > MaxSimplicity {
		return baseVehicle{}, fmt.Errorf("%w: %d", ErrInvalidSimplicity, simplicity)
	}
	return baseVehicle{baseAsset: asset, simplicity: simplicity}, nil
}

// BeginnerSimplicity возвращает простоту для новичка по шкале от 1 до 10.
func (v baseVehicle) BeginnerSimplicity() int { return v.simplicity }

// newElectricVehicle собирает общую часть транспорта с электромотором.
func newElectricVehicle(kind Kind, id, name string, simplicity int, dailyKWh float64) (baseVehicle, energyUsage, error) {
	vehicle, err := newBaseVehicle(kind, id, name, simplicity)
	if err != nil {
		return baseVehicle{}, energyUsage{}, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return baseVehicle{}, energyUsage{}, err
	}
	return vehicle, energy, nil
}

// Bicycle — обычный велосипед: транспорт без потребления энергии.
type Bicycle struct {
	baseVehicle
}

// NewBicycle создаёт велосипед, проверяя номер, название и простоту.
func NewBicycle(id, name string, simplicity int) (*Bicycle, error) {
	vehicle, err := newBaseVehicle(KindBicycle, id, name, simplicity)
	if err != nil {
		return nil, err
	}
	return &Bicycle{baseVehicle: vehicle}, nil
}

// ElectricScooter — электросамокат.
type ElectricScooter struct {
	baseVehicle
	energyUsage
}

// NewElectricScooter создаёт электросамокат, проверяя номер, название, простоту и расход.
func NewElectricScooter(id, name string, simplicity int, dailyKWh float64) (*ElectricScooter, error) {
	vehicle, energy, err := newElectricVehicle(KindElectricScooter, id, name, simplicity, dailyKWh)
	if err != nil {
		return nil, err
	}
	return &ElectricScooter{baseVehicle: vehicle, energyUsage: energy}, nil
}

// ElectricBike — электровелосипед.
type ElectricBike struct {
	baseVehicle
	energyUsage
}

// NewElectricBike создаёт электровелосипед, проверяя номер, название, простоту и расход.
func NewElectricBike(id, name string, simplicity int, dailyKWh float64) (*ElectricBike, error) {
	vehicle, energy, err := newElectricVehicle(KindElectricBike, id, name, simplicity, dailyKWh)
	if err != nil {
		return nil, err
	}
	return &ElectricBike{baseVehicle: vehicle, energyUsage: energy}, nil
}

// Проверки на этапе компиляции: структуры реализуют нужные интерфейсы.
var (
	_ Vehicle        = (*Bicycle)(nil)
	_ Vehicle        = (*ElectricScooter)(nil)
	_ EnergyConsumer = (*ElectricScooter)(nil)
	_ Vehicle        = (*ElectricBike)(nil)
	_ EnergyConsumer = (*ElectricBike)(nil)
)
