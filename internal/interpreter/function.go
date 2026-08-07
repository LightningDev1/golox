package interpreter

import (
	"errors"
	"fmt"

	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/environment"
)

type LoxFunction struct {
	declaration *ast.FunctionStmt
	closure     *environment.Environment
}

func NewLoxFunction(declaration *ast.FunctionStmt, closure *environment.Environment) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(i *Interpreter, arguments []any) (any, error) {
	env := environment.NewEnclosing(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := i.executeBlock(f.declaration.Body, env)
	if ret, ok := errors.AsType[*ReturnError](err); ok {
		return ret.Value, nil
	}

	return nil, err
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
