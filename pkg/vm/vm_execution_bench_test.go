package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkVmDynamic(t *testing.B) {
	tests := []struct {
		name               string
		macro              string
		input              string
		calls              []func(vm VirtualMachine) error
		expectedDefinition *Definition
		options            []Option
	}{
		{
			name: "simple macro",
			macro: `
				macro test {
					kind Syntax
					syntax {
						Hello <hello string>
						World <world number>
					}
				}`,
			input: `
				test test1 {
					Hello "hello"
					World 42
				}`,
			expectedDefinition: &Definition{
				Macro: "test",
				Name:  "test1",
				Data: map[string]interface{}{
					"hello": "hello",
					"world": int64(42),
				},
			},
		},
		{
			name: "complex macro execution",
			macro: `
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
							return Buy("SPY", 100)
						}

						return Buy("SPY", 200)
					}
				}
				`,
			options: []Option{
				WithLocals(map[string]interface{}{
					"Buy": func(symbol string, quantity int64) int64 {
						return quantity
					},
				}),
			},
			calls: []func(vm VirtualMachine) error{
				func(vm VirtualMachine) error {
					definition, err := vm.GetDefinitionByName("VariableHoldUntil4")

					if err != nil {
						return err
					}

					newLocals := map[string]interface{}{}

					for _, indicator := range definition.Data["indicators"].([]interface{}) {
						indicatorMap := indicator.(map[string]interface{})
						newLocals[indicatorMap["alias"].(string)] = indicatorMap["period"].(int64)
					}

					vm.SetLocals(newLocals)

					f, err := vm.LocateCodeBlock(*definition, "strategy")

					if err != nil {
						return err
					}

					result, err := f()

					if err != nil {
						return err
					}

					assert.Equal(t, result, int64(100))

					return nil
				},
			},
		},
		{
			name: "complex macro execution with call chain",
			macro: `
				macro backtest {
					kind Syntax
				
					types {
						Indicator <indicatorName Name> ((<period int>)|(<period int>, <period2 int>)) as <alias Name>
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
					Indicators       [sma(20) as sma20, sma(50) as sma50, sma(200, 300) as sma200]
				
					Strategy {
						return Sub(Add(Add(Add(Add(1, 2), 3), 4), 5), 6)
					}
				}
				`,
			options: []Option{
				WithLocals(map[string]interface{}{
					"Add": func(a, b int64) int64 {
						return a + b
					},
					"Sub": func(a, b int64) int64 {
						return a - b
					},
				}),
			},
			calls: []func(vm VirtualMachine) error{
				func(vm VirtualMachine) error {
					definition, err := vm.GetDefinitionByName("VariableHoldUntil4")

					if err != nil {
						return err
					}

					newLocals := map[string]interface{}{}

					for _, indicator := range definition.Data["indicators"].([]interface{}) {
						indicatorMap := indicator.(map[string]interface{})
						newLocals[indicatorMap["alias"].(string)] = indicatorMap["period"].(int64)
					}

					vm.SetLocals(newLocals)

					f, err := vm.LocateCodeBlock(*definition, "strategy")

					if err != nil {
						return err
					}

					result, err := f()

					if err != nil {
						return err
					}

					assert.Equal(t, result, int64(9))

					return nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			t.RunParallel(func(pb *testing.PB) {
				var g, err = New(tt.options...)

				if err != nil {
					t.Errorf("error: %v", err)
					return
				}

				err = g.LoadMacroContent(tt.macro)

				if err != nil {
					t.Errorf("error: %v", err)
					return
				}

				definitions, err := g.LoadLogiContent(tt.input)

				if err != nil {
					t.Errorf("error: %v", err)
					return
				}

				if tt.expectedDefinition != nil {
					assert.Len(t, definitions, 1)
					if len(definitions) == 1 {
						assert.Equal(t, *tt.expectedDefinition, definitions[0])
					}
				}

				for pb.Next() {
					for _, call := range tt.calls {
						if err := call(g); err != nil {
							t.Errorf("error: %v", err)
						}
					}
				}
			})
		})
	}
}
