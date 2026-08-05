package environment

import (
	"github.com/LightningDev1/golox/internal/scanner"
)

type Environment struct {
	values map[string]any
}

func New() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name scanner.Token) (any, bool) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, true
	}

	return nil, false
}

func (e *Environment) Assign(name scanner.Token, value any) bool {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return true
	}

	return false
}
