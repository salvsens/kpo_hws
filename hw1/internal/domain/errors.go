// Package domain содержит предметную область «ВышКат»: объекты на балансе,
// транспорт и потребителей энергии (ADR-0001).
package domain

import "errors"

// Ошибки предметной области.
var (
	// ErrInvalidAsset возвращается, если у объекта нет инвентарного номера или названия.
	ErrInvalidAsset = errors.New("у объекта должны быть инвентарный номер и название")
	// ErrInvalidSimplicity возвращается, если простота вне диапазона [MinSimplicity, MaxSimplicity].
	ErrInvalidSimplicity = errors.New("простота для новичка должна быть от 1 до 10")
	// ErrInvalidEnergy возвращается, если суточный расход энергии не положительный.
	ErrInvalidEnergy = errors.New("суточный расход энергии должен быть больше нуля")
)
