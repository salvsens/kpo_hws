package cli

import (
	"fmt"
	"io"
)

// demoScript — сценарии 1–7 из спецификации (раздел 4, FR-9).
// Команды выполняются тем же dispatch, что и ввод оператора.
var demoScript = []struct {
	title string
	lines []string
}{
	{"1. Приём транспорта: техосмотр пройден, выдан номер", []string{
		"add-vehicle scooter Ninebot-Max 7 0.4",
		"add-vehicle bike Stels-Navigator 9",
		"add-vehicle scooter Kugoo-G2 4 0.9",
	}},
	{"2. Отказ техосмотра: парк не изменился, номер не потерян", []string{
		"add-vehicle ebike Old-Volt 6 3.5",
		"add-vehicle ebike Eltreco-XT 6 1.2",
	}},
	{"3. Вещи на баланс без техосмотра", []string{
		"add-item helmet Шлем-M",
		"add-item cabinet Шкаф-12 1.5",
		"add-item dock Док-Покровка 0.8",
	}},
	{"4. Отчёт: велосипед и шлем в сумму энергии не входят", []string{cmdReport}},
	{"5. Транспорт для новичков: только простота ≥ 6, без вещей", []string{cmdBeginners}},
	{"6. Инвентарный список: транспорт и вещи с номерами", []string{cmdInventory}},
	{"7. Некорректный ввод: ошибка и подсказка, парк не изменился", []string{
		"add-vehicle monowheel Inmotion 5 1.0",
		"add-vehicle scooter Ninebot-Max 15 0.4",
		"add-item dock Док-2 -1",
		"add-vehicle bike Stels",
		"add-vehicle scooter Xiaomi 7 много",
		"fly",
		cmdReport,
	}},
}

// runDemo выполняет сценарии на текущем парке и печатает каждый шаг.
func (r *REPL) runDemo(out io.Writer) {
	for _, step := range demoScript {
		fmt.Fprintf(out, "\n=== %s\n", step.title)
		for _, line := range step.lines {
			fmt.Fprintf(out, "> %s\n", line)
			r.execute(line, out)
		}
	}
}
