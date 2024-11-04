package logi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"log"
	"strings"
	"testing"
)

func TestParserFull(t *testing.T) {
	tests := map[string]struct {
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
						Properties: []logiAst.Property{
							{
								Name: "id",
								Type: common.TypeDefinition{
									Name: "int",
								},
								Attributes: []logiAst.Attribute{
									{
										Name: "primary",
									},
									{
										Name: "autoincrement",
									},
								},
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.PointerValue(common.StringValue("id")),
									},
									{
										Name:  "propertyType",
										Value: common.PointerValue(common.StringValue("int")),
									},
								},
							},
							{
								Name: "name",
								Type: common.TypeDefinition{
									Name: "string",
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
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyName",
										Value: common.PointerValue(common.StringValue("name")),
									},
									{
										Name:  "propertyType",
										Value: common.PointerValue(common.StringValue("string")),
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
						<propertyType Type> <propertyName name> [primary bool, autoincrement bool, required bool, default string]
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
						Properties: []logiAst.Property{
							{
								Name: "id",
								Type: common.TypeDefinition{
									Name: "int",
								},
								Attributes: []logiAst.Attribute{
									{
										Name: "primary",
									},
									{
										Name: "autoincrement",
									},
								},
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyType",
										Value: common.PointerValue(common.StringValue("int")),
									},
									{
										Name:  "propertyName",
										Value: common.PointerValue(common.StringValue("id")),
									},
								},
							},
							{
								Name: "name",
								Type: common.TypeDefinition{
									Name: "string",
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
								Parameters: []logiAst.Parameter{
									{
										Name:  "propertyType",
										Value: common.PointerValue(common.StringValue("string")),
									},
									{
										Name:  "propertyName",
										Value: common.PointerValue(common.StringValue("name")),
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
					createUser (name string, age int) User
				}
			`,
			expected: &logiAst.Ast{
				Definitions: []logiAst.Definition{
					{
						MacroName: "interface",
						Name:      "UserService",
						MethodSignature: []logiAst.MethodSignature{
							{
								Name: "createUser",
								Type: common.TypeDefinition{
									Name: "User",
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
								Parameters: []logiAst.Parameter{
									{
										Name:  "methodName",
										Value: common.PointerValue(common.StringValue("createUser")),
									},
									{
										Name:  "returnType",
										Value: common.PointerValue(common.StringValue("User")),
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
						<methodName Name> (...[<args Type<string>>]) <returnType Type> { }
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
						Methods: []logiAst.Method{
							{
								Name: "createUser",
								Type: common.TypeDefinition{
									Name: "User",
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
								Parameters: []logiAst.Parameter{
									{
										Name:  "methodName",
										Value: common.PointerValue(common.StringValue("createUser")),
									},
									{
										Name:  "returnType",
										Value: common.PointerValue(common.StringValue("User")),
									},
								},
								CodeBlock: common.CodeBlock{
									Statements: []common.Statement{
										{
											Kind: common.ReturnStatementKind,

											ReturnStmt: &common.ReturnStatement{
												Result: &common.Expression{
													Kind: common.LiteralKind,

													Literal: &common.Literal{
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
		},
		"test arrays with definition": {
			macroInput: `
				macro backtest {
					kind Syntax
				
					types {
						ParamValue <value string>
					}
				
					syntax {
						Param <name string> <value ParamValue>
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
		"test backtest macro": {
			macroInput: `
				macro backtest {
					kind Syntax
				
					definition {
						Indicator <indicatorName string> (<period int>) as <alias string>
					}
				
					syntax {
						InitialCapital <initialCapital int>
						StartTime <startTime string>
						EndTime <endTime string>
						Indicators <indicators array<Indicator>>
						Strategy { }
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
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseFullWithMacro(tt.input, tt.macroInput)

			data, _ := json.MarshalIndent(got, "", "  ")

			log.Print(string(data))

			if got != nil && tt.expected != nil {
				if len(got.Definitions) == len(tt.expected.Definitions) {
					for i, def := range got.Definitions {
						tt.expected.Definitions[i].PlainStatements = def.PlainStatements
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

				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func init() {
	yyDebug = 5
}
