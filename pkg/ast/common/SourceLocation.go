package common

type SourceLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SourceFile struct {
	Url string `json:"url"`
}
