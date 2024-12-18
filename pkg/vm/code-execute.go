package vm

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
)

func (v *vm) Execute(def *logiAst.Definition, implementer Implementer) error {
	for _, statement := range def.Statements {
		if err := implementer.Call(v, statement); err != nil {
			return fmt.Errorf("failed to execute statement: %w at %v", err, statement)
		}
	}
	return nil
}

func (v *vm) Evaluate(expression common.Expression, vars map[string]common.Value, fns map[string]func(args ...common.Value) (common.Value, error)) (common.Value, error) {
	switch expression.Kind {
	case common.LiteralKind:
		return expression.Literal.Value, nil
	case common.VariableKind:
		value, ok := vars[expression.Variable.Name]

		if !ok {
			return common.Value{}, fmt.Errorf("variable %s not found", expression.Variable.Name)
		}

		return value, nil
	case common.BinaryExprKind:
		return v.evaluateBinaryExpression(expression.BinaryExpr, vars, fns)
	case common.FuncCallKind:
		fn, ok := fns[expression.FuncCall.Name]

		if !ok {
			return common.Value{}, fmt.Errorf("function %s not found", expression.FuncCall.Name)
		}

		args := make([]common.Value, 0)
		for _, arg := range expression.FuncCall.Arguments {
			value, err := v.Evaluate(*arg, vars, fns)
			if err != nil {
				return common.Value{}, fmt.Errorf("failed to evaluate argument: %w", err)
			}
			args = append(args, value)
		}

		return fn(args...)
	default:
		return common.Value{}, fmt.Errorf("unknown expression kind: %s", expression.Kind)
	}
}

func (v *vm) evaluateBinaryExpression(expr *common.BinaryExpression, vars map[string]common.Value, fns map[string]func(args ...common.Value) (common.Value, error)) (common.Value, error) {
	leftValue, err := v.Evaluate(*expr.Left, vars, fns)
	if err != nil {
		return common.Value{}, fmt.Errorf("failed to evaluate left expression: %w", err)
	}

	rightValue, err := v.Evaluate(*expr.Right, vars, fns)
	if err != nil {
		return common.Value{}, fmt.Errorf("failed to evaluate right expression: %w", err)
	}

	return v.evaluateBinaryExpressionOnValues(expr.Operator, leftValue, rightValue)
}

func (v *vm) evaluateBinaryExpressionOnValues(operator string, a common.Value, b common.Value) (common.Value, error) {
	switch a.Kind {
	case common.ValueKindString:
		if b.Kind != common.ValueKindString {
			return common.Value{}, fmt.Errorf("left expression is a string, but right expression is not")
		}
		return v.evaluateBinaryExpressionOnString(operator, a.AsString(), b.AsString())
	case common.ValueKindBoolean:
		if b.Kind != common.ValueKindBoolean {
			return common.Value{}, fmt.Errorf("left expression is a boolean, but right expression is not")
		}

		return v.evaluateBinaryExpressionOnBoolean(operator, a.AsBoolean(), b.AsBoolean())
	case common.ValueKindFloat:
		if b.Kind == common.ValueKindFloat {
			return v.evaluateBinaryExpressionOnFloat(operator, a.AsFloat(), b.AsFloat())
		} else if b.Kind == common.ValueKindInteger {
			return v.evaluateBinaryExpressionOnFloat(operator, a.AsFloat(), float64(b.AsInteger()))
		} else {
			return common.Value{}, fmt.Errorf("left expression is a float, but right expression is not")
		}
	case common.ValueKindInteger:
		if b.Kind == common.ValueKindFloat {
			return v.evaluateBinaryExpressionOnFloat(operator, float64(a.AsInteger()), b.AsFloat())
		} else if b.Kind == common.ValueKindInteger {
			return v.evaluateBinaryExpressionOnInteger(operator, a.AsInteger(), b.AsInteger())
		} else {
			return common.Value{}, fmt.Errorf("left expression is an integer, but right expression is not")
		}
	case common.ValueKindArray:
		if b.Kind != common.ValueKindArray {
			return common.Value{}, fmt.Errorf("left expression is an array, but right expression is not")
		}

		return v.evaluateBinaryExpressionOnArray(operator, a.AsArray(), b.AsArray())
	case common.ValueKindMap:
		if b.Kind != common.ValueKindMap {
			return common.Value{}, fmt.Errorf("left expression is a map, but right expression is not")
		}

		return v.evaluateBinaryExpressionOnMap(operator, a.AsMap(), b.AsMap())
	default:
		return common.Value{}, fmt.Errorf("unknown value kind: %s", a.Kind)
	}
}

