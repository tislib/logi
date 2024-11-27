package common

import (
	"fmt"
	"strings"
)

type ValueKind string

const (
	ValueKindString  ValueKind = "String"
	ValueKindBoolean ValueKind = "Boolean"
	ValueKindFloat   ValueKind = "Float"
	ValueKindInteger ValueKind = "Integer"
	ValueKindArray   ValueKind = "Array"
	ValueKindMap     ValueKind = "Map"
)

type Value struct {
	Kind ValueKind `json:"kind"`

	String  *string  `json:"string,omitempty"`
	Boolean *bool    `json:"boolean,omitempty"`
	Float   *float64 `json:"float,omitempty"`
	Integer *int64   `json:"integer,omitempty"`
	Array   []Value  `json:"array,omitempty"`
	Map     map[string]Value

	SourceLocation SourceLocation `json:"sourceLocation,omitempty"`
}

func (v Value) ToDisplayName() string {
	switch v.Kind {
	case ValueKindString:
		return *v.String
	case ValueKindBoolean:
		if *v.Boolean {
			return "true"
		}
		return "false"
	case ValueKindFloat:
		return fmt.Sprintf("%f", *v.Float)
	case ValueKindInteger:
		return fmt.Sprintf("%d", *v.Integer)
	case ValueKindArray:
		var result []string
		for _, value := range v.Array {
			result = append(result, value.ToDisplayName())
		}
		return strings.Join(result, ", ")
	case ValueKindMap:
		var result []string
		for key, value := range v.Map {
			result = append(result, fmt.Sprintf("%s: %s", key, value.ToDisplayName()))
		}
		return strings.Join(result, ", ")
	default:
		return "null"
	}
}

func (v Value) AsInterface() interface{} {
	switch v.Kind {
	case ValueKindString:
		return *v.String
	case ValueKindBoolean:
		return *v.Boolean
	case ValueKindFloat:
		return *v.Float
	case ValueKindInteger:
		return *v.Integer
	case ValueKindArray:
		var result []interface{}
		for _, value := range v.Array {
			result = append(result, value.AsInterface())
		}
		return result
	case ValueKindMap:
		var result = make(map[string]interface{})
		for key, value := range v.Map {
			result[key] = value.AsInterface()
		}
		return result
	default:
		return nil
	}
}

func StringValue(s string) Value {
	return Value{
		Kind:   ValueKindString,
		String: &s,
	}
}

func BooleanValue(b bool) Value {
	return Value{
		Kind:    ValueKindBoolean,
		Boolean: &b,
	}
}

func FloatValue(f float64) Value {
	return Value{
		Kind:  ValueKindFloat,
		Float: &f,
	}
}

func IntegerValue(i int64) Value {
	return Value{
		Kind:    ValueKindInteger,
		Integer: &i,
	}
}

func NullValue() Value {
	return Value{}
}

func PointerValue(value Value) *Value {
	return &value
}

func ArrayValue(arr ...Value) Value {
	return Value{
		Kind:  ValueKindArray,
		Array: arr,
	}
}

func MapValue(m map[string]Value) Value {
	return Value{
		Kind: ValueKindMap,
		Map:  m,
	}
}
