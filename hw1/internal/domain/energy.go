package domain

import "fmt"

// EnergyConsumer — устройство, потребляющее электроэнергию.
// Намеренно не встраивает Asset: расчёту энергии нужен только расход (ISP).
type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

// energyUsage хранит суточный расход и встраивается в «электрические» структуры.
type energyUsage struct {
	dailyKWh float64
}

// newEnergyUsage проверяет, что суточный расход больше нуля.
func newEnergyUsage(dailyKWh float64) (energyUsage, error) {
	if dailyKWh <= 0 {
		return energyUsage{}, fmt.Errorf("%w: %v", ErrInvalidEnergy, dailyKWh)
	}
	return energyUsage{dailyKWh: dailyKWh}, nil
}

// DailyEnergyKWh возвращает суточный расход энергии в кВт·ч.
func (e energyUsage) DailyEnergyKWh() float64 { return e.dailyKWh }
