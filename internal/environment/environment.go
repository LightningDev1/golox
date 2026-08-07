package environment

import (
	"github.com/LightningDev1/golox/internal/scanner"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func New() *Environment {
	return &Environment{values: make(map[string]any)}
}

func NewEnclosing(parent *Environment) *Environment {
	return &Environment{
		enclosing: parent,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) AssignAt(distance int, name scanner.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) Get(name scanner.Token) (any, bool) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, true
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, false
}

func (e *Environment) Assign(name scanner.Token, value any) bool {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return true
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return false
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for range distance {
		environment = environment.enclosing
	}
	return environment
}
