package common

import "fmt"

type SourceLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (s SourceLocation) String() string {
	return fmt.Sprintf("L%d:%d", s.Line, s.Column)
}

type SourceMap map[string]SourceLocation

type SourceFile struct {
	Url string `json:"url"`
}
