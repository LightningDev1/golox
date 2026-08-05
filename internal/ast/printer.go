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
	case *Literal:
		if e.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", e.Value)

	case *Unary:
		return p.parenthesize(e.Operator.Lexeme, e.Right)

	case *Binary:
		return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)

	case *Grouping:
		return p.parenthesize("group", e.Expression)

	case *Variable:
		return e.Name.Lexeme

	case *Assign:
		return p.parenthesize("= "+e.Name.Lexeme, e.Value)

	case *Logical:
		return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)

	case *Call:
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
