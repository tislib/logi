package vm

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVmDynamic(t *testing.T) {
	tests := map[string]struct {
		macro                     string
		input                     string
		calls                     []func(vm VirtualMachine) error
		expectedDefinitions       []Definition
		expectedDefinitionsAsJson string
		options                   []Option
	}{
		"simple macro": {
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
			expectedDefinitions: []Definition{
				{
					Macro: "test",
					Name:  "test1",
					Data: map[string]map[string]interface{}{
						"hello": {
							"hello": "hello",
						},
						"world": {
							"world": int64(42),
						},
					},
				},
			},
		},
		"complex macro execution": {
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

					for _, indicator := range definition.Data["indicators"]["indicators"].([]interface{}) {
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
		"complex macro execution with call chain": {
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

					for _, indicator := range definition.Data["indicators"]["indicators"].([]interface{}) {
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
		"creditRule": {
			macro: `
				macro creditRule {
					kind Syntax
					
					syntax {
					  	creditScore <min int> <max int>
						income <min int> <max int>
						age <min int> <max int>
					}
				}`,
			input: `
				creditRule Rule1 {
					creditScore 500 600
					income 20000 30000
					age 18 65
				}
				
				creditRule Rule2 {
					creditScore 600 700
					income 30000 40000
					age 18 65
				}`,
			expectedDefinitions: []Definition{
				{
					Macro: "creditRule",
					Name:  "Rule1",
					Data: map[string]map[string]interface{}{
						"creditScore": {
							"min": int64(500),
							"max": int64(600),
						},
						"income": {
							"min": int64(20000),
							"max": int64(30000),
						},
						"age": {
							"min": int64(18),
							"max": int64(65),
						},
					},
				},
				{
					Macro: "creditRule",
					Name:  "Rule2",
					Data: map[string]map[string]interface{}{
						"creditScore": {
							"min": int64(600),
							"max": int64(700),
						},
						"income": {
							"min": int64(30000),
							"max": int64(40000),
						},
						"age": {
							"min": int64(18),
							"max": int64(65),
						},
					},
				},
			},
		},
		"chatbot": {
			macro: `
				macro chatbot {
					kind Syntax
					
					syntax {
						intent <name Name> {
							pattern <pattern string>
							response <response string>
						}
					}
				}`,
			input: `
				chatbot MyChatbot {
					intent Greeting {
						pattern "Hello"
						response "Hi there!"
					}
					
					intent Farewell {
						pattern "Goodbye"
						response "See you later!"
					}
				}`,
			expectedDefinitions: []Definition{
				{
					Macro: "chatbot",
					Name:  "MyChatbot",
					Data: map[string]map[string]interface{}{
						"intentGreeting": {
							"name": "Greeting",
							"intentGreeting": map[string]interface{}{
								"pattern": map[string]interface{}{
									"pattern": "Hello",
								},
								"response": map[string]interface{}{
									"response": "Hi there!",
								},
							},
						},
						"intentFarewell": {
							"name": "Farewell",
							"intentFarewell": map[string]interface{}{
								"pattern": map[string]interface{}{
									"pattern": "Goodbye",
								},
								"response": map[string]interface{}{
									"response": "See you later!",
								},
							},
						},
					},
				},
			},
		},
		"person1": {
			macro: `
				macro person {
					kind Syntax
					
					syntax {
						name <name string> as known as <knownAs string>
						age <age int> years old
					}
				}`,
			input: `
				person John {
					name "John" as known as "Johnny"
					age 30 years old
				}`,
			expectedDefinitionsAsJson: `
				[
					{
						"macro": "person",
						"name": "John",
						"data": {
							"name": {
								"name": "John",
								"knownAs": "Johnny"
							},
							"age": {
								"age": 30
							}
						}
					}
				]
`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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

			if tt.expectedDefinitionsAsJson != "" {
				var defs []Definition
				err = json.Unmarshal([]byte(tt.expectedDefinitionsAsJson), &defs)

				if err != nil {
					t.Errorf("error: %v", err)
					return
				}

				expected, _ := json.MarshalIndent(defs, "", "  ")
				actual, _ := json.MarshalIndent(definitions, "", "  ")

				assert.Equal(t, string(expected), string(actual))
			}

			if tt.expectedDefinitions != nil {
				assert.Equal(t, tt.expectedDefinitions, definitions)
			}

			for _, call := range tt.calls {
				if err := call(g); err != nil {
					t.Errorf("error: %v", err)
				}
			}
		})
	}
}
