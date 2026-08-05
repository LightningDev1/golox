package scanner

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FUN
	TOKEN_FOR
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_EOF
)

var keywords = map[string]TokenType{
	"and":    TOKEN_AND,
	"class":  TOKEN_CLASS,
	"else":   TOKEN_ELSE,
	"false":  TOKEN_FALSE,
	"for":    TOKEN_FOR,
	"fun":    TOKEN_FUN,
	"if":     TOKEN_IF,
	"nil":    TOKEN_NIL,
	"or":     TOKEN_OR,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
	"super":  TOKEN_SUPER,
	"this":   TOKEN_THIS,
	"true":   TOKEN_TRUE,
	"var":    TOKEN_VAR,
	"while":  TOKEN_WHILE,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}
