package domain

import (
	"fmt"
	"strings"
)

// Asset — объект на балансе «ВышКат»: транспорт или вещь с инвентарным номером.
type Asset interface {
	Kind() Kind
	ID() string
	Name() string
}

// baseAsset хранит общие для всех объектов на балансе поля.
// Встраивается в каждую структуру домена, чтобы не дублировать Kind(), ID() и Name().
type baseAsset struct {
	kind Kind
	id   string
	name string
}

// newBaseAsset проверяет номер и название; пробелы по краям отбрасываются.
func newBaseAsset(kind Kind, id, name string) (baseAsset, error) {
	id, name = strings.TrimSpace(id), strings.TrimSpace(name)
	if id == "" || name == "" {
		return baseAsset{}, fmt.Errorf("%w: номер %q, название %q", ErrInvalidAsset, id, name)
	}
	return baseAsset{kind: kind, id: id, name: name}, nil
}

// Kind возвращает вид объекта.
func (a baseAsset) Kind() Kind { return a.kind }

// ID возвращает инвентарный номер.
func (a baseAsset) ID() string { return a.id }

// Name возвращает название объекта.
func (a baseAsset) Name() string { return a.name }
