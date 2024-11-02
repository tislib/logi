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
		"syntax macro with simple syntax and definition": {
			input: `
				macro simple {
					kind Syntax
					
					definition {
						Hello2 Hello3
					}
					
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
						Definition: ast.Definition{
							Statements: []ast.SyntaxStatement{
								{
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
		"syntax macro with variable keyword statement": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						hello <userName string>
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
												Name: "hello",
											},
										},
										{
											Kind: ast.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &ast.SyntaxStatementElementVariableKeyword{
												Name: "userName",
												Type: ast.TypeDefinition{
													Name: "string",
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
		},
		"syntax macro with variable keyword with generic type statement": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						hello <userName Type<string>>
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
												Name: "hello",
											},
										},
										{
											Kind: ast.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &ast.SyntaxStatementElementVariableKeyword{
												Name: "userName",
												Type: ast.TypeDefinition{
													Name: "Type",
													SubTypes: []ast.TypeDefinition{
														{
															Name: "string",
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
				},
			},
		},
		"syntax macro with parameter list statement": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						hello (<userName string>, <password string>)
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
												Name: "hello",
											},
										},
										{
											Kind: ast.SyntaxStatementElementKindParameterList,
											ParameterList: &ast.SyntaxStatementElementParameterList{
												Parameters: []ast.SyntaxStatementElementParameter{
													{
														Name: "userName",
														Type: ast.TypeDefinition{
															Name: "string",
														},
													},
													{
														Name: "password",
														Type: ast.TypeDefinition{
															Name: "string",
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
				},
			},
		},
		"syntax macro with argument list statement": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						hello (...[<args Type<string>>])
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
												Name: "hello",
											},
										},
										{
											Kind: ast.SyntaxStatementElementKindArgumentList,
											ArgumentList: &ast.SyntaxStatementElementArgumentList{
												VarArgs: true,
												Arguments: []ast.SyntaxStatementElementArgument{
													{
														Name: "args",
														Type: ast.TypeDefinition{
															Name: "Type",
															SubTypes: []ast.TypeDefinition{
																{
																	Name: "string",
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
