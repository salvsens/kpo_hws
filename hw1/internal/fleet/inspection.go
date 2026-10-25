package fleet

import (
	"fmt"

	"example.com/vyshkat/internal/domain"
)

// DefaultMaxDailyKWh — лимит суточного расхода для техосмотра по умолчанию (FR-3).
const DefaultMaxDailyKWh = 2.0

// Verdict — решение техосмотра.
type Verdict struct {
	Accepted bool
	Reason   string // причина отказа; пусто, если принят
}

// ServiceCenter — техосмотр: решает, принять ли транспорт в парк.
type ServiceCenter interface {
	Inspect(v domain.Vehicle) Verdict
}

// EnergyLimitInspection отклоняет транспорт, расход которого больше лимита.
// Транспорт без потребления энергии (велосипед) по этому правилу проходит.
type EnergyLimitInspection struct {
	MaxDailyKWh float64
}

// Inspect проверяет суточный расход транспорта.
func (s EnergyLimitInspection) Inspect(v domain.Vehicle) Verdict {
	consumer, ok := v.(domain.EnergyConsumer)
	if !ok {
		return Verdict{Accepted: true}
	}
	if kwh := consumer.DailyEnergyKWh(); kwh > s.MaxDailyKWh {
		return Verdict{
			Reason: fmt.Sprintf("суточный расход %.2f кВт·ч больше лимита %.2f кВт·ч — проверьте батарею",
				kwh, s.MaxDailyKWh),
		}
	}
	return Verdict{Accepted: true}
}
