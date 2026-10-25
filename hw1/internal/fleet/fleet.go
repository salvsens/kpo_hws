package fleet

import (
	"fmt"
	"slices"

	"example.com/vyshkat/internal/domain"
)

// idFormat — формат инвентарного номера: VK-<код вида>-<NNNN> (FR-5).
const idFormat = "VK-%s-%04d"

// Report — сводка по парку (FR-6).
type Report struct {
	Total    int     // всего объектов на балансе
	Vehicles int     // из них транспорта
	Items    int     // из них вещей
	DailyKWh float64 // суммарный суточный расход потребителей энергии
}

// Fleet — парк «ВышКат»: объекты на балансе в порядке постановки.
type Fleet struct {
	inspector ServiceCenter
	policy    BeginnerPolicy
	assets    []domain.Asset
	accepted  map[domain.Kind]int // сколько объектов каждого вида на балансе
}

// New создаёт пустой парк с заданным техосмотром и правилом для новичков.
func New(sc ServiceCenter, policy BeginnerPolicy) *Fleet {
	return &Fleet{
		inspector: sc,
		policy:    policy,
		accepted:  make(map[domain.Kind]int),
	}
}

// NextID возвращает следующий свободный номер вида и не занимает его:
// номер считается занятым только после постановки на баланс (ADR-0004).
func (f *Fleet) NextID(k domain.Kind) string {
	return fmt.Sprintf(idFormat, k.Code, f.accepted[k]+1)
}

// AddVehicle проводит техосмотр и при успехе ставит транспорт на баланс.
func (f *Fleet) AddVehicle(v domain.Vehicle) error {
	if err := f.checkID(v); err != nil {
		return err
	}
	if verdict := f.inspector.Inspect(v); !verdict.Accepted {
		return fmt.Errorf("%w: %s", ErrInspectionRejected, verdict.Reason)
	}
	f.register(v)
	return nil
}

// AddItem ставит вещь на баланс без техосмотра.
// Транспорт через AddItem не принимается: он обязан пройти техосмотр.
func (f *Fleet) AddItem(a domain.Asset) error {
	if _, ok := a.(domain.Vehicle); ok {
		return fmt.Errorf("%w: %s — используйте AddVehicle", ErrVehicleAsItem, a.ID())
	}
	if err := f.checkID(a); err != nil {
		return err
	}
	f.register(a)
	return nil
}

// Report считает число объектов и суммарный суточный расход энергии.
func (f *Fleet) Report() Report {
	r := Report{Total: len(f.assets)}
	for _, a := range f.assets {
		if _, ok := a.(domain.Vehicle); ok {
			r.Vehicles++
		} else {
			r.Items++
		}
		if consumer, ok := a.(domain.EnergyConsumer); ok {
			r.DailyKWh += consumer.DailyEnergyKWh()
		}
	}
	return r
}

// Beginners возвращает транспорт, подходящий для новичков по правилу BeginnerPolicy.
func (f *Fleet) Beginners() []domain.Vehicle {
	var result []domain.Vehicle
	for _, a := range f.assets {
		if v, ok := a.(domain.Vehicle); ok && f.policy.Allows(v) {
			result = append(result, v)
		}
	}
	return result
}

// Inventory возвращает копию списка объектов на балансе в порядке постановки.
// Копия — чтобы вызывающий код не мог изменить парк в обход методов.
func (f *Fleet) Inventory() []domain.Asset {
	return slices.Clone(f.assets)
}

// checkID проверяет, что номер объекта — тот, который парк выдал бы сейчас.
func (f *Fleet) checkID(a domain.Asset) error {
	if want := f.NextID(a.Kind()); a.ID() != want {
		return fmt.Errorf("%w: получен %s, ожидался %s", ErrUnexpectedID, a.ID(), want)
	}
	return nil
}

// register ставит объект на баланс и занимает его номер.
func (f *Fleet) register(a domain.Asset) {
	f.assets = append(f.assets, a)
	f.accepted[a.Kind()]++
}
