package lexer

import (
	"errors"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	test := map[string]struct {
		input          string
		config         LexerConfig
		expectedTokens []Token
	}{
		"alpha keywords": {
			input: "hello world",
			config: LexerConfig{
				Tokens: []TokenConfig{
					{
						Id:      1,
						IsAlpha: true,
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    1,
					Value: "hello",
				},
				{
					Id:    1,
					Value: "world",
				},
			},
		},
		"alpha numeric keywords": {
			input: "hello123 world321",
			config: LexerConfig{
				Tokens: []TokenConfig{
					{
						Id:         1,
						IsAlphaNum: true,
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    1,
					Value: "hello123",
				},
				{
					Id:    1,
					Value: "world321",
				},
			},
		},
		"digits": {
			input: "123 321.321 0.123 123.0 -123 -123.0 -0.123",
			config: LexerConfig{
				Tokens: []TokenConfig{
					{
						Id:      1,
						IsDigit: true,
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    1,
					Value: "123",
				},
				{
					Id:    1,
					Value: "321.321",
				},
				{
					Id:    1,
					Value: "0.123",
				},
				{
					Id:    1,
					Value: "123.0",
				},
				{
					Id:    1,
					Value: "-123",
				},
				{
					Id:    1,
					Value: "-123.0",
				},
				{
					Id:    1,
					Value: "-0.123",
				},
			},
		},
		"read combination": {
			input: "hello (123) world",
			config: LexerConfig{
				Tokens: []TokenConfig{
					{
						Id:      1,
						IsDigit: true,
					},
					{
						Id:      2,
						IsAlpha: true,
					},
					{
						Id:         3,
						IsAlphaNum: true,
					},
					{
						Id:     4,
						Equals: "(",
					},
					{
						Id:     5,
						Equals: ")",
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    2,
					Value: "hello",
				},
				{
					Id:    4,
					Value: "(",
				},
				{
					Id:    1,
					Value: "123",
				},
				{
					Id:    5,
					Value: ")",
				},
				{
					Id:    2,
					Value: "world",
				},
			},
		},
		"string and number": {
			input: `"hello" 123 "world"`,
			config: LexerConfig{
				Tokens: []TokenConfig{
					{
						Id:      1,
						IsDigit: true,
					},
					{
						Id:      2,
						IsAlpha: true,
					},
					{
						Id:         3,
						IsAlphaNum: true,
					},
					{
						Id:     4,
						Equals: "(",
					},
					{
						Id:     5,
						Equals: ")",
					},
					{
						Id:       6,
						IsString: true,
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    6,
					Value: "hello",
				},
				{
					Id:    1,
					Value: "123",
				},
				{
					Id:    6,
					Value: "world",
				},
			},
		},
		"comments": {
			input: `
123 321.321 // aaaaa
0.123 123.0 -123 
// 0.123 123.0 -123
/*
ml
ml
*/
-123.0 /* 3.3 */ -0.123
`,
			config: LexerConfig{
				HandleComments: true,
				Tokens: []TokenConfig{
					{
						Id:      1,
						IsDigit: true,
					},
				},
			},
			expectedTokens: []Token{
				{
					Id:    1,
					Value: "123",
				},
				{
					Id:    1,
					Value: "321.321",
				},
				{
					Id:    1,
					Value: "0.123",
				},
				{
					Id:    1,
					Value: "123.0",
				},
				{
					Id:    1,
					Value: "-123",
				},
				{
					Id:    1,
					Value: "-123.0",
				},
				{
					Id:    1,
					Value: "-0.123",
				},
			},
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			lexer := NewLexer(tt.config, strings.NewReader(tt.input), false)

			var i int
			for {
				token, err := lexer.Next()

				if err != nil {
					if errors.Is(err, ErrEOF) {
						if i < len(tt.expectedTokens) {
							t.Errorf("Expected token %d[%s], got EOF", tt.expectedTokens[i].Id, tt.expectedTokens[i].Value)
						}
					} else {
						t.Errorf("Unexpected error: %s; expected token: %d[%s]", err, tt.expectedTokens[i].Id, tt.expectedTokens[i].Value)
					}
					break
				}

				if i >= len(tt.expectedTokens) {
					t.Errorf("Expected EOF , got %d[%s]", token.Id, token.Value)
					return
				}

				if token.Id != tt.expectedTokens[i].Id {
					t.Errorf("Expected token %d, got %d", tt.expectedTokens[i].Id, token.Id)
				}

				if token.Value != tt.expectedTokens[i].Value {
					t.Errorf("Expected token %s, got %s", tt.expectedTokens[i].Value, token.Value)
				}

				if t.Failed() {
					break
				}

				i++
			}
		})
	}

}
