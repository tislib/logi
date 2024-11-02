package logi

import (
	"github.com/stretchr/testify/assert"
	"logi/pkg/ast"
	"strings"
	"testing"
)

func TestSyntaxLogi(t *testing.T) {
	tests := map[string]struct {
		input         string
		expected      *ast.LogiPlainAst
		expectedError string
	}{
		"simple syntax logi": {
			input: `
				entity User {
					id: int [primary, autoincrement]
					name: string [required, default: "John Doe"]
				}
			`,
			expected: &ast.LogiPlainAst{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseLogiPlainContent(tt.input)

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
