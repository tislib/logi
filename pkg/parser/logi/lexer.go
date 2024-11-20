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

type LogiLexer struct {
	lexer lexer.Lexer
	debug bool
	Err   error
}

func NewLogiLexer(r io.Reader, debug bool) *LogiLexer {
	return &LogiLexer{
		lexer: lexer.NewLexer(lexer.LexerConfig{
			HandleComments: true,
			Tokens:         Tokens(),
		}, r, debug),
		debug: debug,
	}
}

func (sc *LogiLexer) Error(s string) {
	sc.Err = errors.New(s)
}

func (sc *LogiLexer) Lex(lval *yySymType) int {
	res := sc.lex(lval)
	if sc.debug {
		log.Printf("lex: %d, %v\n", res, lval)
	}
	return res
}

func (s *LogiLexer) Next() (lexer.Token, error) {
	return s.lexer.Next()
}

func (s *LogiLexer) lex(lval *yySymType) int {
	token, err := s.lexer.Next()

	if err != nil {
		if errors.Is(err, lexer.ErrEOF) {
			return 0
		}
		s.Err = fmt.Errorf("lexer error: %w at [%s]", err, s.lexer.GetReadString())
		return 0
	}

	lval.token = token
	lval.location = s.lexer.GetLastLocation()

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

func Tokens() []lexer.TokenConfig {
	return []lexer.TokenConfig{
		{
			Id:      token_number,
			Label:   "number",
			IsDigit: true,
		},
		{
			Id:         token_bool,
			Label:      "bool",
			EqualOneOf: []string{"true", "false"},
		},
		{
			Id:     BracketOpen,
			Label:  "[",
			Equals: "[",
		},
		{
			Id:     BracketClose,
			Label:  "]",
			Equals: "]",
		},
		{
			Id:     BraceOpen,
			Label:  "{",
			Equals: "{",
		},
		{
			Id:     BraceClose,
			Label:  "}",
			Equals: "}",
		},
		{
			Id:     Comma,
			Label:  ",",
			Equals: ",",
		},
		{
			Id:     Colon,
			Label:  ":",
			Equals: ":",
		},
		{
			Id:     Semicolon,
			Label:  ";",
			Equals: ";",
		},
		{
			Id:     ParenOpen,
			Label:  "(",
			Equals: "(",
		},
		{
			Id:     ParenClose,
			Label:  ")",
			Equals: ")",
		},
		{
			Id:    Eol,
			Label: "Eol",
			IsEol: true,
		},
		{
			Id:     Equal,
			Label:  "=",
			Equals: "=",
		},
		{
			Id:     GreaterThan,
			Label:  ">",
			Equals: ">",
		},
		{
			Id:     LessThan,
			Label:  "<",
			Equals: "<",
		},
		{
			Id:     Dot,
			Label:  ".",
			Equals: ".",
		},
		{
			Id:     Arrow,
			Label:  "->",
			Equals: "->",
		},
		{
			Id:     Or,
			Label:  "|",
			Equals: "|",
		},
		{
			Id:     And,
			Label:  "&",
			Equals: "&",
		},
		{
			Id:     Exclamation,
			Label:  "!",
			Equals: "!",
		},
		{
			Id:     Plus,
			Label:  "+",
			Equals: "+",
		},
		{
			Id:     Minus,
			Label:  "-",
			Equals: "-",
		},
		{
			Id:     Star,
			Label:  "*",
			Equals: "*",
		},
		{
			Id:     Slash,
			Label:  "/",
			Equals: "/",
		},
		{
			Id:     Percent,
			Label:  "%",
			Equals: "%",
		},
		{
			Id:     Xor,
			Label:  "^",
			Equals: "^",
		},
		{
			Id:     FuncKeyword,
			Label:  "func",
			Equals: "func",
		},
		{
			Id:     ElseKeyword,
			Label:  "else",
			Equals: "else",
		},
		{
			Id:     VarKeyword,
			Label:  "var",
			Equals: "var",
		},
		{
			Id:     IfKeyword,
			Label:  "if",
			Equals: "if",
		},
		{
			Id:     ReturnKeyword,
			Label:  "return",
			Equals: "return",
		},
		{
			Id:     SwitchKeyword,
			Label:  "switch",
			Equals: "switch",
		},
		{
			Id:     CaseKeyword,
			Label:  "case",
			Equals: "case",
		},
		{
			Id:       token_string,
			Label:    "string",
			IsString: true,
		},
		{
			Id:           token_identifier,
			Label:        "identifier",
			IsIdentifier: true,
		},
	}
}

func (sc *LogiLexer) GetReadString() any {
	return sc.lexer.GetReadString()
}

func (sc *LogiLexer) GetLastLocation() lexer.Location {
	return sc.lexer.GetLastLocation()
}
