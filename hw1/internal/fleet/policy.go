package fleet

import "example.com/vyshkat/internal/domain"

// DefaultMinSimplicity — минимальная простота для выдачи новичку по умолчанию (FR-7).
const DefaultMinSimplicity = 6

// BeginnerPolicy — правило, можно ли выдавать транспорт новичку.
type BeginnerPolicy interface {
	Allows(v domain.Vehicle) bool
}

// MinSimplicityPolicy разрешает транспорт с простотой не ниже Min.
type MinSimplicityPolicy struct {
	Min int
}

// Allows сообщает, достаточно ли прост транспорт для новичка.
func (p MinSimplicityPolicy) Allows(v domain.Vehicle) bool {
	return v.BeginnerSimplicity() >= p.Min
}
