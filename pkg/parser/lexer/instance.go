package lexer

import (
	"bufio"
	"io"
)

func NewLexer(config LexerConfig, r io.Reader, debug bool) Lexer {
	return &lexer{
		config: config,
		buf:    bufio.NewReader(r),
		debug:  debug,
	}
}
