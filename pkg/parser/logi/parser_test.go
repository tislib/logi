package logi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"logi/pkg/ast/common"
	"logi/pkg/ast/plain"
	"strings"
	"testing"
)

func TestSyntaxLogi(t *testing.T) {
	tests := map[string]struct {
		input         string
		expected      *plain.Ast
		expectedError string
	}{
		"simple syntax logi": {
			input: `
				entity User {
					id int [primary, autoincrement]
					name string [required, default "John Doe"]
				}
			`,
			expected: &plain.Ast{
				Definitions: []plain.Definition{
					{
						MacroName: "entity",
						Name:      "User",
						Statements: []plain.DefinitionStatement{
							{
								Elements: []plain.DefinitionStatementElement{
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "id",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "int",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindAttributeList,
										AttributeList: &plain.DefinitionStatementElementAttributeList{
											Attributes: []plain.DefinitionStatementElementAttribute{
												{
													Name: "primary",
												},
												{
													Name: "autoincrement",
												},
											},
										},
									},
								},
							},
							{
								Elements: []plain.DefinitionStatementElement{
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "name",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "string",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindAttributeList,
										AttributeList: &plain.DefinitionStatementElementAttributeList{
											Attributes: []plain.DefinitionStatementElementAttribute{
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
				},
			},
		},
		"definition with arguments": {
			input: `
				service UserService {
					createUser (name string, age int)
				}
			`,
			expected: &plain.Ast{
				Definitions: []plain.Definition{
					{
						MacroName: "service",
						Name:      "UserService",
						Statements: []plain.DefinitionStatement{
							{
								Elements: []plain.DefinitionStatementElement{
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "createUser",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindArgumentList,
										ArgumentList: &plain.DefinitionStatementElementArgumentList{
											Arguments: []plain.DefinitionStatementElementArgument{
												{
													Name: "name",
													TypeDefinition: &common.TypeDefinition{
														Name: "string",
													},
												},
												{
													Name: "age",
													TypeDefinition: &common.TypeDefinition{
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
				},
			},
		},
		"definition with code block and arguments": {
			input: `
				service UserService {
					createUser (name string, age int) int {
						if age < 18 {
							return 0
						}

						return 1
					}
				}
			`,
			expected: &plain.Ast{
				Definitions: []plain.Definition{
					{
						MacroName: "service",
						Name:      "UserService",
						Statements: []plain.DefinitionStatement{
							{
								Elements: []plain.DefinitionStatementElement{
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "createUser",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindArgumentList,
										ArgumentList: &plain.DefinitionStatementElementArgumentList{
											Arguments: []plain.DefinitionStatementElementArgument{
												{
													Name: "name",
													TypeDefinition: &common.TypeDefinition{
														Name: "string",
													},
												},
												{
													Name: "age",
													TypeDefinition: &common.TypeDefinition{
														Name: "int",
													},
												},
											},
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "int",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindCodeBlock,

										CodeBlock: &plain.DefinitionStatementElementCodeBlock{
											CodeBlock: plain.CodeBlock{
												Statements: []plain.Statement{
													{
														Kind: plain.IfStatementKind,
														IfStmt: &plain.IfStatement{
															Condition: &plain.Expression{
																Kind: plain.BinaryExprKind,
																BinaryExpr: &plain.BinaryExpression{
																	Left: &plain.Expression{
																		Kind: plain.VariableKind,
																		Variable: &plain.Variable{
																			Name: "age",
																		},
																	},
																	Operator: "<",
																	Right: &plain.Expression{
																		Kind: plain.LiteralKind,

																		Literal: &plain.Literal{
																			Value: common.IntegerValue(18),
																		},
																	},
																},
															},
															ThenBlock: &plain.CodeBlock{
																Statements: []plain.Statement{
																	{
																		Kind: plain.ReturnStatementKind,

																		ReturnStmt: &plain.ReturnStatement{
																			Result: &plain.Expression{
																				Kind: plain.LiteralKind,

																				Literal: &plain.Literal{
																					Value: common.IntegerValue(0),
																				},
																			},
																		},
																	},
																},
															},
														},
													},

													{
														Kind: plain.ReturnStatementKind,

														ReturnStmt: &plain.ReturnStatement{
															Result: &plain.Expression{
																Kind: plain.LiteralKind,

																Literal: &plain.Literal{
																	Value: common.IntegerValue(1),
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
			got, err := ParseLogiPlainContent(tt.input)

			gotJ, _ := json.Marshal(got)

			expJ, _ := json.Marshal(tt.expected)

			log.Print(string(gotJ))
			log.Print(string(expJ))

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
