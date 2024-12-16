package logi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"strings"
	"testing"
)

func TestParserFull(t *testing.T) {
	tests := map[string]struct {
		skipped       bool
		macroInput    string
		input         string
		expected      *logiAst.Ast
		expectedError string
	}{
		"simple parse": {
			macroInput: `
				macro entity {
					kind Syntax

					syntax {
						<propertyName Name> <propertyType Type> [primary bool, autoincrement bool, required bool, default string]
					}
				}
`,
			input: `
				entity User {
					id int <[primary, autoincrement]>
					name string <[required, default "John Doe"]>
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
						Statements: []logiAst.Statement{
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.StringValue("id"),
									},
									{
										Name:  "propertyType",
										Value: common.StringValue("int"),
									},
								},
								Attributes: []logiAst.Attribute{
									{
										Name: "primary",
									},
									{
										Name: "autoincrement",
									},
								},
							},
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.StringValue("name"),
									},
									{
										Name:  "propertyType",
										Value: common.StringValue("string"),
									},
								},
								Attributes: []logiAst.Attribute{
									{
										Name: "required",
									},
									{
										Name:  "default",
										Value: common.PointerValue(common.StringValue("John Doe")),
									},
								},
							},
						},
					},
				},
			},
		},
		"simple parse comments": {
			macroInput: `
				macro entity {
					kind Syntax // this is a comment
					// this is another comment
					syntax {
						/* this is a block comment */ <propertyName Name> <propertyType Type> [primary bool, autoincrement bool, required bool, default string]
					}
					/* this is 
						a block 
						comment */
				}
`,
			input: `
				entity User {
					id int <[primary, autoincrement]>
					name string <[required, default "John Doe"]>
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
					},
				},
			},
		},
		"simple parse combination": {
			macroInput: `
				macro entity {
					kind Syntax
					syntax {
						(Hello | World) <propertyName Name> <propertyType Type>
					}
				}
`,
			input: `
				entity User {
					Hello id int
					World name string
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
					},
				},
			},
		},
		"complex parse combination": {
			macroInput: `
				macro entity {
					kind Syntax
					
					types {
						ParamType1 <value1 string> <value2 string>
						ParamType2 <value3 int> <value4 int>
					}

					syntax {
						<propertyName Name> (<value ParamType1> | <value ParamType2>)
					}
				}
`,
			input: `
				entity User {
					param1 "11" "22"
					param2 1 2
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
					},
				},
			},
		},
		"parse nested structure": {
			macroInput: `
				macro simple {
					kind Syntax

					syntax {
						Main { code }
						Auth {
							Username <username string>
							Password <password string>
						}
					}
				}
`,
			input: `
				simple User1 {
					Main {
						return 123
					}
					Auth {
						Username "user1"
						Password "password1"
					}
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "simple",
						Name:      "User1",
					},
				},
			},
		},
		"simple type reference": {
			macroInput: `
				macro entity {
					kind Syntax

					types {
						World <value1 string> <value2 string> <value3 string>
					}

					syntax {
						Hello <World>
					}
				}`,
			input: `
				entity User {
					Hello "from the other side" "hello" "world"
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
					},
				},
			},
		},
		"simple type first parse": {
			macroInput: `
				macro entity {
					kind Syntax

					syntax {
						<propertyType Type> <propertyName Name> [primary bool, autoincrement bool, required bool, default string]
					}
				}
`,
			input: `
				entity User {
					int id <[primary, autoincrement]>
					string name <[required, default "John Doe"]>
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "entity",
						Name:      "User",
					},
				},
			},
		},
		"simple signature definition for interface": {
			macroInput: `
				macro interface {
					kind Syntax

					syntax {
						<methodName Name> (...[<args Type<string>>]) <returnType Type>
					}
				}
`,
			input: `
				interface UserService {
					createUser (name string, age int) User
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "interface",
						Name:      "UserService",
					},
				},
			},
		},
		"simple method call": {
			macroInput: `
				macro implementation {
					kind Syntax

					syntax {
						<methodName Name> (...[<args Type<string>>]) <returnType Type> { code }
					}
				}
`,
			input: `
				implementation UserServiceImpl {
					createUser (name string, age int) User {
						return 0
					}
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "implementation",
						Name:      "UserServiceImpl",
					},
				},
			},
		},
		"test simple type definition": {
			macroInput: `
				macro complexArray {
					kind Syntax
				
					types {
						ParamValue 	<value11 string> <value21 string>
					}
				
					syntax {
						Param1 <value ParamValue>
					}
				}`,
			input: `
				complexArray ComplexArray1 {
					Param1 "value1" "value2"
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "complexArray",
						Name:      "ComplexArray1",
					},
				},
			},
		},
		"test simple type array definition": {
			macroInput: `
				macro complexArray {
					kind Syntax
				
					types {
						ParamValue <value11 string> <value21 string>
					}
				
					syntax {
						Param1 <value array<ParamValue>>
					}
				}`,
			input: `
				complexArray ComplexArray1 {
					Param1 ["value1" "value2", "value3" "value4"]
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "complexArray",
						Name:      "ComplexArray1",
					},
				},
			},
		},
		"test simple type array definition with EOL inside": {
			macroInput: `
				macro complexArray {
					kind Syntax
				
					types {
						ParamValue <value11 string> <value21 string>
					}
				
					syntax {
						Param1 <value array<ParamValue>>
					}
				}`,
			input: `
				complexArray ComplexArray1 {
					Param1 [
						"value1" "value2", 
						"value3" "value4"
					]
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "complexArray",
						Name:      "ComplexArray1",
					},
				},
			},
		},
		"test backtest macro": {
			macroInput: `
				macro backtest {
					kind Syntax
				
					types {
						Indicator <indicatorName Name> (<period int>) as <alias Name>
					}
				
					syntax {
						InitialCapital <initialCapital int>
						StartTime <startTime string>
						EndTime <endTime string>
						Indicators <indicators array<Indicator>>
						Strategy { code }
					}
				}`,
			input: `
				backtest VariableHoldUntil4 {
					InitialCapital  10000
					StartTime       "2010-01-01"
					EndTime         "2010-12-31"
					Indicators       [sma(20) as sma20, sma(50) as sma50, sma(200) as sma200]
				
					Strategy {
						if (sma20 < sma50) {
							Buy("SPY", 100)
						}
					}
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "backtest",
						Name:      "VariableHoldUntil4",
					},
				},
			},
		},
		"circuit-test": {
			macroInput: `
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
							on_press(<component Name>, <count int>) { command }
							on_release(<component Name>, <count int>) { command }
							while_held(<component Name>) { command }
						}
					}
				}
