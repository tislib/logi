package logi

import (
	"encoding/json"
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
		"json element": {
			input: `
				entity User {
					data { 
						"key1": "value1", 
						"nestedKey": [
							{ 
								"nestedKey": "nestedValue",
								"nestedKey2": null
							}
						] 
					}
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
											Identifier: "data",
										},
									},
									{
										Kind: plain.DefinitionStatementElementKindValue,
										Value: &plain.DefinitionStatementElementValue{
											Value: common.Value{
												Kind: common.ValueKindMap,
												Map: map[string]common.Value{
													"key1": common.StringValue("value1"),
													"nestedKey": {
														Kind: common.ValueKindArray,
														Array: []common.Value{
															{
																Kind: common.ValueKindMap,
																Map: map[string]common.Value{
																	"nestedKey":  common.StringValue("nestedValue"),
																	"nestedKey2": common.NullValue(),
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
						if (age < 18) {
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
										Kind: plain.DefinitionStatementElementKindStruct,

										Struct: &plain.DefinitionStatementElementStruct{
											Statements: []plain.DefinitionStatement{
												{
													Elements: []plain.DefinitionStatementElement{
														{
															Kind: plain.DefinitionStatementElementKindIdentifier,
															Identifier: &plain.DefinitionStatementElementIdentifier{
																Identifier: "if",
															},
														},
														{
															Kind: plain.DefinitionStatementElementKindParameterList,
															ParameterList: &plain.DefinitionStatementElementParameterList{
																Parameters: []common.Expression{
																	common.BinaryExpr("<", common.Var("age"), common.Lit(common.IntegerValue(18))),
																},
															},
														},
														{
															Kind: plain.DefinitionStatementElementKindStruct,
															Struct: &plain.DefinitionStatementElementStruct{
																Statements: []plain.DefinitionStatement{
																	{
																		Elements: []plain.DefinitionStatementElement{
																			{
																				Kind: plain.DefinitionStatementElementKindIdentifier,
																				Identifier: &plain.DefinitionStatementElementIdentifier{
																					Identifier: "return",
																				},
																			},
																			{
																				Kind:  plain.DefinitionStatementElementKindValue,
																				Value: &plain.DefinitionStatementElementValue{Value: common.IntegerValue(0)},
																			},
																		},
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
																Identifier: "return",
															},
														},
														{
															Kind:  plain.DefinitionStatementElementKindValue,
															Value: &plain.DefinitionStatementElementValue{Value: common.IntegerValue(1)},
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
			got, err := ParsePlainContent(tt.input, false)

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
