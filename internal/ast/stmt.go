package ast

import "github.com/LightningDev1/golox/internal/scanner"

type Stmt interface {
	stmtNode()
}

type ExpressionStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

type BlockStmt struct {
	Statements []Stmt
}

func (*ExpressionStmt) stmtNode() {}
func (*PrintStmt) stmtNode()      {}
func (*VarStmt) stmtNode()        {}
func (*BlockStmt) stmtNode()      {}
