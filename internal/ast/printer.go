package ast

import (
	"fmt"
	"strings"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(expr Expr) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *LiteralExpr:
		if e.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", e.Value)

	case *UnaryExpr:
		return p.parenthesize(e.Operator.Lexeme, e.Right)

	case *BinaryExpr:
		return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)

	case *GroupingExpr:
		return p.parenthesize("group", e.Expression)

	case *VariableExpr:
		return e.Name.Lexeme

	case *AssignExpr:
		return p.parenthesize("= "+e.Name.Lexeme, e.Value)

	case *LogicalExpr:
		return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)

	case *CallExpr:
		args := append([]Expr{e.Callee}, e.Arguments...)
		return p.parenthesize("call", args...)

	default:
		return fmt.Sprintf("(unknown %T)", expr)
	}
}

func (p *Printer) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(p.Print(expr))
	}

	builder.WriteString(")")
	return builder.String()
}
