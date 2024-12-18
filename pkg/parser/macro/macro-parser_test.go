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
		"simple syntax with examples": {
			input: `
				macro simple {
					kind Syntax

					types {
						World <value string>
					}

					syntax {
						Hello <World> # Hello "World", Hello "World2"
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
									Examples: []string{
										`Hello "World"`,
										`Hello "World2"`,
									},
								},
							},
						},
					},
				},
			},
		},
		"simple syntax with examples2": {
			input: `
				macro simple {
					kind Syntax

					types {
						World <value string>
					}

					syntax {
						Hello <World> # [Hello "World", Hello "World2"]
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
									Examples: []string{
										`[Hello "World", Hello "World2"]`,
									},
								},
							},
						},
					},
				},
			},
		},
		"simple syntax with scopes": {
			input: `
				macro circuit {
					kind Syntax
				
					syntax {
						components { components }
						actions { command | handler }
					}
				
					scopes {
						components {
							<component Name> Led <pin int>
							<component Name> Button <pin int>
						}
						command {
							// Basic commands
							on(<component Name>)
							off(<component Name>)
							blink(<component Name>, <count int>, <seconds float>)
							wait(<seconds float>)
							brightness(<component Name>, <value float>)
							fade_in(<component Name>, <seconds float>)
							fade_out(<component Name>, <seconds float>)
							// Conditional commands
							if (<condition bool>) { command | handler }
							if (<condition bool>) { command | handler } else { command | handler }
						}
						handler {
							// Event handlers
							on_click(<component Name>) { command }
							on_click(<component Name>, <count int>) { command }
							on_press(<component Name>) { command }
							on_release(<component Name>) { command }
							while_held(<component Name>) { command }
						}
					}
				}
			`,
			expected: &astMacro.Ast{
				Macros: []astMacro.Macro{
					{
						Name: "circuit",
						Kind: astMacro.KindSyntax,
						Syntax: astMacro.Syntax{
							Statements: []astMacro.SyntaxStatement{
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "components",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindScope,
											ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
												Scopes: []string{"components"},
											},
										},
									},
								},
								{
									Elements: []astMacro.SyntaxStatementElement{
										{
											Kind: astMacro.SyntaxStatementElementKindKeyword,
											KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
												Name: "actions",
											},
										},
										{
											Kind: astMacro.SyntaxStatementElementKindScope,
											ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
												Scopes: []string{"command", "handler"},
											},
										},
									},
								},
							},
						},
						Scopes: astMacro.Scopes{
							Scopes: []astMacro.ScopeItem{
								{
									Name: "components",
									Statements: []astMacro.SyntaxStatement{
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
													VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
														Name: "component",
														Type: common.TypeDefinition{
															Name: "Name",
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "Led",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
													VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
														Name: "pin",
														Type: common.TypeDefinition{
															Name: "int",
														},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
													VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
														Name: "component",
														Type: common.TypeDefinition{
															Name: "Name",
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "Button",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindVariableKeyword,
													VariableKeyword: &astMacro.SyntaxStatementElementVariableKeyword{
														Name: "pin",
														Type: common.TypeDefinition{
															Name: "int",
														},
													},
												},
											},
										},
									},
								},
								{
									Name: "command",
									Statements: []astMacro.SyntaxStatement{
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "on",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
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
														Name: "off",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
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
														Name: "blink",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
															{
																Name: "count",
																Type: common.TypeDefinition{
																	Name: "int",
																},
															},
															{
																Name: "seconds",
																Type: common.TypeDefinition{
																	Name: "float",
																},
															},
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
														Name: "wait",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "seconds",
																Type: common.TypeDefinition{
																	Name: "float",
																},
															},
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
														Name: "brightness",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
															{
																Name: "value",
																Type: common.TypeDefinition{
																	Name: "float",
																},
															},
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
														Name: "fade_in",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
															{
																Name: "seconds",
																Type: common.TypeDefinition{
																	Name: "float",
																},
															},
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
														Name: "fade_out",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
															{
																Name: "seconds",
																Type: common.TypeDefinition{
																	Name: "float",
																},
															},
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
														Name: "if",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "condition",
																Type: common.TypeDefinition{
																	Name: "bool",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command", "handler"},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "if",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "condition",
																Type: common.TypeDefinition{
																	Name: "bool",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command", "handler"},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "else",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command", "handler"},
													},
												},
											},
										},
									},
								},
								{
									Name: "handler",
									Statements: []astMacro.SyntaxStatement{
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "on_click",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command"},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "on_click",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
															{
																Name: "count",
																Type: common.TypeDefinition{
																	Name: "int",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command"},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "on_press",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command"},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "on_release",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command"},
													},
												},
											},
										},
										{
											Elements: []astMacro.SyntaxStatementElement{
												{
													Kind: astMacro.SyntaxStatementElementKindKeyword,
													KeywordDef: &astMacro.SyntaxStatementElementKeywordDef{
														Name: "while_held",
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindParameterList,
													ParameterList: &astMacro.SyntaxStatementElementParameterList{
														Parameters: []astMacro.SyntaxStatementElementParameter{
															{
																Name: "component",
																Type: common.TypeDefinition{
																	Name: "Name",
																},
															},
														},
													},
												},
												{
													Kind: astMacro.SyntaxStatementElementKindScope,
													ScopeDef: &astMacro.SyntaxStatementElementScopeDef{
														Scopes: []string{"command"},
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
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseMacroContent(tt.input, false)

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
