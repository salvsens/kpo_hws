package cli

import (
	"errors"
	"fmt"
	"io"

	"example.com/vyshkat/internal/domain"
	"example.com/vyshkat/internal/fleet"
)

// Ошибки ввода консоли.
// Сравнивайте их через errors.Is, не по тексту сообщения.
var (
	// errUnknownCommand возвращается, если первое слово строки — не команда из справки.
	errUnknownCommand = errors.New("неизвестная команда")
	// errUnknownKind возвращается, если вида нет в vehicleSpecs или itemSpecs.
	errUnknownKind = errors.New("неизвестный вид")
	// errUsage возвращается, если у команды не хватает аргументов или их слишком много.
	// Текст обрывается на «нужно»: при оборачивании к нему дописывается правильный формат команды.
	errUsage = errors.New("неверное число аргументов, нужно")
	// errBadNumber возвращается, если простоту или расход не удалось прочитать как число.
	errBadNumber = errors.New("неверное число")
)

// hints сопоставляет ошибку с подсказкой «что делать».
// Порядок важен: берётся первая подходящая строка.
var hints = []struct {
	target error
	hint   string
}{
	{fleet.ErrInspectionRejected, "транспорт не принят; устраните причину и повторите приёмку"},
	{domain.ErrInvalidSimplicity, "укажите простоту целым числом от 1 до 10"},
	{domain.ErrInvalidEnergy, "укажите суточный расход больше 0, например 0.4"},
	{domain.ErrInvalidAsset, "укажите непустое название"},
	{errUnknownCommand, "введите help, чтобы увидеть список команд"},
	{errUnknownKind, "выберите вид из списка"},
	{errUsage, "повторите команду с нужными аргументами"},
	{errBadNumber, "числа вводите цифрами, дробные — через точку: 0.4"},
}

// printError печатает ошибку по схеме «что произошло → что делать».
func printError(out io.Writer, err error) {
	hint := "это ошибка программы; сообщите разработчику"
	for _, h := range hints {
		if errors.Is(err, h.target) {
			hint = h.hint
			break
		}
	}
	fmt.Fprintf(out, "Ошибка: %v → %s\n", err, hint)
}
