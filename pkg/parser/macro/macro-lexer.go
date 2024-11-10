package macro

import (
	"errors"
	"github.com/tislib/logi/pkg/parser/lexer"
	"io"
	"log"
	"strconv"
	"strings"
)

type macroLexer struct {
	lexer lexer.Lexer
	Err   error
	debug bool
}

func newMacroLexer(r io.Reader, debug bool) *macroLexer {
	return &macroLexer{
		lexer: lexer.NewLexer(lexer.LexerConfig{
			HandleComments: true,
			Tokens:         tokens(),
		}, r, debug),
		debug: debug,
	}
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
			Id:     CodeBlock,
			Equals: "{ code }",
		},
		{
			Id:     ExpressionBlock,
			Equals: "{ expr }",
		},
		{
			Id:     TypesKeyword,
			Equals: "types",
		},
		{
			Id:     SyntaxKeyword,
			Equals: "syntax",
		},
		{
			Id:     MacroKeyword,
			Equals: "macro",
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
			Id:     Dash,
			Equals: "-",
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
			Id:       token_string,
			IsString: true,
		},
		{
			Id:           token_identifier,
			IsIdentifier: true,
		},
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
	token, err := s.lexer.Next()

	if err != nil {
		if errors.Is(err, lexer.ErrEOF) {
			return 0
		}
		s.Err = err
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

func (sc *macroLexer) GetReadString() any {
	return sc.lexer.GetReadString()
}