func (v *vm) evaluateBinaryExpressionOnString(operator string, left string, right string) (common.Value, error) {
	switch operator {
	case "==":
		return common.BooleanValue(left == right), nil
	case "!=":
		return common.BooleanValue(left != right), nil
	case "+":
		return common.StringValue(left + right), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s", operator)
	}
}

func (v *vm) evaluateBinaryExpressionOnBoolean(operator string, left bool, right bool) (common.Value, error) {
	switch operator {
	case "==":
		return common.BooleanValue(left == right), nil
	case "!=":
		return common.BooleanValue(left != right), nil
	case "&&":
		return common.BooleanValue(left && right), nil
	case "||":
		return common.BooleanValue(left || right), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s, allowed: ==, !=, &&, || for booleans", operator)
	}
}

func (v *vm) evaluateBinaryExpressionOnFloat(operator string, left float64, right float64) (common.Value, error) {
	switch operator {
	case "+":
		return common.FloatValue(left + right), nil
	case "-":
		return common.FloatValue(left - right), nil
	case "*":
		return common.FloatValue(left * right), nil
	case "/":
		return common.FloatValue(left / right), nil
	case "==":
		return common.BooleanValue(left == right), nil
	case "!=":
		return common.BooleanValue(left != right), nil
	case ">":
		return common.BooleanValue(left > right), nil
	case "<":
		return common.BooleanValue(left < right), nil
	case ">=":
		return common.BooleanValue(left >= right), nil
	case "<=":
		return common.BooleanValue(left <= right), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s", operator)
	}
}

func (v *vm) evaluateBinaryExpressionOnInteger(operator string, left int64, right int64) (common.Value, error) {
	switch operator {
	case "+":
		return common.IntegerValue(left + right), nil
	case "-":
		return common.IntegerValue(left - right), nil
	case "*":
		return common.IntegerValue(left * right), nil
	case "/":
		return common.IntegerValue(left / right), nil
	case "==":
		return common.BooleanValue(left == right), nil
	case "!=":
		return common.BooleanValue(left != right), nil
	case ">":
		return common.BooleanValue(left > right), nil
	case "<":
		return common.BooleanValue(left < right), nil
	case ">=":
		return common.BooleanValue(left >= right), nil
	case "<=":
		return common.BooleanValue(left <= right), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s", operator)
	}
}

func (v *vm) evaluateBinaryExpressionOnArray(operator string, a []common.Value, b []common.Value) (common.Value, error) {
	switch operator {
	case "==":
		if len(a) != len(b) {
			return common.BooleanValue(false), nil
		}

		for i, value := range a {
			ev, err := v.evaluateBinaryExpressionOnValues("==", value, b[i])

			if err != nil {
				return common.BooleanValue(false), err
			}
			if ev.AsBoolean() {
				return common.BooleanValue(false), nil
			}
		}

		return common.BooleanValue(true), nil
	case "!=":
		ev, err := v.evaluateBinaryExpressionOnArray("==", a, b)

		if err != nil {
			return common.BooleanValue(false), err
		}

		return common.BooleanValue(!ev.AsBoolean()), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s", operator)
	}
}

func (v *vm) evaluateBinaryExpressionOnMap(operator string, a map[string]common.Value, b map[string]common.Value) (common.Value, error) {
	switch operator {
	case "==":
		if len(a) != len(b) {
			return common.BooleanValue(false), nil
		}

		for i, value := range a {
			bValue, ok := b[i]

			if !ok {
				return common.BooleanValue(false), fmt.Errorf("key %s not found in right map", i)
			}
			ev, err := v.evaluateBinaryExpressionOnValues("==", value, bValue)

			if err != nil {
				return common.BooleanValue(false), err
			}
			if ev.AsBoolean() {
				return common.BooleanValue(false), nil
			}
		}

		return common.BooleanValue(true), nil
	case "!=":
		ev, err := v.evaluateBinaryExpressionOnMap("==", a, b)

		if err != nil {
			return common.BooleanValue(false), err
		}

		return common.BooleanValue(!ev.AsBoolean()), nil
	default:
		return common.NullValue(), fmt.Errorf("unknown operator: %s", operator)
	}
}
