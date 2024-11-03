package common

import "fmt"

type ValueKind string

const (
	ValueKindString  ValueKind = "String"
	ValueKindBoolean ValueKind = "Boolean"
	ValueKindFloat   ValueKind = "Float"
	ValueKindInteger ValueKind = "Integer"
)

type Value struct {
	Kind ValueKind `json:"kind"`

	String  *string  `json:"string,omitempty"`
	Boolean *bool    `json:"boolean,omitempty"`
	Float   *float64 `json:"float,omitempty"`
	Integer *int64   `json:"integer,omitempty"`
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
	}
	return ""
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

func PointerValue(value Value) *Value {
	return &value
}
