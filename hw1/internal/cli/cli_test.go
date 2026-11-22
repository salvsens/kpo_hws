package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"example.com/vyshkat/internal/cli"
	"example.com/vyshkat/internal/domain"
	"example.com/vyshkat/internal/fleet"
)

// stubServiceCenter — заглушка техосмотра: всегда возвращает одно решение.
type stubServiceCenter struct {
	verdict fleet.Verdict
}

func (s stubServiceCenter) Inspect(domain.Vehicle) fleet.Verdict { return s.verdict }

// realServiceCenter — техосмотр по умолчанию, как в приложении.
var realServiceCenter = fleet.EnergyLimitInspection{MaxDailyKWh: fleet.DefaultMaxDailyKWh}

// runREPL запускает REPL на строках ввода и возвращает весь вывод.
func runREPL(t *testing.T, sc fleet.ServiceCenter, lines ...string) string {
	t.Helper()
	f := fleet.New(sc, fleet.MinSimplicityPolicy{Min: fleet.DefaultMinSimplicity})
	var out bytes.Buffer
	if err := cli.New(f).Run(strings.NewReader(strings.Join(lines, "\n")), &out); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	return out.String()
}

// TestCommands — table-driven тест команд: ввод → фрагменты, которые должны быть в выводе.
func TestCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sc      fleet.ServiceCenter
		input   []string
		want    []string
		notWant []string
	}{
		{
			name:  "приём самоката",
			input: []string{"add-vehicle scooter Ninebot 7 0.4"},
			want:  []string{"Техосмотр пройден. Принят: Ninebot, инв. № VK-SC-0001"},
		},
		{
			name:  "отказ техосмотра — причина и подсказка, парк пуст",
			sc:    stubServiceCenter{verdict: fleet.Verdict{Reason: "тормоза неисправны"}},
			input: []string{"add-vehicle bike Stels 9", "inventory"},
			want: []string{
				"Ошибка: транспорт не прошёл техосмотр: тормоза неисправны → транспорт не принят",
				"На балансе пока ничего нет.",
			},
			notWant: []string{"Техосмотр пройден"},
		},
		{
			name:  "вещи и инвентарь",
			input: []string{"add-item helmet Шлем-M", "add-item dock Док-1 0.8", "inventory"},
			want: []string{
				"На баланс поставлен: Шлем-M, инв. № VK-HL-0001",
				"VK-DS-0001  док-станция  Док-1",
			},
		},
		{
			name:  "отчёт",
			input: []string{"add-vehicle scooter Ninebot 7 0.4", "add-item helmet Шлем-M", "report"},
			want:  []string{"Всего: 2 (транспорт: 1, вещи: 1)", "Суточное энергопотребление: 0.40 кВт·ч"},
		},
		{
			name:  "для новичков — граница 6",
			input: []string{"add-vehicle bike Easy 6", "add-vehicle bike Hard 5", "beginners"},
			// Названия есть и в строках «Принят: …», поэтому проверяем строки списка целиком.
			want:    []string{"VK-BK-0001  велосипед  Easy  простота 6"},
			notWant: []string{"Hard  простота 5"},
		},
		{
			name:  "пустой список для новичков",
			input: []string{"beginners"},
			want:  []string{"Транспорта для новичков нет."},
		},
		{
			name:  "неизвестная команда",
			input: []string{"fly"},
			want:  []string{`Ошибка: неизвестная команда: "fly" → введите help`},
		},
		{
			name:  "неизвестный вид транспорта",
			input: []string{"add-vehicle monowheel X 5 1"},
			want:  []string{`неизвестный вид "monowheel", доступны: bike, ebike, scooter`},
		},
		{
			name:  "неизвестный вид вещи",
			input: []string{"add-item drone X"},
			want:  []string{"доступны: cabinet, dock, helmet"},
		},
		{
			name:  "лишний аргумент у велосипеда",
			input: []string{"add-vehicle bike Stels 9 0.4"},
			want:  []string{"нужно: add-vehicle bike <название> <простота> →"},
		},
		{
			name:  "нет аргументов",
			input: []string{"add-vehicle", "add-item"},
			want:  []string{"нужно: add-vehicle <вид>", "нужно: add-item <вид>"},
		},
		{
			name:  "не хватает расхода у шкафа",
			input: []string{"add-item cabinet Шкаф"},
			want:  []string{"нужно: add-item cabinet <название> <кВт·ч>"},
		},
		{
			name:  "простота вне шкалы",
			input: []string{"add-vehicle scooter Ninebot 15 0.4"},
			want:  []string{"→ укажите простоту целым числом от 1 до 10"},
		},
		{
			name:  "простота не число",
			input: []string{"add-vehicle bike Stels семь"},
			want:  []string{`простота "семь" — не целое число →`},
		},
		{
			name:  "расход не число",
			input: []string{"add-item dock Док-1 много"},
			want:  []string{`расход "много" — не число →`},
		},
		{
			name:  "отрицательный расход",
			input: []string{"add-vehicle ebike Eltreco 6 -1"},
			want:  []string{"→ укажите суточный расход больше 0"},
		},
		{
			name:  "справка и выход",
			input: []string{"", "help", "exit", "report"},
			want:  []string{"Команды:", "Пока!"},
			// После exit команды не выполняются.
			notWant: []string{"Всего:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := tt.sc
			if sc == nil {
				sc = realServiceCenter
			}

			out := runREPL(t, sc, tt.input...)

			for _, w := range tt.want {
				if !strings.Contains(out, w) {
					t.Errorf("в выводе нет %q:\n%s", w, out)
				}
			}
			for _, nw := range tt.notWant {
				if strings.Contains(out, nw) {
					t.Errorf("в выводе не должно быть %q:\n%s", nw, out)
				}
			}
		})
	}
}

// TestRunStopsAtEOF: конец ввода без exit завершает REPL без ошибки и без прощания.
func TestRunStopsAtEOF(t *testing.T) {
	t.Parallel()

	out := runREPL(t, realServiceCenter, "report")

	if !strings.Contains(out, "Всего: 0") || strings.Contains(out, "Пока!") {
		t.Errorf("неожиданный вывод:\n%s", out)
	}
}

// TestDemo проверяет, что demo проходит все сценарии спецификации (FR-9).
func TestDemo(t *testing.T) {
	t.Parallel()

	out := runREPL(t, realServiceCenter, "demo", "exit")

	want := []string{
		"=== 1.", "=== 2.", "=== 3.", "=== 4.", "=== 5.", "=== 6.", "=== 7.",
		"инв. № VK-SC-0001",
		"транспорт не прошёл техосмотр",
		"Принят: Eltreco-XT, инв. № VK-EB-0001",
		"Всего: 7 (транспорт: 4, вещи: 3)",
		"Суточное энергопотребление: 4.80 кВт·ч",
		`неизвестный вид "monowheel"`,
		"Пока!",
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("в выводе demo нет %q", w)
		}
	}
	// Kugoo-G2 с простотой 4 не должен попасть в список для новичков (сценарий 5).
	_, afterFive, ok := strings.Cut(out, "=== 5.")
	if !ok {
		t.Fatal("в выводе demo нет сценария 5")
	}
	beginners, _, _ := strings.Cut(afterFive, "=== 6.")
	if strings.Contains(beginners, "Kugoo-G2") {
		t.Errorf("Kugoo-G2 (простота 4) в списке для новичков:\n%s", beginners)
	}
}
