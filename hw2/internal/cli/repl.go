// Package cli provides commands for the gradebook.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
)

const (
	// studentsCommand lists the students.
	studentsCommand = "students"
	// criteriaCommand lists the criteria.
	criteriaCommand = "criteria"
)

var (
	// errInvalidArguments indicates that the command has the wrong number of arguments.
	errInvalidArguments = errors.New("invalid arguments")
	// errUnknownCommand indicates that the command is not in the list.
	errUnknownCommand = errors.New("unknown command")
)

// command defines the syntax and operation for one command.
type command struct {
	name   string
	usage  string
	run    func(string) (string, error)
	noArgs bool
	exits  bool
}

// commandResult gives the output, error, and status of one command.
type commandResult struct {
	text string
	err  error
	done bool
}

// REPL reads commands and returns results from the gradebook.
type REPL struct {
	gradebook *app.Gradebook
	newDemo   func() *app.Gradebook
	commands  []command
}

// New creates a console with a factory for empty demo journals.
func New(gradebook *app.Gradebook, newDemo func() *app.Gradebook) *REPL {
	repl := &REPL{
		gradebook: gradebook,
		newDemo:   newDemo,
	}
	repl.commands = []command{
		{name: "add-student", usage: "add-student <имя>", run: repl.addStudent},
		{name: studentsCommand, usage: studentsCommand, run: repl.listStudents, noArgs: true},
		{name: "add-criterion", usage: "add-criterion <максимум> <вес> <название>", run: repl.addCriterion},
		{name: criteriaCommand, usage: criteriaCommand, run: repl.listCriteria, noArgs: true},
		{name: "set-grade", usage: "set-grade <студент> <критерий> <балл>", run: repl.setGrade},
		{name: "grades", usage: "grades <студент>", run: repl.listGrades},
		{name: "edit-grade", usage: "edit-grade <студент> <критерий> <балл>", run: repl.editGrade},
		{name: "cancel-grade", usage: "cancel-grade <студент> <критерий>", run: repl.cancelGrade},
		{name: "result", usage: "result <студент> <current|final>", run: repl.showResult},
		{name: "export", usage: "export <студент> <current|final> <файл.csv>", run: repl.exportResult},
		{name: "demo", usage: "demo <файл.csv>", run: repl.runDemo},
		{name: "help", usage: "help", run: repl.help, noArgs: true},
		{name: "exit", usage: "exit", run: exit, noArgs: true, exits: true},
	}
	return repl
}

// Run reads commands until exit or the end of input.
// It prints command errors and returns input or output errors.
func (r *REPL) Run(in io.Reader, out io.Writer) error {
	if _, err := io.WriteString(out, "ПОС. Введите help, чтобы увидеть команды.\n"); err != nil {
		return fmt.Errorf("write welcome message: %w", err)
	}

	reader := bufio.NewReader(in)
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return fmt.Errorf("read command: %w", readErr)
		}
		if text := strings.TrimSpace(line); text != "" {
			response := r.execute(text)
			if _, err := io.WriteString(out, response.text+"\n"); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			if response.done {
				return nil
			}
		}
		if readErr != nil {
			return nil
		}
	}
}

// execute runs one command and prepares its output.
func (r *REPL) execute(line string) commandResult {
	name, args := splitWord(line)
	index := slices.IndexFunc(r.commands, func(cmd command) bool {
		return cmd.name == name
	})
	if index < 0 {
		return commandResult{
			text: fmt.Sprintf("Ошибка. Неизвестная команда %q. Введите help.", name),
			err:  errUnknownCommand,
		}
	}
	cmd := r.commands[index]
	if cmd.noArgs && args != "" {
		return usageError(cmd.usage)
	}
	text, err := cmd.run(args)
	if errors.Is(err, errInvalidArguments) {
		return usageError(cmd.usage)
	}
	if err != nil {
		return commandResult{
			text: text + "Ошибка. " + explainError(err),
			err:  err,
		}
	}
	return commandResult{text: text, done: cmd.exits}
}

// help returns the current list of commands.
func (r *REPL) help(_ string) (string, error) {
	lines := []string{"Доступные команды:"}
	for _, cmd := range r.commands {
		lines = append(lines, "  "+cmd.usage)
	}
	return strings.Join(lines, "\n"), nil
}

// exit returns the message for the exit command.
func exit(_ string) (string, error) {
	return "Работа завершена.", nil
}

// splitWord returns the first word and the rest of the line.
func splitWord(line string) (string, string) {
	index := strings.IndexFunc(line, unicode.IsSpace)
	if index < 0 {
		return line, ""
	}
	return line[:index], strings.TrimSpace(line[index:])
}

// usageError returns the syntax for a command with invalid arguments.
func usageError(usage string) commandResult {
	return commandResult{
		text: "Ошибка. Используйте команду: " + usage,
		err:  errInvalidArguments,
	}
}
