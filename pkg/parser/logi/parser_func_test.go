package logi

import (
	"github.com/stretchr/testify/assert"
	"logi/pkg/ast/common"
	logiAst "logi/pkg/ast/logi"
	"strings"
	"testing"
)

func TestParserFunc(t *testing.T) {
	tests := map[string]struct {
		macroInput    string
		input         string
		expected      *logiAst.Ast
		expectedError string
	}{
		"parse main func": {
			macroInput: ``,
			input: `
				func main() {
					println("Hello, World!")
				}
			`,
			expected: &logiAst.Ast{
				Functions: []logiAst.Function{
					{
						Name: "main",
						CodeBlock: common.CodeBlock{
							Statements: []common.Statement{
								{
									Kind: common.FuncCallStatementKind,

									FuncCall: &common.FunctionCallStatement{
										Call: &common.FunctionCall{
											Name: "println",
											Arguments: []*common.Expression{
												{
													Kind: common.LiteralKind,
													Literal: &common.Literal{
														Value: common.StringValue("Hello, World!"),
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
			got, err := ParseFullWithMacro(tt.input, tt.macroInput)

			if got != nil && tt.expected != nil {
				if len(got.Definitions) == len(tt.expected.Definitions) {
					for i, def := range got.Definitions {
						tt.expected.Definitions[i].PlainStatements = def.PlainStatements
					}
				} else {
					assert.Fail(t, "expected %d definitions, got %d", len(tt.expected.Definitions), len(got.Definitions))
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