`,
			input: `
				circuit simple1 {
					components {
						yellowLed Led 5
						redLed Led 6
						blueLed Led 13
						button1 Button 17
						button2 Button 17
					}
				
					actions {
						on(yellowLed)
						on(redLed)
				
						on_click(button1) {
							if (status(button2) == 'on') {
								on(blueLed)
								on(yellowLed)
								on(redLed)
							}
						}
				
						on_click(button2) {
							off(blueLed)
							on(yellowLed)
							on(redLed)
						}
					}
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "circuit",
						Name:      "simple1",
						Statements: []logiAst.Statement{
							{
								Command: "components",
								SubStatements: []logiAst.Statement{
									{
										Command: "Led",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("yellowLed"),
											},
											{
												Name:  "pin",
												Value: common.IntegerValue(5),
											},
										},
									},
									{
										Command: "Led",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("redLed"),
											},
											{
												Name:  "pin",
												Value: common.IntegerValue(6),
											},
										},
									},
									{
										Command: "Led",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("blueLed"),
											},
											{
												Name:  "pin",
												Value: common.IntegerValue(13),
											},
										},
									},
									{
										Command: "Button",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("button1"),
											},
											{
												Name:  "pin",
												Value: common.IntegerValue(17),
											},
										},
									},
									{
										Command: "Button",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("button2"),
											},
											{
												Name:  "pin",
												Value: common.IntegerValue(17),
											},
										},
									},
								},
							},
							{
								Command: "actions",
								SubStatements: []logiAst.Statement{
									{
										Command: "on",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("yellowLed"),
												Expression: &common.Expression{
													Kind: common.VariableKind,

													Variable: &common.Variable{
														Name: "yellowLed",
													},
												},
											},
										},
									},
									{
										Command: "on",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("redLed"),
												Expression: &common.Expression{
													Kind: common.VariableKind,

													Variable: &common.Variable{
														Name: "redLed",
													},
												},
											},
										},
									},
									{
										Command: "on_click",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("button1"),
												Expression: &common.Expression{
													Kind: common.VariableKind,

													Variable: &common.Variable{
														Name: "button1",
													},
												},
											},
										},
										SubStatements: []logiAst.Statement{
											{
												Command: "if",
												Parameters: []logiAst.Parameter{
													{
														Name: "condition",
														Value: common.MapValue(map[string]common.Value{
															"left": common.MapValue(map[string]common.Value{
																"name":      common.StringValue("status"),
																"arguments": common.ArrayValue(common.StringValue("button2")),
															}),
															"operator": common.StringValue("=="),
															"right":    common.StringValue("on"),
														}),
														Expression: &common.Expression{
															Kind: common.BinaryExprKind,
															BinaryExpr: &common.BinaryExpression{
																Left: &common.Expression{
																	Kind: common.FuncCallKind,
																	FuncCall: &common.FunctionCall{
																		Name: "status",
																		Arguments: []*common.Expression{
																			{
																				Kind: common.VariableKind,
																				Variable: &common.Variable{
																					Name: "button2",
																				},
																			},
																		},
																	},
																},
																Operator: "==",
																Right: &common.Expression{
																	Kind: common.LiteralKind,
																	Literal: &common.Literal{
																		Value: common.StringValue("on"),
																	},
																},
															},
														},
													},
												},
												SubStatements: []logiAst.Statement{
													{
														Command: "on",
														Parameters: []logiAst.Parameter{
															{
																Name:  "component",
																Value: common.StringValue("blueLed"),
																Expression: &common.Expression{
																	Kind: common.VariableKind,

																	Variable: &common.Variable{
																		Name: "blueLed",
																	},
																},
															},
														},
													},
													{
														Command: "on",
														Parameters: []logiAst.Parameter{
															{
																Name:  "component",
																Value: common.StringValue("yellowLed"),
																Expression: &common.Expression{
																	Kind: common.VariableKind,

																	Variable: &common.Variable{
																		Name: "yellowLed",
																	},
																},
															},
														},
													},
													{
														Command: "on",
														Parameters: []logiAst.Parameter{
															{
																Name:  "component",
																Value: common.StringValue("redLed"),
																Expression: &common.Expression{
																	Kind: common.VariableKind,

																	Variable: &common.Variable{
																		Name: "redLed",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									{
										Command: "on_click",
										Parameters: []logiAst.Parameter{
											{
												Name:  "component",
												Value: common.StringValue("button2"),
												Expression: &common.Expression{
													Kind: common.VariableKind,

													Variable: &common.Variable{
														Name: "button2",
													},
												},
											},
										},
										SubStatements: []logiAst.Statement{
											{
												Command: "off",
												Parameters: []logiAst.Parameter{
													{
														Name:  "component",
														Value: common.StringValue("blueLed"),
														Expression: &common.Expression{
															Kind: common.VariableKind,

															Variable: &common.Variable{
																Name: "blueLed",
															},
														},
													},
												},
											},
											{
												Command: "on",
												Parameters: []logiAst.Parameter{
													{
														Name:  "component",
														Value: common.StringValue("yellowLed"),
														Expression: &common.Expression{
															Kind: common.VariableKind,

															Variable: &common.Variable{
																Name: "yellowLed",
															},
														},
													},
												},
											},
											{
												Command: "on",
												Parameters: []logiAst.Parameter{
													{
														Name:  "component",
														Value: common.StringValue("redLed"),
														Expression: &common.Expression{
															Kind: common.VariableKind,

															Variable: &common.Variable{
																Name: "redLed",
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
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseFullWithMacro(tt.input, tt.macroInput, false)

			if tt.skipped {
				t.Skip()
				return
			}

			if got != nil && tt.expected != nil {
				if len(got.Definitions) == len(tt.expected.Definitions) {
					for i, def := range got.Definitions {
						if tt.expected.Definitions[i].PlainStatements == nil {
							tt.expected.Definitions[i].PlainStatements = def.PlainStatements
						}
					}
				} else {
					a, b := len(tt.expected.Definitions), len(got.Definitions)
					assert.Failf(t, "expected %d definitions, got %d", "", a, b)
				}
			}

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
