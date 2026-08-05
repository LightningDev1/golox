package scanner

import "fmt"

type ScanError struct {
	Line    int
	Message string
}

func NewScanError(line int, message string) ScanError {
	return ScanError{Line: line, Message: message}
}

func (e ScanError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}
