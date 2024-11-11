package vm

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	"reflect"
)

func (v *vm) executableFunc(block common.CodeBlock) ExecutableFunc {
	return func(args ...interface{}) (interface{}, error) {
		for _, statement := range block.Statements {
			switch statement.Kind {
			case common.ExprStatementKind:
				_, err := v.evaluateExpression(statement.ExprStmt.Expr)

				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression: %w", err)
				}
			case common.IfStatementKind:
				con, err := v.evaluateBool(statement.IfStmt.Condition)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate condition: %w", err)
				}
				if con {
					return v.executableFunc(*statement.IfStmt.ThenBlock)(args...)
				} else if statement.IfStmt.ElseBlock != nil {
					return v.executableFunc(*statement.IfStmt.ElseBlock)(args...)
				}
			case common.ReturnStatementKind:
				return v.evaluateExpression(statement.ReturnStmt.Result)
			case common.VarDeclKind:
				exprRes, err := v.evaluateExpression(statement.VarDecl.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression: %w", err)
				}
				v.vars[statement.VarDecl.Name] = exprRes
				v.types[statement.VarDecl.Name] = statement.VarDecl.Type
			case common.FuncCallStatementKind:
				_, err := v.evaluateExpression(&common.Expression{
					Kind:     common.FuncCallKind,
					FuncCall: statement.FuncCall.Call,
				})

				if err != nil {
					return nil, fmt.Errorf("failed to evaluate function call: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown statement kind: %s", statement.Kind)
			}
		}

		return nil, nil
	}
}

