package interpreter

import (
	"fmt"

	"github.com/LightningDev1/golox/internal/scanner"
)

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func NewRuntimeError(token scanner.Token, message string) RuntimeError {
	return RuntimeError{Token: token, Message: message}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at %s: %s",
		e.Token.Line, e.Token.Lexeme, e.Message)
}
