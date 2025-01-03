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
					id int [primary, autoincrement]
					name string [required, default "John Doe"]
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
		"simple symbol parse": {
			macroInput: `
				macro simple {
					kind Syntax
					
					syntax {
						hello: (...)
					}
				}
`,
			input: `
				simple User {
					hello: (a: b)
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "simple",
						Name:      "User",
						Statements: []logiAst.Statement{
							{
								Command: "hello",
								Parameters: []logiAst.Parameter{
									{
										Name:  "a",
										Value: common.StringValue("b"),
										Expression: &common.Expression{
											Kind: common.VariableKind,
											Variable: &common.Variable{
												Name: "b",
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
					id int [primary, autoincrement]
					name string [required, default "John Doe"]
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
						Statements: []logiAst.Statement{
							{
								Command: "Hello",
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
							},
							{
								Command: "World",
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
							},
						},
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
						Statements: []logiAst.Statement{
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.StringValue("param1"),
									},
									{
										Name:  "value1",
										Value: common.StringValue("11"),
									},
									{
										Name:  "value2",
										Value: common.StringValue("22"),
									},
								},
							},
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.StringValue("param2"),
									},
									{
										Name:  "value3",
										Value: common.IntegerValue(1),
									},
									{
										Name:  "value4",
										Value: common.IntegerValue(2),
									},
								},
							},
						},
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
						Auth { auth }
					}

					scopes { 
						code {
							return <value int>
						}
						auth {
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
						Statements: []logiAst.Statement{
							{
								Command: "Main",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "code",
											Command: "return",
											Parameters: []logiAst.Parameter{
												{
													Name:  "value",
													Value: common.IntegerValue(123),
												},
											},
										},
									},
								},
							},
							{
								Command: "Auth",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "auth",
											Command: "Username",
											Parameters: []logiAst.Parameter{
												{
													Name:  "username",
													Value: common.StringValue("user1"),
												},
											},
										},
										{
											Scope:   "auth",
											Command: "Password",
											Parameters: []logiAst.Parameter{
												{
													Name:  "password",
													Value: common.StringValue("password1"),
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
						Statements: []logiAst.Statement{
							{
								Command: "Hello",
								Parameters: []logiAst.Parameter{
									{
										Name:  "value1",
										Value: common.StringValue("from the other side"),
									},
									{
										Name:  "value2",
										Value: common.StringValue("hello"),
									},
									{
										Name:  "value3",
										Value: common.StringValue("world"),
									},
								},
							},
						},
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
					int id [primary, autoincrement]
					string name [required, default "John Doe"]
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
										Name:  "propertyType",
										Value: common.StringValue("int"),
									},
									{
										Name:  "propertyName",
										Value: common.StringValue("id"),
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
										Name:  "propertyType",
										Value: common.StringValue("string"),
									},
									{
										Name:  "propertyName",
										Value: common.StringValue("name"),
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
					createUser ((name string, age int)) User
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "interface",
						Name:      "UserService",
						Statements: []logiAst.Statement{
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "methodName",
										Value: common.StringValue("createUser"),
									},
									{
										Name:  "returnType",
										Value: common.StringValue("User"),
									},
								},
								Arguments: []logiAst.Argument{
									{
										Name: "name",
										Type: common.TypeDefinition{
											Name: "string",
										},
									},
									{
										Name: "age",
										Type: common.TypeDefinition{
											Name: "int",
										},
									},
								},
							},
						},
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

					scopes {
						code {
							return <value int>
						}
					}
				}
`,
			input: `
				implementation UserServiceImpl {
					createUser ((name string, age int)) User {
						return 0
					}
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "implementation",
						Name:      "UserServiceImpl",
						Statements: []logiAst.Statement{
							{
								Parameters: []logiAst.Parameter{
									{
										Name:  "methodName",
										Value: common.StringValue("createUser"),
									},
									{
										Name:  "returnType",
										Value: common.StringValue("User"),
									},
								},
								Arguments: []logiAst.Argument{
									{
										Name: "name",
										Type: common.TypeDefinition{
											Name: "string",
										},
									},
									{
										Name: "age",
										Type: common.TypeDefinition{
											Name: "int",
										},
									},
								},
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "code",
											Command: "return",
											Parameters: []logiAst.Parameter{
												{
													Name:  "value",
													Value: common.IntegerValue(0),
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
						Statements: []logiAst.Statement{
							{
								Command: "Param1",
								Parameters: []logiAst.Parameter{
									{
										Name:  "value11",
										Value: common.StringValue("value1"),
									},
									{
										Name:  "value21",
										Value: common.StringValue("value2"),
									},
								},
							},
						},
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
						Statements: []logiAst.Statement{
							{
								Command: "Param1",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "value11",
													Value: common.StringValue("value1"),
												},
												{
													Name:  "value21",
													Value: common.StringValue("value2"),
												},
											},
										},
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "value11",
													Value: common.StringValue("value3"),
												},
												{
													Name:  "value21",
													Value: common.StringValue("value4"),
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
						Statements: []logiAst.Statement{
							{
								Command: "Param1",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "value11",
													Value: common.StringValue("value1"),
												},
												{
													Name:  "value21",
													Value: common.StringValue("value2"),
												},
											},
										},
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "value11",
													Value: common.StringValue("value3"),
												},
												{
													Name:  "value21",
													Value: common.StringValue("value4"),
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
						Strategy { strategy }
					}

					scopes { 
						strategy {
							if (<condition bool>) { strategy }
							Buy(<symbol string>, <quantity int>)	
							Sell(<symbol string>, <quantity int>)	
						}
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
							Buy(quantity: 100, symbol: "SPY")
						}
					}
				}
				`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "backtest",
						Name:      "VariableHoldUntil4",
						Statements: []logiAst.Statement{
							{
								Command: "InitialCapital",
								Parameters: []logiAst.Parameter{
									{
										Name:  "initialCapital",
										Value: common.IntegerValue(10000),
									},
								},
							},
							{
								Command: "StartTime",
								Parameters: []logiAst.Parameter{
									{
										Name:  "startTime",
										Value: common.StringValue("2010-01-01"),
									},
								},
							},
							{
								Command: "EndTime",
								Parameters: []logiAst.Parameter{
									{
										Name:  "endTime",
										Value: common.StringValue("2010-12-31"),
									},
								},
							},
							{
								Command: "Indicators",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "indicatorName",
													Value: common.StringValue("sma"),
												},
												{
													Name:  "period",
													Value: common.IntegerValue(20),
													Expression: &common.Expression{
														Kind: common.LiteralKind,
														Literal: &common.Literal{
															Value: common.IntegerValue(20),
														},
													},
												},
												{
													Name:  "alias",
													Value: common.StringValue("sma20"),
												},
											},
										},
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "indicatorName",
													Value: common.StringValue("sma"),
												},
												{
													Name:  "period",
													Value: common.IntegerValue(50),
													Expression: &common.Expression{
														Kind: common.LiteralKind,
														Literal: &common.Literal{
															Value: common.IntegerValue(50),
														},
													},
												},
												{
													Name:  "alias",
													Value: common.StringValue("sma50"),
												},
											},
										},
										{
											Parameters: []logiAst.Parameter{
												{
													Name:  "indicatorName",
													Value: common.StringValue("sma"),
												},
												{
													Name:  "period",
													Value: common.IntegerValue(200),
													Expression: &common.Expression{
														Kind: common.LiteralKind,
														Literal: &common.Literal{
															Value: common.IntegerValue(200),
														},
													},
												},
												{
													Name:  "alias",
													Value: common.StringValue("sma200"),
												},
											},
										},
									},
								},
							},
							{
								Command: "Strategy",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "strategy",
											Command: "if",
											Parameters: []logiAst.Parameter{
												{
													Name: "condition",
													Value: common.MapValue(map[string]common.Value{
														"left":     common.StringValue("sma20"),
														"operator": common.StringValue("<"),
														"right":    common.StringValue("sma50"),
													}),
													Expression: &common.Expression{
														Kind: common.BinaryExprKind,
														BinaryExpr: &common.BinaryExpression{
															Left: &common.Expression{
																Kind: common.VariableKind,
																Variable: &common.Variable{
																	Name: "sma20",
																},
															},
															Operator: "<",
															Right: &common.Expression{
																Kind: common.VariableKind,
																Variable: &common.Variable{
																	Name: "sma50",
																},
															},
														},
													},
												},
											},
											SubStatements: [][]logiAst.Statement{
												{
													{
														Scope:   "strategy",
														Command: "Buy",
														Parameters: []logiAst.Parameter{
															{
																Name:  "symbol",
																Value: common.StringValue("SPY"),
																Expression: &common.Expression{
																	Kind: common.LiteralKind,
																	Literal: &common.Literal{
																		Value: common.StringValue("SPY"),
																	},
																},
															},
															{
																Name:  "quantity",
																Value: common.IntegerValue(100),
																Expression: &common.Expression{
																	Kind: common.LiteralKind,
																	Literal: &common.Literal{
																		Value: common.IntegerValue(100),
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
							Led 	<component Name> <pin int>
							Button 	<component Name> <pin int>
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
						Led 	yellowLed 5
						Led 	redLed 6
						Led 	blueLed 13
						Button 	button1 17
						Button 	button2 17
					}
				
					actions {
						on(yellowLed)
						on(redLed)
				
						on_click(button1) {
							if (status(button2) == 'on') {
								on(blueLed)
								on(yellowLed)
								on(redLed)
							} else {
								off(blueLed)
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
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "components",
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
											Scope:   "components",
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
											Scope:   "components",
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
											Scope:   "components",
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
											Scope:   "components",
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
							},
							{
								Command: "actions",
								SubStatements: [][]logiAst.Statement{
									{
										{
											Scope:   "command",
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
											Scope:   "command",
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
											Scope:   "handler",
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
											SubStatements: [][]logiAst.Statement{
												{
													{
														Scope:   "command",
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
														SubStatements: [][]logiAst.Statement{
															{
																{
																	Scope:   "command",
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
																	Scope:   "command",
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
																	Scope:   "command",
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
															{
																{
																	Scope:   "command",
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
															},
														},
													},
												},
											},
										},
										{
											Scope:   "handler",
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
											SubStatements: [][]logiAst.Statement{
												{
													{
														Scope:   "command",
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
														Scope:   "command",
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
														Scope:   "command",
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
