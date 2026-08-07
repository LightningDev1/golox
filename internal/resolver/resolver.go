package resolver

import (
	"slices"

	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/interpreter"
	"github.com/LightningDev1/golox/internal/scanner"
)

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func New(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          nil,
		currentFunction: FunctionTypeNone,
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) error {
	for _, statement := range statements {
		if err := r.resolveStmt(statement); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveStmt(statement ast.Stmt) error {
	switch stmt := statement.(type) {
	case *ast.BlockStmt:
		r.beginScope()
		defer r.endScope()

		return r.Resolve(stmt.Statements)

	case *ast.ExpressionStmt:
		return r.resolveExpr(stmt.Expression)

	case *ast.VarStmt:
		if err := r.declare(stmt.Name); err != nil {
			return err
		}
		if stmt.Initializer != nil {
			if err := r.resolveExpr(stmt.Initializer); err != nil {
				return err
			}
		}
		r.define(stmt.Name)

	case *ast.WhileStmt:
		if err := r.resolveExpr(stmt.Condition); err != nil {
			return err
		}
		return r.resolveStmt(stmt.Body)

	case *ast.FunctionStmt:
		if err := r.declare(stmt.Name); err != nil {
			return err
		}
		r.define(stmt.Name)
		return r.resolveFunction(stmt, FunctionTypeFunction)

	case *ast.IfStmt:
		if err := r.resolveExpr(stmt.Condition); err != nil {
			return err
		}
		if err := r.resolveStmt(stmt.ThenBranch); err != nil {
			return err
		}
		if stmt.ElseBranch != nil {
			return r.resolveStmt(stmt.ElseBranch)
		}

	case *ast.PrintStmt:
		return r.resolveExpr(stmt.Expression)

	case *ast.ReturnStmt:
		if r.currentFunction == FunctionTypeNone {
			return NewResolveError(stmt.Keyword, "Can't return from top-level code.")
		}
		if stmt.Value != nil {
			return r.resolveExpr(stmt.Value)
		}
	}
	return nil
}

func (r *Resolver) resolveExpr(expression ast.Expr) error {
	switch expr := expression.(type) {
	case *ast.Variable:
		if scope := r.getTopScope(); scope != nil {
			if defined, ok := scope[expr.Name.Lexeme]; ok && !defined {
				return NewResolveError(expr.Name, "Can't read local variable in its own initializer.")
			}
		}
		r.resolveLocal(expr, expr.Name)

	case *ast.Assign:
		if err := r.resolveExpr(expr.Value); err != nil {
			return err
		}
		r.resolveLocal(expr, expr.Name)

	case *ast.Binary:
		if err := r.resolveExpr(expr.Left); err != nil {
			return err
		}
		return r.resolveExpr(expr.Right)

	case *ast.Call:
		if err := r.resolveExpr(expr.Callee); err != nil {
			return err
		}
		for _, argument := range expr.Arguments {
			if err := r.resolveExpr(argument); err != nil {
				return err
			}
		}

	case *ast.Grouping:
		return r.resolveExpr(expr.Expression)

	case *ast.Literal:
		return nil

	case *ast.Logical:
		if err := r.resolveExpr(expr.Left); err != nil {
			return err
		}
		return r.resolveExpr(expr.Right)

	case *ast.Unary:
		return r.resolveExpr(expr.Right)
	}
	return nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name scanner.Token) {
	for i, scope := range slices.Backward(r.scopes) {
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *ast.FunctionStmt, funcType FunctionType) error {
	enclosingFunc := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	defer func() {
		r.endScope()
		r.currentFunction = enclosingFunc
	}()

	for _, param := range function.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param)
	}

	return r.Resolve(function.Body)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) getTopScope() map[string]bool {
	if len(r.scopes) == 0 {
		return nil
	}
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) declare(name scanner.Token) error {
	if scope := r.getTopScope(); scope != nil {
		if _, ok := scope[name.Lexeme]; ok {
			return NewResolveError(name, "Already a variable with this name in this scope.")
		}
		scope[name.Lexeme] = false
	}
	return nil
}

func (r *Resolver) define(name scanner.Token) {
	if scope := r.getTopScope(); scope != nil {
		scope[name.Lexeme] = true
	}
}
