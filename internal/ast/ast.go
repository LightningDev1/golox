package ast

import "github.com/LightningDev1/golox/internal/scanner"

type Expr interface {
	exprNode()
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

type Variable struct {
	Name scanner.Token
}

type Assign struct {
	Name  scanner.Token
	Value Expr
}

type Logical struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type Call struct {
	Callee    Expr
	Paren     scanner.Token
	Arguments []Expr
}

func (*Binary) exprNode()   {}
func (*Grouping) exprNode() {}
func (*Literal) exprNode()  {}
func (*Unary) exprNode()    {}
func (*Variable) exprNode() {}
func (*Assign) exprNode()   {}
func (*Logical) exprNode()  {}
func (*Call) exprNode()     {}
