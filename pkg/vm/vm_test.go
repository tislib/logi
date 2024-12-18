package vm

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"testing"
)

type TestLedImplementor struct {
	leds        map[string]int64
	buttons     map[string]int64
	ledState    map[string]bool
	buttonState map[string]bool

	buttonHandlers map[string]func() error
}

func (t *TestLedImplementor) Call(vm VirtualMachine, statement logiAst.Statement) error {
	var funcs = map[string]func(args ...common.Value) (common.Value, error){
		"status": func(args ...common.Value) (common.Value, error) {
			if len(args) != 1 {
				return common.Value{}, fmt.Errorf("expected 1 argument")
			}

			var buttonName = args[0].AsMap()["name"].AsString()

			if _, ok := t.buttons[buttonName]; !ok {
				return common.Value{}, fmt.Errorf("unknown button %s", buttonName)
			}

			if t.buttonState[buttonName] {
				return common.StringValue("on"), nil
			} else {
				return common.StringValue("off"), nil
			}
		},
	}

	var vars = map[string]common.Value{}

	for ledName := range t.leds {
		vars[ledName] = common.MapValue(map[string]common.Value{
			"kind": common.StringValue("led"),
			"name": common.StringValue(ledName),
		})
	}

	for buttonName := range t.buttons {
		vars[buttonName] = common.MapValue(map[string]common.Value{
			"kind": common.StringValue("button"),
			"name": common.StringValue(buttonName),
		})
	}

	switch statement.Scope {
	case "components":
		switch statement.Command {
		case "Led":
			t.leds[statement.GetParameter("component").AsString()] = statement.GetParameter("pin").AsInteger()
			return nil
		case "Button":
			t.buttons[statement.GetParameter("component").AsString()] = statement.GetParameter("pin").AsInteger()
			return nil
		default:
			return fmt.Errorf("unknown command %s", statement.Command)
		}
	case "command":
		switch statement.Command {
		case "on":
			var component = statement.GetParameter("component").AsString()
			if _, ok := t.leds[component]; !ok {
				return fmt.Errorf("unknown led %s", component)
			}
			t.ledState[component] = true
			return nil
		case "off":
			var component = statement.GetParameter("component").AsString()
			if _, ok := t.leds[component]; !ok {
				return fmt.Errorf("unknown led %s", component)
			}
			t.ledState[component] = false
			return nil
		case "if":
			var condition *common.Expression
			for _, parameter := range statement.Parameters {
				if parameter.Name == "condition" {
					condition = parameter.Expression
				}
			}

			if condition == nil {
				return fmt.Errorf("condition is required")
			}

			conditionRes, err := vm.Evaluate(*condition, vars, funcs)

			if err != nil {
				return err
			}

			if conditionRes.Kind != common.ValueKindBoolean {
				return fmt.Errorf("condition must be boolean")
			}

			var ifPass = conditionRes.AsBoolean()
			if ifPass {
				for _, subStatement := range statement.SubStatements[0] {
					err := t.Call(vm, subStatement)
					if err != nil {
						return err
					}
				}
			} else if len(statement.SubStatements) > 1 {
				for _, subStatement := range statement.SubStatements[1] {
					err := t.Call(vm, subStatement)
					if err != nil {
						return err
					}
				}
			}
			return nil
		default:
			return fmt.Errorf("unknown command %s", statement.Command)
		}
	case "handler":
		switch statement.Command {
		case "on_click":
			var component = statement.GetParameter("component").AsString()
			if _, ok := t.buttons[component]; !ok {
				return fmt.Errorf("unknown button %s", component)
			}

			t.buttonHandlers[component] = func() error {
				for _, subStatement := range statement.SubStatements[0] {
					err := t.Call(vm, subStatement)
					if err != nil {
						return err
					}
				}

				return nil
			}

			return nil
		default:
			return fmt.Errorf("unknown command %s", statement.Command)
		}
	case "":
		for _, statement := range statement.SubStatements[0] {
			err := t.Call(vm, statement)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown scope %s", statement.Scope)
	}
}

//go:embed test_data/circuit.lgm
var circuitLgm string

//go:embed test_data/circuit.lg
var circuitLg string

func TestVmDynamic(t *testing.T) {
	tests := map[string]struct {
		macro       string
		input       string
		Implementer Implementer
		Assert      func(t *testing.T, implementer Implementer)
	}{
		"simple macro": {
			macro: circuitLgm,
			input: circuitLg,
			Implementer: &TestLedImplementor{
				leds:           make(map[string]int64),
				buttons:        make(map[string]int64),
				ledState:       make(map[string]bool),
				buttonState:    make(map[string]bool),
				buttonHandlers: make(map[string]func() error),
			},
			Assert: func(t *testing.T, implementer Implementer) {
				ledImplementor := implementer.(*TestLedImplementor)

				assert.Equal(t, ledImplementor.leds, map[string]int64{
					"blueLed":   13,
					"redLed":    6,
					"yellowLed": 5,
				})

				assert.Equal(t, ledImplementor.buttons, map[string]int64{
					"button1": 17,
					"button2": 19,
				})

				assert.Equal(t, ledImplementor.ledState, map[string]bool{
					"redLed":    true,
					"yellowLed": true,
				})
			},
		},
		"simple macro, button click": {
			macro: circuitLgm,
			input: circuitLg,
			Implementer: &TestLedImplementor{
				leds:           make(map[string]int64),
				buttons:        make(map[string]int64),
				ledState:       make(map[string]bool),
				buttonState:    make(map[string]bool),
				buttonHandlers: make(map[string]func() error),
			},
			Assert: func(t *testing.T, implementer Implementer) {
				ledImplementor := implementer.(*TestLedImplementor)

				ledImplementor.buttonState["button2"] = true

				var err = ledImplementor.buttonHandlers["button1"]()
				assert.NoErrorf(t, err, fmt.Sprintf("error: %v", err))

				assert.Equal(t, ledImplementor.ledState, map[string]bool{
					"blueLed":   true,
					"redLed":    true,
					"yellowLed": true,
				})

				ledImplementor.buttonState["button2"] = false

				err = ledImplementor.buttonHandlers["button1"]()
				assert.NoErrorf(t, err, fmt.Sprintf("error: %v", err))

				assert.Equal(t, ledImplementor.ledState, map[string]bool{
					"blueLed":   false,
					"redLed":    true,
					"yellowLed": true,
				})
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var g = New()

			err := g.LoadMacroContent(tt.macro)

			if err != nil {
				t.Errorf("error: %v", err)
				return
			}

			definitions, err := g.LoadLogiContent(tt.input)

			if err != nil {
				t.Errorf("error: %v", err)
				return
			}

			for _, definition := range definitions {
				err := g.Execute(&definition, tt.Implementer)

				if err != nil {
					t.Errorf("error: %v", err)
					return
				}
			}

			if tt.Assert != nil {
				tt.Assert(t, tt.Implementer)
			}
		})
	}
}
