package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type macroLexer struct {
	buf   *bufio.Reader
	Err   error
	debug bool
}

func newMacroLexer(r io.Reader, debug bool) *macroLexer {
	return &macroLexer{
		buf:   bufio.NewReader(r),
		debug: debug,
	}
}

func (sc *macroLexer) Error(s string) {
	sc.Err = errors.New(s)
}

func (sc *macroLexer) Lex(lval *yySymType) int {
	res := sc.lex(lval)
	if sc.debug {
		log.Printf("lex: %d, %v\n", res, lval)
	}
	return res
}

func (s *macroLexer) lex(lval *yySymType) int {
	for {
		r := s.read()
		if r == 0 { // EOF
			log.Print("EOF")
			return 0
		}
		if isEol(r) {
			return Eol
		}
		if isWhitespace(r) {
			continue
		}

		if isDigit(r) {
			s.unread()
			lval.number = s.scanNumber()
			return token_number
		}

		if r == '"' {
			s.unread()
			lval.string = s.scanStr()
			return token_string
		}

		switch r {
		case 't':
			if s.scanTrue() {
				lval.bool = true
				return token_bool
			}
		case 'f':
			if s.scanFalse() {
				lval.bool = false
				return token_bool
			}
		case '{':
			return BraceOpen
		case '}':
			return BraceClose
		case '[':
			return BracketOpen
		case ']':
			return BracketClose
		case '(':
			return ParenOpen
		case ')':
			return ParenClose
		case ',':
			return Comma
		case ':':
			return Colon
		case '=':
			return Equal
		case '>':
			return GreaterThan
		case '<':
			return LessThan
		case '-':
			next := s.read()
			if next == '>' {
				return Arrow
			}
			s.unread()
			return Dash
		case '.':
			return Dot
		default:
			if isAlpha(r) {
				s.unread()
				var identifier = s.scanIdentifier()

				switch identifier {
				case "macro":
					return MacroKeyword
				case "definition":
					return DefinitionKeyword
				case "syntax":
					return SyntaxKeyword
				default:
					lval.string = identifier
					return token_identifier
				}
			}

			log.Println("Error: Unrecognized character ", r)
			s.Err = errors.New("error: unrecognized character")
			return 0
		}
	}
}

func (s *macroLexer) scanTrue() bool {
	return s.matchWord("true")
}

func (s *macroLexer) scanFalse() bool {
	return s.matchWord("false")
}

func (s *macroLexer) scanNull() bool {
	return s.matchWord("null")
}

func (s *macroLexer) matchWord(word string) bool {
	for _, ch := range word {
		if s.read() != ch {
			s.Err = fmt.Errorf("expected %s", word)
			return false
		}
	}
	return true
}

func (s *macroLexer) scanStr() string {
	var str []rune
	if s.read() != '"' {
		return ""
	}
	for {
		r := s.read()
		if r == '"' || r == 0 {
			break
		}
		str = append(str, r)
	}
	return string(str)
}

func (s *macroLexer) scanNumber() interface{} {
	var number []rune
	var isFloat bool
	for {
		r := s.read()
		if r == '.' && len(number) > 0 && !isFloat {
			isFloat = true
			number = append(number, r)
			continue
		}

		if isWhitespace(r) || r == ',' || r == '}' || r == ']' {
			s.unread()
			break
		}
		if r == 0 || !isDigit(r) {
			s.unread()
			break
		}
		number = append(number, r)
	}
	if isFloat {
		f, _ := strconv.ParseFloat(string(number), 64)
		return f
	}
	i, _ := strconv.Atoi(string(number))
	return i
}

func (s *macroLexer) scanIdentifier() string {
	var identifier []rune
	for {
		r := s.read()
		if !isAlphaNum(r) {
			s.unread()
			break
		}
		identifier = append(identifier, r)
	}
	return string(identifier)
}

func (s *macroLexer) read() rune {
	ch, _, _ := s.buf.ReadRune()
	return ch
}

func (s *macroLexer) unread() { _ = s.buf.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isEol(ch rune) bool        { return ch == '\n' || ch == '\r' }
func isDigit(r rune) bool       { return r >= '0' && r <= '9' }
func isAlpha(r rune) bool       { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') }
func isAlphaNum(r rune) bool    { return isAlpha(r) || isDigit(r) }
