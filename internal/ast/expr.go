package ast

import "github.com/LightningDev1/golox/internal/scanner"

type Expr interface {
	exprNode()
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type GroupingExpr struct {
	Expression Expr
}

type LiteralExpr struct {
	Value any
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

type VariableExpr struct {
	Name scanner.Token
}

type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

type LogicalExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type CallExpr struct {
	Callee    Expr
	Paren     scanner.Token
	Arguments []Expr
}

type GetExpr struct {
	Object Expr
	Name   scanner.Token
}

type SetExpr struct {
	Object Expr
	Name   scanner.Token
	Value  Expr
}

type ThisExpr struct {
	Keyword scanner.Token
}

func (*BinaryExpr) exprNode()   {}
func (*GroupingExpr) exprNode() {}
func (*LiteralExpr) exprNode()  {}
func (*UnaryExpr) exprNode()    {}
func (*VariableExpr) exprNode() {}
func (*AssignExpr) exprNode()   {}
func (*LogicalExpr) exprNode()  {}
func (*CallExpr) exprNode()     {}
func (*GetExpr) exprNode()      {}
func (*SetExpr) exprNode()      {}
func (*ThisExpr) exprNode()     {}
