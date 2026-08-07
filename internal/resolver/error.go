package resolver

import (
	"fmt"

	"github.com/LightningDev1/golox/internal/scanner"
)

type ResolveError struct {
	Token   scanner.Token
	Message string
}

func NewResolveError(token scanner.Token, message string) ResolveError {
	return ResolveError{Token: token, Message: message}
}

func (e ResolveError) Error() string {
	return fmt.Sprintf("[line %d] Error at %s: %s",
		e.Token.Line, e.Token.Lexeme, e.Message)
}
