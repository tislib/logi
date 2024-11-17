package macro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyntaxMacroErrors(t *testing.T) {
	tests := map[string]struct {
		input         string
		expectedError string
	}{
		"unexpected start": {
			input: `
				simple {
					kind Syntax
				}
			`,
			expectedError: "syntax error at or near \"simple\" at line 2 column 5: unexpected identifier \"simple\", expecting macro keyword",
		},
		"macro incorrect name": {
			input: `
				macro simple Simple2 {
					kind Syntax
				}
			`,
			expectedError: "syntax error at or near \"Simple2\" at line 2 column 18: unexpected identifier \"Simple2\", expecting {",
		},
		"macro incorrect name [2]": {
			input: `
				macro "simple"" {
					kind Syntax
				}
			`,
			expectedError: "syntax error at or near \"simple\" at line 2 column 11: unexpected string \"simple\", expecting identifier",
		},
		"macro missing syntax": {
			input: `
				macro simple {
					syntax {
					}
				}
			`,
			expectedError: "syntax error at or near \"syntax\" at line 3 column 6: unexpected syntax keyword \"syntax\", expecting identifier",
		},
		"macro incorrect syntax": {
			input: `
				macro simple {
					kind Syntax1
					syntax {
					}
				}
			`,
			expectedError: "syntax error at or near \"Syntax1\" at line 3 column 11: unexpected kind value: \"Syntax1\", expecting \"Syntax\"",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ParseMacroContent(tt.input)

			assert.Equal(t, tt.expectedError, err.Error())
		})
	}
}
