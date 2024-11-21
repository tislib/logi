package common

type SourceLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SourceMap map[string]SourceLocation

type SourceFile struct {
	Url string `json:"url"`
}
