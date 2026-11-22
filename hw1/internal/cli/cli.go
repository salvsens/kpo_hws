// Package cli реализует текстовый интерфейс (REPL) парка «ВышКат».
package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"example.com/vyshkat/internal/fleet"
)

const helpText = `Команды:
  add-vehicle <вид> <название> <простота> [кВт·ч]  принять транспорт после техосмотра
                                                   вид: scooter, bike, ebike; кВт·ч — для scooter и ebike
  add-item <вид> <название> [кВт·ч]                поставить вещь на баланс
                                                   вид: helmet, dock, cabinet; кВт·ч — для dock и cabinet
  report                                           число единиц и суточное энергопотребление
  beginners                                        транспорт для новичков
  inventory                                        все объекты на балансе с инвентарными номерами
  demo                                             показать все сценарии на тестовых данных
  help                                             эта справка
  exit                                             выйти
Название — одно слово без пробелов, дробные числа — через точку: 0.4`

// Команды консоли (спецификация, раздел 4).
const (
	cmdAddVehicle = "add-vehicle"
	cmdAddItem    = "add-item"
	cmdReport     = "report"
	cmdBeginners  = "beginners"
	cmdInventory  = "inventory"
	cmdDemo       = "demo"
	cmdHelp       = "help"
	cmdExit       = "exit"
	cmdQuit       = "quit"
)

// REPL — консольный интерфейс парка: читает команды и печатает результат.
type REPL struct {
	fleet *fleet.Fleet
}

// New создаёт REPL поверх парка.
func New(f *fleet.Fleet) *REPL {
	return &REPL{fleet: f}
}

// Run запускает цикл «команда — ответ» до команды exit или конца ввода.
func (r *REPL) Run(in io.Reader, out io.Writer) error {
	fmt.Fprintln(out, "«ВышКат» — учёт парка. Введите help для списка команд.")
	sc := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, "> ")
		if !sc.Scan() {
			return sc.Err()
		}
		if r.execute(sc.Text(), out) {
			fmt.Fprintln(out, "Пока!")
			return nil
		}
	}
}

// execute выполняет одну строку и печатает ошибку, если она есть.
// Возвращает true, если пользователь попросил выйти.
func (r *REPL) execute(line string, out io.Writer) (done bool) {
	done, err := r.dispatch(line, out)
	if err != nil {
		printError(out, err)
	}
	return done
}

// dispatch разбирает одну строку и выполняет команду.
func (r *REPL) dispatch(line string, out io.Writer) (done bool, err error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, nil
	}

	switch cmd, args := fields[0], fields[1:]; cmd {
	case cmdExit, cmdQuit:
		return true, nil
	case cmdHelp:
		fmt.Fprintln(out, helpText)
	case cmdAddVehicle:
		return false, r.addVehicle(args, out)
	case cmdAddItem:
		return false, r.addItem(args, out)
	case cmdReport:
		r.printReport(out)
	case cmdBeginners:
		return false, r.printBeginners(out)
	case cmdInventory:
		return false, r.printInventory(out)
	case cmdDemo:
		r.runDemo(out)
	default:
		return false, fmt.Errorf("%w: %q", errUnknownCommand, cmd)
	}
	return false, nil
}