func (v *vm) evaluateBool(expression *common.Expression) (bool, error) {
	result, err := v.evaluateExpression(expression)

	if err != nil {
		return false, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	switch result.(type) {
	case bool:
		return result.(bool), nil
	default:
		return false, fmt.Errorf("expression result is not a boolean")
	}
}

func (v *vm) evaluateExpression(expression *common.Expression) (interface{}, error) {
	switch expression.Kind {
	case common.LiteralKind:
		return expression.Literal.Value.AsInterface(), nil
	case common.VariableKind:
		if v.vars[expression.Variable.Name] != nil {
			return v.vars[expression.Variable.Name], nil
		} else if v.locals[expression.Variable.Name] != nil {
			return v.locals[expression.Variable.Name], nil
		} else {
			return nil, fmt.Errorf("variable %s not found", expression.Variable.Name)
		}
	case common.BinaryExprKind:
		return v.evaluateBinaryExpression(expression.BinaryExpr)
	case common.FuncCallKind:
		return v.evaluateFunctionCall(expression.FuncCall)
	default:
		return nil, fmt.Errorf("unknown expression kind: %s", expression.Kind)
	}
}

func (v *vm) evaluateBinaryExpression(binaryExpr *common.BinaryExpression) (interface{}, error) {
	leftE, err := v.evaluateExpression(binaryExpr.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate left expression: %w", err)
	}
	rightE, err := v.evaluateExpression(binaryExpr.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate right expression: %w", err)
	}

	leftType, err := v.getTypeInformation(leftE)
	if err != nil {
		return nil, fmt.Errorf("failed to get type information for left expression: %w", err)
	}
	rightType, err := v.getTypeInformation(rightE)
	if err != nil {
		return nil, fmt.Errorf("failed to get type information for left expression: %w", err)
	}

	// As numeric
	if leftType.numeric != rightType.numeric {
		return nil, fmt.Errorf("left and right expressions are not of the same type; one is numberic and the other is not: %v, %v", leftType.kind.String(), rightType.kind.String())
	}

	if leftType.numeric {
		if leftType.floating || rightType.floating {
			var leftFloat, rightFloat float64
			if leftType.floating {
				leftFloat = leftType.floatValue
			} else {
				leftFloat = float64(leftType.intValue)
			}
			if rightType.floating {
				rightFloat = rightType.floatValue
			} else {
				rightFloat = float64(rightType.intValue)
			}
			return v.evaluateBinaryExpressionFloat(binaryExpr, leftFloat, rightFloat)
		}
		return v.evaluateBinaryExpressionInt(binaryExpr, leftType.intValue, rightType.intValue)
	}
	// As string
	if leftType.string != rightType.string {
		return nil, fmt.Errorf("left and right expressions are not of the same type; one is string and the other is not: %v, %v", leftType.kind.String(), rightType.kind.String())
	}

	if leftType.string {
		return v.evaluateBinaryExpressionString(binaryExpr, leftType.stringValue, rightType.stringValue)
	}

	// As boolean
	if leftType.boolean != rightType.boolean {
		return nil, fmt.Errorf("left and right expressions are not of the same type; one is boolean and the other is not: %v, %v", leftType.kind.String(), rightType.kind.String())
	}

	if leftType.boolean {
		return v.evaluateBinaryExpressionBool(binaryExpr, leftType.boolValue, rightType.boolValue)
	}

	return nil, fmt.Errorf("unknown expression type: %v", leftType.kind.String())
}

func (v *vm) evaluateBinaryExpressionFloat(expr *common.BinaryExpression, left float64, right float64) (interface{}, error) {
	switch expr.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		return left / right, nil
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case ">=":
		return left >= right, nil
	case "<=":
		return left <= right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

func (v *vm) evaluateBinaryExpressionInt(expr *common.BinaryExpression, left int64, right int64) (interface{}, error) {
	switch expr.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		return left / right, nil
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case ">=":
		return left >= right, nil
	case "<=":
		return left <= right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

func (v *vm) evaluateBinaryExpressionString(expr *common.BinaryExpression, left string, right string) (interface{}, error) {
	switch expr.Operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

func (v *vm) evaluateBinaryExpressionBool(expr *common.BinaryExpression, left bool, right bool) (interface{}, error) {
	switch expr.Operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case "&&":
		return left && right, nil
	case "||":
		return left || right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

type typeInformation struct {
	kind     reflect.Kind
	numeric  bool
	boolean  bool
	string   bool
	callable bool
	floating bool

	boolValue     bool
	intValue      int64
	floatValue    float64
	stringValue   string
	callableValue ExecutableFunc
}

func (v *vm) getTypeInformation(value interface{}) (typeInformation, error) {
	refValue := reflect.ValueOf(value)

	var res typeInformation

	res.kind = refValue.Kind()

	switch refValue.Kind() {
	case reflect.Bool:
		res.boolean = true
		res.boolValue = refValue.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res.numeric = true
		res.intValue = refValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		res.numeric = true
		res.intValue = int64(refValue.Uint())
	case reflect.Float32, reflect.Float64:
		res.floating = true
		res.numeric = true
		res.floatValue = refValue.Float()
	case reflect.String:
		res.string = true
		res.stringValue = refValue.String()
	case reflect.Func:
		res.callable = true
		execFunc, ok := value.(ExecutableFunc)

		if !ok {
			return res, fmt.Errorf("value is not a function in correct format: func(args ...interface{}) (interface{}, error)")
		}

		res.callableValue = execFunc
	default:
		return res, fmt.Errorf("unknown type: %s", refValue.Kind())
	}

	return res, nil
}

func (v *vm) evaluateFunctionCall(call *common.FunctionCall) (interface{}, error) {
	// locate function

	var fnObj interface{}
	if v.vars[call.Name] != nil {
		fnObj = v.vars[call.Name]
	} else if v.locals[call.Name] != nil {
		fnObj = v.locals[call.Name]
	} else {
		return nil, fmt.Errorf("function %s not found", call.Name)
	}

	// check is func
	refValue := reflect.ValueOf(fnObj)

	if refValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("function %s is not a function", call.Name)
	}

	fn, ok := fnObj.(ExecutableFunc)
	if !ok {
		refType := refValue.Type()

		// convert to func
		fn = func(args ...interface{}) (interface{}, error) {
			if refType.NumIn() != len(args) {
				return nil, fmt.Errorf("function %s expected %d arguments, got %d", call.Name, refType.NumIn(), len(args))
			}

			argsValues := make([]reflect.Value, len(args))
			for i, arg := range args {
				argExpectedType := refType.In(i)

				if argExpectedType.Kind() != reflect.TypeOf(arg).Kind() {
					return nil, fmt.Errorf("function %s expected argument %d to be of type %s, got %s", call.Name, i, argExpectedType.Kind(), reflect.TypeOf(arg).Kind())
				}
				argsValues[i] = reflect.ValueOf(arg)
			}
			result := refValue.Call(argsValues)

			if len(result) == 0 {
				return nil, nil
			}

			if len(result) == 1 {
				return result[0].Interface(), nil
			}

			if len(result) == 2 {
				if result[1].Interface() != nil {
					return nil, result[1].Interface().(error)
				}
				return result[0].Interface(), nil
			}

			return nil, fmt.Errorf("function %s returned more than 2 values", call.Name)
		}
	}

	// evaluate arguments
	var args []interface{}
	for _, arg := range call.Arguments {
		argRes, err := v.evaluateExpression(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate argument: %w", err)
		}
		args = append(args, argRes)
	}

	// call function
	return fn(args...)
}
