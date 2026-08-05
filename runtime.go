package golox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/LightningDev1/golox/internal/interpreter"
	"github.com/LightningDev1/golox/internal/parser"
	"github.com/LightningDev1/golox/internal/scanner"
)

type Runtime struct {
	HadError        bool
	HadRuntimeError bool

	interpreter *interpreter.Interpreter
}

func NewRuntime() *Runtime {
	return &Runtime{interpreter: interpreter.New()}
}

func (r *Runtime) Run(source string) {
	sc := scanner.New(source)
	tokens, errs := sc.ScanTokens()

	for _, err := range errs {
		r.Report(err.Line, "", err.Message)
	}

	p := parser.New(tokens)
	statements, err := p.Parse()
	if err != nil {
		r.Report(0, "", err.Error())
	}

	if r.HadError {
		return
	}

	if err = r.interpreter.Interpret(statements); err != nil {
		r.HadRuntimeError = true
	}
}

func (r *Runtime) RunFile(file string) error {
	source, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	r.Run(string(source))

	if r.HadError {
		return fmt.Errorf("error occurred")
	}
	if r.HadRuntimeError {
		return fmt.Errorf("runtime error occurred")
	}

	return nil
}

func (r *Runtime) RunPrompt() error {
	stdin := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !stdin.Scan() {
			break
		}

		r.Run(stdin.Text())
		r.HadError = false
	}

	return stdin.Err()
}

func (r *Runtime) Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	r.HadError = true
}
