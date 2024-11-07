package macro

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tislib/logi/pkg/ast/common"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
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
					kind Syntax
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
					kind Syntax
				}

				macro simple2 {
					kind Syntax
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
					kind Syntax
					
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
		"syntax macro with variable keyword statement": {
			input: `
				macro simple {
					kind Syntax
					
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
												Type: common.TypeDefinition{
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
												Type: common.TypeDefinition{
													Name: "Type",
													SubTypes: []common.TypeDefinition{
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
														Type: common.TypeDefinition{
															Name: "string",
														},
													},
													{
														Name: "password",
														Type: common.TypeDefinition{
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
														Type: common.TypeDefinition{
															Name: "Type",
															SubTypes: []common.TypeDefinition{
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
		"syntax macro with code block and expression statement": {
			input: `
				macro simple {
					kind Syntax
					
					syntax {
						hello (...[<args Type<string>>]) { code }
						hello (...[<args Type<string>>]) { expr }
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
														Type: common.TypeDefinition{
															Name: "Type",
															SubTypes: []common.TypeDefinition{
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
														Type: common.TypeDefinition{
															Name: "Type",
															SubTypes: []common.TypeDefinition{
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
											Kind:            astMacro.SyntaxStatementElementKindExpressionBlock,
											ExpressionBlock: &astMacro.SyntaxStatementElementExpressionBlock{},
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
					kind Syntax
					
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
														Type: common.TypeDefinition{
															Name: "bool",
														},
													},
													{
														Name: "default",
														Type: common.TypeDefinition{Name: "string"},
													},
													{
														Name: "number",
														Type: common.TypeDefinition{Name: "float"},
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
		"syntax macro with types definition": {
			input: `
				macro simple {
					kind Syntax

					types {
						LatLong <lat float> <long float> 
					}

					syntax {
						Location <value LatLong>
					}
				}
`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,

						Types: astMacro.Types{
							Types: []astMacro.TypeStatement{
								{
									Name: "LatLong",
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "lat",
												Type: common.TypeDefinition{
													Name: "float",
												},
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "long",
												Type: common.TypeDefinition{
													Name: "float",
												},
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
												Name: "Location",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "value",
												Type: common.TypeDefinition{
													Name: "LatLong",
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
		"simple syntax with Or symbol": {
			input: `
				macro simple {
					kind Syntax

					syntax {
						(Hello | World) <name string>
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
											Kind: astMacro.SyntaxStatementElementKindCombination,
											Combination: &astMacro.SyntaxStatementElementCombination{
												Elements: []astMacro.SyntaxStatementElement{
													{
														Kind: astMacro.SyntaxStatementElementKindKeyword,
														KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
															Name: "Hello",
														},
													},
													{
														Kind: astMacro.SyntaxStatementElementKindKeyword,
														KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
															Name: "World",
														},
													},
												},
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "name",
												Type: common.TypeDefinition{
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
		"simple syntax with nested structure": {
			input: `
				macro simple {
					kind Syntax

					syntax {
						Auth {
							Username <username string>
							Password <password string>
						}
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
												Name: "Auth",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindStructure,
											Structure: &astMacro.SyntaxStatementElementStructure{
												Statements: []astMacro.SyntaxStatement{
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindKeyword,
																KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																	Name: "Username",
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																	Name: "username",
																	Type: common.TypeDefinition{
																		Name: "string",
																	},
																},
															},
														},
													},
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindKeyword,
																KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																	Name: "Password",
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																	Name: "password",
																	Type: common.TypeDefinition{
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
		},
		"simple syntax with double nested structure": {
			input: `
				macro simple {
					kind Syntax

					syntax {
						Auth {
							Credentials {
								Username <username string>
								Password <password string>
							}
						}
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
												Name: "Auth",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindStructure,
											Structure: &astMacro.SyntaxStatementElementStructure{
												Statements: []astMacro.SyntaxStatement{
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindKeyword,
																KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																	Name: "Credentials",
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindStructure,
																Structure: &astMacro.SyntaxStatementElementStructure{
																	Statements: []astMacro.SyntaxStatement{
																		{
																			Elements: []astMacro.SyntaxStatementElement{
																				{
																					Kind: astMacro.SyntaxStatementElementKindKeyword,
																					KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																						Name: "Username",
																					},
																				},
																				{
																					Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																					VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																						Name: "username",
																						Type: common.TypeDefinition{
																							Name: "string",
																						},
																					},
																				},
																			},
																		},
																		{
																			Elements: []astMacro.SyntaxStatementElement{
																				{
																					Kind: astMacro.SyntaxStatementElementKindKeyword,
																					KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																						Name: "Password",
																					},
																				},
																				{
																					Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																					VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																						Name: "password",
																						Type: common.TypeDefinition{
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
							},
						},
					},
				},
			},
		},
		"simple syntax with nested structure with or": {
			input: `
				macro simple {
					kind Syntax

					syntax {
						Auth {
							(Hello | World) <name string>
							Username <username string>
							Password <password string>
						}
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
												Name: "Auth",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindStructure,
											Structure: &astMacro.SyntaxStatementElementStructure{
												Statements: []astMacro.SyntaxStatement{
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindCombination,
																Combination: &astMacro.SyntaxStatementElementCombination{
																	Elements: []astMacro.SyntaxStatementElement{
																		{
																			Kind: astMacro.SyntaxStatementElementKindKeyword,
																			KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																				Name: "Hello",
																			},
																		},
																		{
																			Kind: astMacro.SyntaxStatementElementKindKeyword,
																			KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																				Name: "World",
																			},
																		},
																	},
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																	Name: "name",
																	Type: common.TypeDefinition{
																		Name: "string",
																	},
																},
															},
														},
													},
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindKeyword,
																KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																	Name: "Username",
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																	Name: "username",
																	Type: common.TypeDefinition{
																		Name: "string",
																	},
																},
															},
														},
													},
													{
														Elements: []astMacro.SyntaxStatementElement{
															{
																Kind: astMacro.SyntaxStatementElementKindKeyword,
																KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
																	Name: "Password",
																},
															},
															{
																Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
																VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
																	Name: "password",
																	Type: common.TypeDefinition{
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
		},
		"simple syntax with type reference": {
			input: `
				macro simple {
					kind Syntax

					types {
						World <value string>
					}

					syntax {
						Hello <World>
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "simple",
						Kind: astMacro.KindSyntax,
						Types: astMacro.Types{
							Types: []astMacro.TypeStatement{
								{
									Name: "World",
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
											VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
												Name: "value",
												Type: common.TypeDefinition{
													Name: "string",
												},
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
										{
											Kind: astMacro.SyntaxStatementElementKindTypeReference,
											TypeReference: &astMacro.SyntaxStatementElementTypeReference{
												Name: "World",
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
					assert.Fail(t, "unexpected error: %w", err)
				}

				if got == nil {
					assert.Fail(t, "expected non-nil, got nil")
				}

				expectedJson, _ := json.MarshalIndent(tt.expected, "", "  ")
				gotJson, _ := json.MarshalIndent(got, "", "  ")

				assert.Equal(t, string(expectedJson), string(gotJson))
			}
		})
	}
}
