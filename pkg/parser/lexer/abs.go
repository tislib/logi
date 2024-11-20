package lexer

type TokenConfig struct {
	Id                    int
	Label                 string
	StartsWith            string
	EndsWith              string
	Equals                string
	EqualOneOf            []string
	EqualsCaseInsensitive string
	IsIdentifier          bool
	IsAlphaNum            bool
	IsAlpha               bool
	IsDigit               bool
	IsEol                 bool
	IsString              bool
}

type Token struct {
	Id    int
	Value string
}

type LexerConfig struct {
	HandleComments bool
	Tokens         []TokenConfig
}

type Union struct {
	bool   bool
	number interface{}
	string string
}

type Location struct {
	Line   int
	Column int
}

type Lexer interface {
	Next() (Token, error)
	GetReadString() any
	GetLastToken() Token
	GetLastLocation() Location
}
