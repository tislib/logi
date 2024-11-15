package golang

import (
	"github.com/tislib/logi/pkg/parser/macro"
	"go/format"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name          string
		macroInput    string
		expectedCode  string
		expectedError string
		skip          bool
	}{
		{
			name: "simple macro",
			macroInput: `
				macro test {
					kind Syntax
					syntax {
						Hello <hello string>
						World <world number>
					}
				}`,
			expectedCode: `package model
				type Test struct {	
					Hello string
					World int
				}`,
		},
		{
			name: "macro with nested struct",
			skip: true,
			macroInput: `
				macro user {
					kind Syntax
					syntax {
						fullName <fullName string>
						Auth {
							username <username string>
							password <password string>
							SessionConfig {
								sessionTimeout <sessionTimeout int>
								sessionToken <sessionToken string>
							}	
						}
					}
				}`,
			expectedCode: `package model
				type User struct {
					FullName string
					Auth struct {
						Username string
						Password string
						SessionConfig struct {
							SessionTimeout int
							SessionToken string
						}
					}
				}`,
		},
	}

	var g = NewGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			var macroAst, err = macro.ParseMacroContent(tt.macroInput)

			if err != nil {
				t.Errorf("failed to parse macro, %s", err)
				return
			}

			code, err := g.Generate(macroAst, "model")
			if err != nil {
				if err.Error() != tt.expectedError {
					t.Errorf("expected error %s, got %s", tt.expectedError, err.Error())
				}
			}

			expectedCodeFormatted, err := format.Source([]byte(tt.expectedCode))

			if err != nil {
				t.Errorf("failed to format expected code, %s", err)
				return
			}

			if code != string(expectedCodeFormatted) {
				t.Errorf("expected code %s, got %s", tt.expectedCode, code)
			}
		})
	}
}
