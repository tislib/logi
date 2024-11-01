package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Scanner struct {
	buf   *bufio.Reader
	data  interface{}
	err   error
	debug bool
}

func NewScanner(r io.Reader, debug bool) *Scanner {
	return &Scanner{
		buf:   bufio.NewReader(r),
		debug: debug,
	}
}

func (sc *Scanner) Error(s string) {
	sc.err = errors.New(s)
}

func (sc *Scanner) Reduced(rule, state int, lval *yySymType) bool {
	if sc.debug {
		fmt.Printf("rule: %v; state %v; lval: %v\n", rule, state, lval)
	}
	return false
}

func (s *Scanner) Lex(lval *yySymType) int {
	return s.lex(lval)
}

func (s *Scanner) lex(lval *yySymType) int {
	for {
		r := s.read()
		if r == 0 {
			log.Print("EOF")
			return 0
		}
		if isWhitespace(r) {
			continue
		}

		if isDigit(r) {
			s.unread()
			lval.val = s.scanNumber()
			log.Println("number: ", lval.val)
			return NUMBER
		}

		switch r {
		case '^':
			log.Println("POW")
			return POW
		case '-':
			log.Println("MINUS")
			return MINUS
		case '+':
			log.Println("PLUS")
			return PLUS
		case '/':
			log.Println("DIVIDE")
			return DIVIDE
		case '*':
			log.Println("TIMES")
			return TIMES
		case '(':
			log.Println("LP")
			return LP
		case ')':
			log.Println("RP")
			return RP
		default:
			log.Println("Error: ", r)
			s.err = errors.New("error")
			return 0
		}
	}
}

func (s *Scanner) scanTrue() bool {
	t := []rune{'t', 'r', 'u', 'e'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("true is error")
			return false
		}
	}
	return true
}

func (s *Scanner) scanFalse() bool {
	t := []rune{'f', 'a', 'l', 's', 'e'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("false is error")
			return false
		}
	}
	return true
}

func (s *Scanner) scanNull() bool {
	t := []rune{'n', 'u', 'l', 'l'}
	for _, i := range t {
		r := s.read()
		if r != i {
			s.err = errors.New("null is error")
			return false
		}
	}
	return true
}

func (s *Scanner) scanStr() string {
	var str []rune
	//begin with ", end with "
	r := s.read()
	if r != '"' {
		os.Exit(1)
	}

	for {
		r := s.read()
		if r == '"' || r == 1 {
			break
		}
		str = append(str, r)
	}
	return string(str)
}

func (s *Scanner) scanNumber() interface{} {
	var number []rune
	var isFloat bool
	for {
		r := s.read()
		log.Println("nr: ", r)
		if r == '.' && len(number) > 0 && !isFloat {
			isFloat = true
			number = append(number, r)
			continue
		}

		if isWhitespace(r) || r == ',' || r == '}' || r == ']' {
			s.unread()
			break
		}
		if r == 0 {
			break
		}
		if !isDigit(r) {
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

func (s *Scanner) read() rune {
	ch, _, _ := s.buf.ReadRune()
	return ch
}

func (s *Scanner) unread() { _ = s.buf.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
