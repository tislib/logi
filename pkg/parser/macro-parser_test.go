package parser

import (
	"github.com/stretchr/testify/assert"
	"logi/pkg/ast"
	"strings"
	"testing"
)

func TestSyntaxMacro(t *testing.T) {
	tests := map[string]struct {
		input         string
		expected      *ast.MacroAst
		expectedError string
	}{
		"simple syntax macro": {
			input: `
				macro simple {
					kind Syntax
				}
			`,
			expected: &ast.MacroAst{
				Macros: []ast.Macro{
					{
						Name: "simple",
						Kind: ast.MacroKindSyntax,
					},
				},
			},
		},
		"multiple syntax macro": {
			input: `
				macro simple {
					kind Syntax
				}

				macro simple2 {
					kind Syntax
				}
			`,
			expected: &ast.MacroAst{
				Macros: []ast.Macro{
					{
						Name: "simple",
						Kind: ast.MacroKindSyntax,
					},
					{
						Name: "simple2",
						Kind: ast.MacroKindSyntax,
					},
				},
			},
		},
		"syntax macro with simple syntax": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						Hello
						Hello2 Hello3
					}
				}
			`,
			expected: &ast.MacroAst{
				Macros: []ast.Macro{
					{
						Name: "simple",
						Kind: ast.MacroKindSyntax,
						Syntax: ast.Syntax{
							Statements: []ast.SyntaxStatement{
								{
									Elements: []ast.SyntaxStatementElement{
										{
											Kind: ast.SyntaxStatementElementKindKeyword,
											KeywordDef: &ast.SyntaxStatementElementKeywordDef{
												Name: "Hello",
											},
										},
									},
								}, {
									Elements: []ast.SyntaxStatementElement{
										{
											Kind: ast.SyntaxStatementElementKindKeyword,
											KeywordDef: &ast.SyntaxStatementElementKeywordDef{
												Name: "Hello2",
											},
										},
										{
											Kind: ast.SyntaxStatementElementKindKeyword,
											KeywordDef: &ast.SyntaxStatementElementKeywordDef{
												Name: "Hello3",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"fail macro if kind is missing": {
			input: `
				macro simple {
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is missing": {
			input: `
				macro {
					kind Syntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format": {
			input: `
				macro simple!{
					kind Syntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format[2]": {
			input: `
				macro simple simple{
					kind Syntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format[3]": {
			input: `
				macro SimPlEE{
					kind Syntax
				}
			`,
			expectedError: "failed to convert base macro: unexpected name value: SimPlEE",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseMacroContent(tt.input)

			if tt.expectedError != "" {
				if err == nil {
					assert.Fail(t, "expected error, got nil")
					return
				}
				if strings.Contains(tt.expectedError, err.Error()) {
					assert.Fail(t, "expected error %q, got %q", tt.expectedError, err.Error())
				}

				if got != nil {
					assert.Fail(t, "expected nil, got %v", got)
				}
				return
			} else {
				if err != nil {
					assert.Fail(t, "unexpected error: %s", err)
				}

				if got == nil {
					assert.Fail(t, "expected non-nil, got nil")
				}

				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
