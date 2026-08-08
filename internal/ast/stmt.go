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

type ReturnStmt struct {
	Keyword scanner.Token
	Value   Expr
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type FunctionStmt struct {
	Name   scanner.Token
	Params []scanner.Token
	Body   []Stmt
}

type ClassStmt struct {
	Name    scanner.Token
	Methods []*FunctionStmt
}

func (*ExpressionStmt) stmtNode() {}
func (*PrintStmt) stmtNode()      {}
func (*ReturnStmt) stmtNode()     {}
func (*VarStmt) stmtNode()        {}
func (*BlockStmt) stmtNode()      {}
func (*IfStmt) stmtNode()         {}
func (*WhileStmt) stmtNode()      {}
func (*FunctionStmt) stmtNode()   {}
func (*ClassStmt) stmtNode()      {}
