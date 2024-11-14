package logi

import (
	"errors"
	"fmt"
	"github.com/tislib/logi/pkg/parser/lexer"
	"io"
	"log"
	"strconv"
	"strings"
)

type logiLexer struct {
	lexer lexer.Lexer
	debug bool
	Err   error
}

func newLogiLexer(r io.Reader, debug bool) *logiLexer {
	return &logiLexer{
		lexer: lexer.NewLexer(lexer.LexerConfig{
			HandleComments: true,
			Tokens:         tokens(),
		}, r, debug),
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

func (s *logiLexer) Next() (lexer.Token, error) {
	return s.lexer.Next()
}

func (s *logiLexer) lex(lval *yySymType) int {
	token, err := s.lexer.Next()

	if err != nil {
		if errors.Is(err, lexer.ErrEOF) {
			return 0
		}
		s.Err = fmt.Errorf("lexer error: %w at [%s]", err, s.lexer.GetReadString())
		return 0
	}

	switch token.Id {
	case token_number:
		if strings.Contains(token.Value, ".") {
			number, err := strconv.ParseFloat(token.Value, 64)

			if err != nil {
				panic(err)
			}

			lval.number = number
		} else {
			number, err := strconv.Atoi(token.Value)

			if err != nil {
				panic(err)
			}

			lval.number = number
		}
	case token_string:
		lval.string = token.Value
	case token_bool:
		lval.bool = token.Value == "true"
	case token_identifier:
		lval.string = token.Value
	}

	return token.Id
}

func tokens() []lexer.TokenConfig {

	return []lexer.TokenConfig{
		{
			Id:      token_number,
			IsDigit: true,
		},
		{
			Id:         token_bool,
			EqualOneOf: []string{"true", "false"},
		},
		{
			Id:     BracketOpen,
			Equals: "[",
		},
		{
			Id:     BracketClose,
			Equals: "]",
		},
		{
			Id:     BraceOpen,
			Equals: "{",
		},
		{
			Id:     BraceClose,
			Equals: "}",
		},
		{
			Id:     Comma,
			Equals: ",",
		},
		{
			Id:     Colon,
			Equals: ":",
		},
		{
			Id:     Semicolon,
			Equals: ";",
		},
		{
			Id:     ParenOpen,
			Equals: "(",
		},
		{
			Id:     ParenClose,
			Equals: ")",
		},
		{
			Id:    Eol,
			IsEol: true,
		},
		{
			Id:     Equal,
			Equals: "=",
		},
		{
			Id:     GreaterThan,
			Equals: ">",
		},
		{
			Id:     LessThan,
			Equals: "<",
		},
		{
			Id:     Dot,
			Equals: ".",
		},
		{
			Id:     Arrow,
			Equals: "->",
		},
		{
			Id:     Or,
			Equals: "|",
		},
		{
			Id:     And,
			Equals: "&",
		},
		{
			Id:     Exclamation,
			Equals: "!",
		},
		{
			Id:     Plus,
			Equals: "+",
		},
		{
			Id:     Minus,
			Equals: "-",
		},
		{
			Id:     Star,
			Equals: "*",
		},
		{
			Id:     Slash,
			Equals: "/",
		},
		{
			Id:     Percent,
			Equals: "%",
		},
		{
			Id:     Xor,
			Equals: "^",
		},
		{
			Id:     FuncKeyword,
			Equals: "func",
		},
		{
			Id:     ElseKeyword,
			Equals: "else",
		},
		{
			Id:     VarKeyword,
			Equals: "var",
		},
		{
			Id:     IfKeyword,
			Equals: "if",
		},
		{
			Id:     ReturnKeyword,
			Equals: "return",
		},
		{
			Id:     SwitchKeyword,
			Equals: "switch",
		},
		{
			Id:     CaseKeyword,
			Equals: "case",
		},
		{
			Id:       token_string,
			IsString: true,
		},
		{
			Id:           token_identifier,
			IsIdentifier: true,
		},
	}
}

func (sc *logiLexer) GetReadString() any {
	return sc.lexer.GetReadString()
}
