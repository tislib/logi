package logi

import (
	"github.com/stretchr/testify/assert"
	"github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
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
					id int <[primary, autoincrement]>
					name string <[required, default "John Doe"]>
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
									{
										Kind: plain.DefinitionStatementElementKindIdentifier,
										Identifier: &plain.DefinitionStatementElementIdentifier{
											Identifier: "int",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindCodeBlock,

										CodeBlock: &plain.DefinitionStatementElementCodeBlock{
											CodeBlock: common.CodeBlock{
												Statements: []common.Statement{
													{
														Kind: common.IfStatementKind,
														IfStmt: &common.IfStatement{
															Condition: &common.Expression{
																Kind: common.BinaryExprKind,
																BinaryExpr: &common.BinaryExpression{
																	Left: &common.Expression{
																		Kind: common.VariableKind,
																		Variable: &common.Variable{
																			Name: "age",
																		},
																	},
																	Operator: "<",
																	Right: &common.Expression{
																		Kind: common.LiteralKind,

																		Literal: &common.Literal{
																			Value: common.IntegerValue(18),
																		},
																	},
																},
															},
															ThenBlock: &common.CodeBlock{
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

													{
														Kind: common.ReturnStatementKind,

														ReturnStmt: &common.ReturnStatement{
															Result: &common.Expression{
																Kind: common.LiteralKind,

																Literal: &common.Literal{
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
			got, err := ParsePlainContent(tt.input)

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
