package domain

// Helmet — шлем: вещь на балансе без потребления энергии.
type Helmet struct {
	baseAsset
}

// NewHelmet создаёт шлем, проверяя номер и название.
func NewHelmet(id, name string) (*Helmet, error) {
	asset, err := newBaseAsset(KindHelmet, id, name)
	if err != nil {
		return nil, err
	}
	return &Helmet{baseAsset: asset}, nil
}

// DockStation — док-станция для парковки и зарядки транспорта.
type DockStation struct {
	baseAsset
	energyUsage
}

// NewDockStation создаёт док-станцию, проверяя номер, название и расход.
func NewDockStation(id, name string, dailyKWh float64) (*DockStation, error) {
	asset, energy, err := newElectricItem(KindDockStation, id, name, dailyKWh)
	if err != nil {
		return nil, err
	}
	return &DockStation{baseAsset: asset, energyUsage: energy}, nil
}

// ChargingCabinet — зарядный шкаф для аккумуляторов.
type ChargingCabinet struct {
	baseAsset
	energyUsage
}

// NewChargingCabinet создаёт зарядный шкаф, проверяя номер, название и расход.
func NewChargingCabinet(id, name string, dailyKWh float64) (*ChargingCabinet, error) {
	asset, energy, err := newElectricItem(KindChargingCabinet, id, name, dailyKWh)
	if err != nil {
		return nil, err
	}
	return &ChargingCabinet{baseAsset: asset, energyUsage: energy}, nil
}

// newElectricItem собирает общую часть вещи, потребляющей энергию.
func newElectricItem(kind Kind, id, name string, dailyKWh float64) (baseAsset, energyUsage, error) {
	asset, err := newBaseAsset(kind, id, name)
	if err != nil {
		return baseAsset{}, energyUsage{}, err
	}
	energy, err := newEnergyUsage(dailyKWh)
	if err != nil {
		return baseAsset{}, energyUsage{}, err
	}
	return asset, energy, nil
}

// Проверки на этапе компиляции: структуры реализуют нужные интерфейсы.
var (
	_ Asset          = (*Helmet)(nil)
	_ Asset          = (*DockStation)(nil)
	_ EnergyConsumer = (*DockStation)(nil)
	_ Asset          = (*ChargingCabinet)(nil)
	_ EnergyConsumer = (*ChargingCabinet)(nil)
)
