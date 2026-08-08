package interpreter

import (
	"errors"
	"fmt"

	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/environment"
)

type LoxFunction struct {
	declaration   *ast.FunctionStmt
	closure       *environment.Environment
	isInitializer bool
}

func NewLoxFunction(
	declaration *ast.FunctionStmt,
	closure *environment.Environment,
	isInitializer bool,
) *LoxFunction {
	return &LoxFunction{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := environment.NewEnclosing(f.closure)
	env.Define("this", instance)
	return &LoxFunction{
		declaration:   f.declaration,
		closure:       env,
		isInitializer: f.isInitializer,
	}
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
		if f.isInitializer {
			return f.closure.GetAt(0, "this"), err
		}

		return ret.Value, nil
	} else if err != nil {
		return nil, err
	}

	if f.isInitializer {
		return f.closure.GetAt(0, "this"), err
	}

	return nil, nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
