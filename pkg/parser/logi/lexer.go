package logi

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
)

type logiLexer struct {
	buf   *bufio.Reader
	Err   error
	debug bool
}

func newLogiLexer(r io.Reader, debug bool) *logiLexer {
	return &logiLexer{
		buf:   bufio.NewReader(r),
		debug: debug,
	}
}

func (sc *logiLexer) Error(s string) {
	sc.Err = errors.New(s)
}

func (sc *logiLexer) Lex(lval *yySymType) int {
	res := sc.lex(lval)
	if sc.debug {
		log.Printf("lex: %d, %v\n", res, lval)
	}
	return res
}

func (s *logiLexer) lex(lval *yySymType) int {
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
		case '+':
			return Plus
		case '-':
			return Minus
		case '*':
			return Star
		case '/':
			return Slash
		case '%':
			return Percent
		case '!':
			return Exclamation
		case '&':
			return And
		case '|':
			return Or
		case '^':
			return Xor
		case '.':
			return Dot
		default:
			if isAlpha(r) {
				s.unread()
				var identifier = s.scanIdentifier()

				switch identifier {
				case "if":
					return IfKeyword
				case "var":
					return VarKeyword
				//case "for":
				//	return ForKeyword
				case "return":
					return ReturnKeyword
				case "switch":
					return SwitchKeyword
				case "case":
					return CaseKeyword
				case "logi":
					return LogiKeyword
				case "definition":
					return DefinitionKeyword
				case "syntax":
					return SyntaxKeyword
				case "false":
					lval.bool = false
					return token_bool
				case "true":
					lval.bool = true
					return token_bool
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

func (s *logiLexer) scanStr() string {
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

func (s *logiLexer) scanNumber() interface{} {
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

func (s *logiLexer) scanIdentifier() string {
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

func (s *logiLexer) read() rune {
	ch, _, _ := s.buf.ReadRune()
	return ch
}

func (s *logiLexer) unread() { _ = s.buf.UnreadRune() }
