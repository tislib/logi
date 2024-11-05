package macro

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
)

type macroLexer struct {
	buf     *bufio.Reader
	Err     error
	debug   bool
	readStr string
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
			return 0
		}
		if isEol(r) {
			return Eol
		}
		if isWhitespace(r) {
			continue
		}

		// handle comments
		if r == '/' {
			s.handleComments(r)
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
		case '|':
			return Or
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
				case "types":
					return TypesKeyword
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
	s.readStr += string(ch)
	return ch
}

func (s *macroLexer) unread() {
	_ = s.buf.UnreadRune()
	s.readStr = s.readStr[:len(s.readStr)-1]
}

func (sc *macroLexer) handleComments(r rune) {
	// single line comments
	if sc.read() == '/' {
		for {
			r = sc.read()
			if isEol(r) || r == 0 {
				sc.unread()
				return
			}
		}
	} else {
		sc.unread()
	}

	// multi line comments
	if sc.read() == '*' {
		for {
			r = sc.read()
			if r == '*' {
				if sc.read() == '/' {
					sc.unread()
					return
				}
			}
		}
	} else {
		sc.unread()
	}
}
