package interpreter

import (
	"fmt"

	"github.com/LightningDev1/golox/internal/scanner"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]any)}
}

func (i *LoxInstance) Get(name scanner.Token) (any, error) {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value, nil
	}

	if method := i.class.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i), nil
	}

	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme))
}

func (i *LoxInstance) Set(name scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) String() string {
	return i.class.name + " instance"
}
