package macro

import "fmt"

type Error struct {
	Line   int
	Column int
	At     string
	Msg    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("syntax error at or near \"%s\" at line %d column %d: %s", e.At, e.Line, e.Column, e.Msg)
}

func newError(line, column int, at, msg string) *Error {
	return &Error{
		Line:   line,
		Column: column,
		At:     at,
		Msg:    msg,
	}
}
