package scanner

import "strconv"

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  nil,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, []ScanError) {
	var errs []ScanError

	for !s.isAtEnd() {
		s.start = s.current

		if err, ok := s.scanToken(); !ok {
			errs = append(errs, err)
		}
	}

	s.tokens = append(s.tokens, NewToken(TOKEN_EOF, "", nil, s.line))

	return s.tokens, errs
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() (ScanError, bool) {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(TOKEN_LEFT_PAREN)
	case ')':
		s.addToken(TOKEN_RIGHT_PAREN)
	case '{':
		s.addToken(TOKEN_LEFT_BRACE)
	case '}':
		s.addToken(TOKEN_RIGHT_BRACE)
	case ',':
		s.addToken(TOKEN_COMMA)
	case '.':
		s.addToken(TOKEN_DOT)
	case '-':
		s.addToken(TOKEN_MINUS)
	case '+':
		s.addToken(TOKEN_PLUS)
	case ';':
		s.addToken(TOKEN_SEMICOLON)
	case '*':
		s.addToken(TOKEN_STAR)

	case '!':
		if s.match('=') {
			s.addToken(TOKEN_BANG_EQUAL)
		} else {
			s.addToken(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(TOKEN_EQUAL_EQUAL)
		} else {
			s.addToken(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(TOKEN_LESS_EQUAL)
		} else {
			s.addToken(TOKEN_LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(TOKEN_GREATER_EQUAL)
		} else {
			s.addToken(TOKEN_GREATER)
		}

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TOKEN_SLASH)
		}

	case ' ', '\r', '\t':
		break

	case '\n':
		s.line++

	case '"':
		return s.string()

	default:
		if s.isDigit(c) {
			return s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			return NewScanError(s.line, "unexpected character"), false
		}
	}

	return ScanError{}, true
}

func (s *Scanner) advance() byte {
	val := s.source[s.current]
	s.current++
	return val
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s *Scanner) string() (ScanError, bool) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return NewScanError(s.line, "unterminated string"), false
	}

	// The closing ".
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(TOKEN_STRING, value)

	return ScanError{}, true
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) number() (ScanError, bool) {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	numberString := s.source[s.start:s.current]

	number, err := strconv.ParseFloat(numberString, 64)
	if err != nil {
		return NewScanError(s.line, "invalid number"), false
	}

	s.addLiteralToken(TOKEN_NUMBER, number)

	return ScanError{}, true
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = TOKEN_IDENTIFIER
	}

	s.addToken(tokenType)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addLiteralToken(tokenType, nil)
}

func (s *Scanner) addLiteralToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}
