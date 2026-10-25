package cli

import (
	"maps"
	"slices"
	"strings"

	"example.com/vyshkat/internal/domain"
)

// vehicleSpec описывает, как консоль создаёт транспорт одного вида.
type vehicleSpec struct {
	kind     domain.Kind
	electric bool // нужен ли аргумент кВт·ч
	build    func(id, name string, simplicity int, dailyKWh float64) (domain.Vehicle, error)
}

// itemSpec описывает, как консоль создаёт вещь одного вида.
type itemSpec struct {
	kind     domain.Kind
	electric bool // нужен ли аргумент кВт·ч
	build    func(id, name string, dailyKWh float64) (domain.Asset, error)
}

// vehicleSpecs — слова консоли для видов транспорта.
// Новый вид транспорта — одна строка здесь и тип в domain; парк не меняется.
var vehicleSpecs = map[string]vehicleSpec{
	"scooter": {
		kind: domain.KindElectricScooter, electric: true,
		build: func(id, name string, simplicity int, kwh float64) (domain.Vehicle, error) {
			return asVehicle(domain.NewElectricScooter(id, name, simplicity, kwh))
		},
	},
	"bike": {
		kind: domain.KindBicycle,
		build: func(id, name string, simplicity int, _ float64) (domain.Vehicle, error) {
			return asVehicle(domain.NewBicycle(id, name, simplicity))
		},
	},
	"ebike": {
		kind: domain.KindElectricBike, electric: true,
		build: func(id, name string, simplicity int, kwh float64) (domain.Vehicle, error) {
			return asVehicle(domain.NewElectricBike(id, name, simplicity, kwh))
		},
	},
}

// itemSpecs — слова консоли для видов вещей.
var itemSpecs = map[string]itemSpec{
	"helmet": {
		kind: domain.KindHelmet,
		build: func(id, name string, _ float64) (domain.Asset, error) {
			return asAsset(domain.NewHelmet(id, name))
		},
	},
	"dock": {
		kind: domain.KindDockStation, electric: true,
		build: func(id, name string, kwh float64) (domain.Asset, error) {
			return asAsset(domain.NewDockStation(id, name, kwh))
		},
	},
	"cabinet": {
		kind: domain.KindChargingCabinet, electric: true,
		build: func(id, name string, kwh float64) (domain.Asset, error) {
			return asAsset(domain.NewChargingCabinet(id, name, kwh))
		},
	},
}

// asVehicle приводит результат конструктора к интерфейсу.
// При ошибке возвращает настоящий nil, а не интерфейс с nil-указателем внутри.
func asVehicle[T domain.Vehicle](v T, err error) (domain.Vehicle, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

// asAsset — то же, что asVehicle, для вещей.
func asAsset[T domain.Asset](a T, err error) (domain.Asset, error) {
	if err != nil {
		return nil, err
	}
	return a, nil
}

// kindNames возвращает слова консоли через запятую в алфавитном порядке.
func kindNames[S any](specs map[string]S) string {
	return strings.Join(slices.Sorted(maps.Keys(specs)), ", ")
}
