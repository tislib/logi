package macro

import (
	"github.com/stretchr/testify/assert"
	astMacro "logi/pkg/ast/macro"
	"strings"
	"testing"
)

func TestSyntaxMacro(t *testing.T) {
	tests := map[string]struct {
		input         string
		expected      *astMacro.Ast
		expectedError string
	}{
		"simple syntax macro": {
			input: `
				macro simple {
					kind MacroSyntax
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
					},
				},
			},
		},
		"multiple syntax macro": {
			input: `
				macro simple {
					kind MacroSyntax
				}

				macro simple2 {
					kind MacroSyntax
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
					},
					{
						Name: "simple2",
						Kind: astMacro.KindSyntax,
					},
				},
			},
		},
		"syntax macro with simple syntax": {
			input: `
				macro simple {
					kind MacroSyntax
					
					syntax {
						Hello
						Hello2 Hello3
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello",
											},
										},
									},
								}, {
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello2",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
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
					kind MacroSyntax
					
					definition {
						Hello2 Hello3
					}
					
					syntax {
						Hello
						Hello2 Hello3
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Definition: astMacro.Definition{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello2",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello3",
											},
										},
									},
								},
							},
						},
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello",
											},
										},
									},
								}, {
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "Hello2",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
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
					kind MacroSyntax
					
					syntax {
						hello <userName string>
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "userName",
												Type: astMacro.TypeDefinition{
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
					kind MacroSyntax
					
					syntax {
						hello <userName Type<string>>
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "userName",
												Type: astMacro.TypeDefinition{
													Name: "Type",
													SubTypes: []astMacro.TypeDefinition{
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
					kind MacroSyntax
					
					syntax {
						hello (<userName string>, <password string>)
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindParameterList,
											ParameterList: &astMacro.SyntaxStatementElementParameterList{
												Parameters: []astMacro.SyntaxStatementElementParameter{
													{
														Name: "userName",
														Type: astMacro.TypeDefinition{
															Name: "string",
														},
													},
													{
														Name: "password",
														Type: astMacro.TypeDefinition{
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
					kind MacroSyntax
					
					syntax {
						hello (...[<args Type<string>>])
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindArgumentList,
											ArgumentList: &astMacro.SyntaxStatementElementArgumentList{
												VarArgs: true,
												Arguments: []astMacro.SyntaxStatementElementArgument{
													{
														Name: "args",
														Type: astMacro.TypeDefinition{
															Name: "Type",
															SubTypes: []astMacro.TypeDefinition{
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
		"syntax macro with code block statement": {
			input: `
				macro simple {
					kind MacroSyntax
					
					syntax {
						hello (...[<args Type<string>>]) { }
						hello (...[<args Type<string>>]) { string }
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindArgumentList,
											ArgumentList: &astMacro.SyntaxStatementElementArgumentList{
												VarArgs: true,
												Arguments: []astMacro.SyntaxStatementElementArgument{
													{
														Name: "args",
														Type: astMacro.TypeDefinition{
															Name: "Type",
															SubTypes: []astMacro.TypeDefinition{
																{
																	Name: "string",
																},
															},
														},
													},
												},
											},
										},
										{
											Kind:      astMacro.SyntaxStatementElementKindCodeBlock,
											CodeBlock: &astMacro.SyntaxStatementElementCodeBlock{},
										},
									},
								},
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "hello",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindArgumentList,
											ArgumentList: &astMacro.SyntaxStatementElementArgumentList{
												VarArgs: true,
												Arguments: []astMacro.SyntaxStatementElementArgument{
													{
														Name: "args",
														Type: astMacro.TypeDefinition{
															Name: "Type",
															SubTypes: []astMacro.TypeDefinition{
																{
																	Name: "string",
																},
															},
														},
													},
												},
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindCodeBlock,
											CodeBlock: &astMacro.SyntaxStatementElementCodeBlock{
												ReturnType: astMacro.TypeDefinition{
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
		"syntax macro with attributes statement": {
			input: `
				macro simple {
					kind MacroSyntax
					
					syntax {
						details [required bool, default string, number float]
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "details",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindAttributeList,
											AttributeList: &astMacro.SyntaxStatementElementAttributeList{
												Attributes: []astMacro.SyntaxStatementElementAttribute{
													{
														Name: "required",
														Type: astMacro.TypeDefinition{
															Name: "bool",
														},
													},
													{
														Name: "default",
														Type: astMacro.TypeDefinition{Name: "string"},
													},
													{
														Name: "number",
														Type: astMacro.TypeDefinition{Name: "float"},
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
					kind MacroSyntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format": {
			input: `
				macro simple!{
					kind MacroSyntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format[2]": {
			input: `
				macro simple simple{
					kind MacroSyntax
				}
			`,
			expectedError: "syntax error",
		},
		"fail macro if name is in incorrect format[3]": {
			input: `
				macro SimPlEE{
					kind MacroSyntax
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
					assert.Fail(t, "unexpected error: %w", err)
				}

				if got == nil {
					assert.Fail(t, "expected non-nil, got nil")
				}

				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
