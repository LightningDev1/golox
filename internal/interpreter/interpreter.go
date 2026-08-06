package interpreter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/environment"
	"github.com/LightningDev1/golox/internal/scanner"
)

type Interpreter struct {
	globals     *environment.Environment
	environment *environment.Environment
}

func New() *Interpreter {
	globals := environment.New()
	globals.Define("clock", &ClockFn{})
	return &Interpreter{globals: globals, environment: globals}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(statement ast.Stmt) error {
	switch stmt := statement.(type) {
	case *ast.ExpressionStmt:
		_, err := i.evaluate(stmt.Expression)
		return err

	case *ast.PrintStmt:
		value, err := i.evaluate(stmt.Expression)
		if err != nil {
			return err
		}

		fmt.Println(i.stringify(value))

	case *ast.VarStmt:
		var value any
		if stmt.Initializer != nil {
			var err error
			value, err = i.evaluate(stmt.Initializer)
			if err != nil {
				return err
			}
		}

		i.environment.Define(stmt.Name.Lexeme, value)

	case *ast.BlockStmt:
		enclosingEnv := environment.NewEnclosing(i.environment)
		return i.executeBlock(stmt.Statements, enclosingEnv)

	case *ast.IfStmt:
		condValue, err := i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}

		if i.isTruthy(condValue) {
			return i.execute(stmt.ThenBranch)
		} else if stmt.ElseBranch != nil {
			return i.execute(stmt.ThenBranch)
		}

	case *ast.WhileStmt:
		for {
			condValue, err := i.evaluate(stmt.Condition)
			if err != nil {
				return err
			}

			if !i.isTruthy(condValue) {
				break
			}

			if err = i.execute(stmt.Body); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, env *environment.Environment) error {
	previous := i.environment
	i.environment = env
	defer func() { i.environment = previous }()

	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	if expr == nil {
		return nil, nil
	}

	switch e := expr.(type) {
	case *ast.Literal:
		return e.Value, nil

	case *ast.Unary:
		right, err := i.evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Type {
		case scanner.TOKEN_MINUS:
			rightDouble, ok := right.(float64)
			if !ok {
				return nil, NewRuntimeError(e.Operator, "Operand must be a number.")
			}
			return -rightDouble, nil
		case scanner.TOKEN_BANG:
			value := !i.isTruthy(right)
			return value, nil
		}

	case *ast.Variable:
		value, ok := i.environment.Get(e.Name)
		if !ok {
			return nil, NewRuntimeError(e.Name,
				fmt.Sprintf("Undefined variable '%s'.", e.Name.Lexeme))
		}
		return value, nil

	case *ast.Assign:
		value, err := i.evaluate(e.Value)
		if err != nil {
			return nil, err
		}

		ok := i.environment.Assign(e.Name, value)
		if !ok {
			return nil, NewRuntimeError(e.Name,
				fmt.Sprintf("Undefined variable '%s'.", e.Name.Lexeme))
		}

		return value, nil

	case *ast.Binary:
		left, err := i.evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := i.evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Type {
		case scanner.TOKEN_MINUS:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l - r, nil

		case scanner.TOKEN_PLUS:
			if l, ok := left.(float64); ok {
				if r, ok := right.(float64); ok {
					return l + r, nil
				}
			}
			if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r, nil
				}
			}
			return nil, NewRuntimeError(e.Operator, "Operands must be two numbers or two strings.")

		case scanner.TOKEN_SLASH:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l / r, nil

		case scanner.TOKEN_STAR:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l * r, nil

		case scanner.TOKEN_GREATER:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l > r, nil

		case scanner.TOKEN_GREATER_EQUAL:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l >= r, nil

		case scanner.TOKEN_LESS:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l < r, nil

		case scanner.TOKEN_LESS_EQUAL:
			l, r, err := i.checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return l <= r, nil

		case scanner.TOKEN_BANG_EQUAL:
			return !i.isEqual(left, right), nil

		case scanner.TOKEN_EQUAL_EQUAL:
			return i.isEqual(left, right), nil
		}

	case *ast.Grouping:
		return i.evaluate(e.Expression)

	case *ast.Logical:
		left, err := i.evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Type {
		case scanner.TOKEN_OR:
			if i.isTruthy(left) {
				return left, nil
			}
		case scanner.TOKEN_AND:
			if !i.isTruthy(left) {
				return left, nil
			}
		}

		return i.evaluate(e.Right)

	case *ast.Call:
		callee, err := i.evaluate(e.Callee)
		if err != nil {
			return nil, err
		}

		var arguments []any
		for _, arg := range e.Arguments {
			argValue, err := i.evaluate(arg)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argValue)
		}

		function, ok := callee.(LoxCallable)
		if !ok {
			return nil, NewRuntimeError(e.Paren, "Can only call functions and classes.")
		}

		if len(arguments) != function.Arity() {
			return nil, NewRuntimeError(
				e.Paren,
				fmt.Sprintf("Expected %d arguments but got %d.",
					function.Arity(), len(arguments)),
			)
		}

		return function.Call(i, arguments)
	}
	return nil, nil
}

func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if valueBool, ok := value.(bool); ok {
		return valueBool
	}

	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return reflect.DeepEqual(a, b)
}

func (i *Interpreter) checkNumberOperands(op scanner.Token, left, right any) (float64, float64, error) {
	l, okL := left.(float64)
	r, okR := right.(float64)
	if okL && okR {
		return l, r, nil
	}
	return 0, 0, NewRuntimeError(op, "Operands must be numbers.")
}
