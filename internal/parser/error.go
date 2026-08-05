package parser

import (
	"fmt"

	"github.com/LightningDev1/golox/internal/scanner"
)

type ParseError struct {
	Token   scanner.Token
	Message string
}

func NewParseError(token scanner.Token, message string) ParseError {
	return ParseError{Token: token, Message: message}
}

func (e ParseError) Error() string {
	where := fmt.Sprintf("'%s'", e.Token.Lexeme)
	if e.Token.Type == scanner.TOKEN_EOF {
		where = "end"
	}

	return fmt.Sprintf("[line %d] Error at %s: %s",
		e.Token.Line, where, e.Message)
}
