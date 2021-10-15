package lexer

import (
	"github.com/ZephroC/go-interpreter/token"
	"unicode"
)

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	tok.Literal = string(l.ch)
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = "=="
		}
		tok.Type = token.ASSIGN
	case ';':
		tok.Type = token.SEMICOLON
	case '(':
		tok.Type = token.LPAREN
	case ')':
		tok.Type = token.RPAREN
	case '{':
		tok.Type = token.LBRACE
	case '}':
		tok.Type = token.RBRACE
	case ',':
		tok.Type = token.COMMA
	case '+':
		tok.Type =token.PLUS
	case '-':
		tok.Type = token.MINUS
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.NOT_EQ
			tok.Literal = "!="
		}
		tok.Type = token.BANG
	case '*':
		tok.Type = token.ASTERISK
	case '/':
		tok.Type = token.SLASH
	case '<':
		tok.Type = token.LESS_THAN
	case '>':
		tok.Type = token.GREATER_THAN
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else  {
			tok.Type = token.ILLEGAL
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}