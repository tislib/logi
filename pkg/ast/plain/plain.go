package plain

import (
	"github.com/tislib/logi/pkg/ast/common"
)

type Ast struct {
	SourceFile common.SourceFile `json:"sourceFile"`
	// The Macros of the package
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	MacroName string `json:"macroName"`
	Name      string `json:"name"`

	Statements []DefinitionStatement `json:"elements"`

	MacroNameSourceLocation common.SourceLocation `json:"macroNameSourceLocation"`
	NameSourceLocation      common.SourceLocation `json:"nameSourceLocation"`
}

type DefinitionStatement struct {
	Elements []DefinitionStatementElement `json:"elements"`

	SourceLocation common.SourceLocation `json:"sourceLocation"`
}

func (s DefinitionStatement) AsValue() common.Value {
	if len(s.Elements) == 0 {
		return common.NullValue()
	}

	if len(s.Elements) == 1 {
		return s.Elements[0].AsValue()
	}

	var arr []common.Value

	for _, e := range s.Elements {
		arr = append(arr, e.AsValue())
	}

	return common.ArrayValue(arr...)
}

type DefinitionStatementElementKind string

const (
	DefinitionStatementElementKindIdentifier    DefinitionStatementElementKind = "Identifier"
	DefinitionStatementElementKindValue         DefinitionStatementElementKind = "Value"
	DefinitionStatementElementKindArray         DefinitionStatementElementKind = "Array"
	DefinitionStatementElementKindStruct        DefinitionStatementElementKind = "Struct"
	DefinitionStatementElementKindAttributeList DefinitionStatementElementKind = "AttributeList"
	DefinitionStatementElementKindArgumentList  DefinitionStatementElementKind = "ArgumentList"
	DefinitionStatementElementKindParameterList DefinitionStatementElementKind = "ParameterList"
	DefinitionStatementElementKindExpression    DefinitionStatementElementKind = "Expression"
)

type DefinitionStatementElement struct {
	Kind DefinitionStatementElementKind `json:"kind"`

	Identifier    *DefinitionStatementElementIdentifier    `json:"identifier,omitempty"`
	Value         *DefinitionStatementElementValue         `json:"value,omitempty"`
	Array         *DefinitionStatementElementArray         `json:"array,omitempty"`
	Struct        *DefinitionStatementElementStruct        `json:"struct,omitempty"`
	AttributeList *DefinitionStatementElementAttributeList `json:"attributeList,omitempty"`
	ArgumentList  *DefinitionStatementElementArgumentList  `json:"argumentList,omitempty"`
	ParameterList *DefinitionStatementElementParameterList `json:"parameterList,omitempty"`
	Expression    *common.Expression                       `json:"expression,omitempty"`

	SourceLocation common.SourceLocation `json:"sourceLocation"`
}

func (e DefinitionStatementElement) AsValue() common.Value {
	switch e.Kind {
	case DefinitionStatementElementKindIdentifier:
		return e.Identifier.AsValue()
	case DefinitionStatementElementKindValue:
		return e.Value.Value
	case DefinitionStatementElementKindArray:
		return e.Array.AsValue()
	case DefinitionStatementElementKindAttributeList:
		return e.AttributeList.AsValue()
	case DefinitionStatementElementKindArgumentList:
		return e.ArgumentList.AsValue()
	case DefinitionStatementElementKindParameterList:
		return e.ParameterList.AsValue()
	}

	return common.NullValue()
}

type DefinitionStatementElementIdentifier struct {
	Identifier string `json:"identifier"`
}

func (i DefinitionStatementElementIdentifier) AsValue() common.Value {
	return common.StringValue(i.Identifier)
}

type DefinitionStatementElementValue struct {
	Value common.Value `json:"value"`
}

type DefinitionStatementElementArray struct {
	Items []DefinitionStatement `json:"items"`
}

type DefinitionStatementElementStruct struct {
	Statements []DefinitionStatement `json:"statements"`
}

func (l DefinitionStatementElementArray) AsValue() common.Value {
	var arr []common.Value

	for _, e := range l.Items {
		arr = append(arr, e.AsValue())
	}

	return common.ArrayValue(arr...)
}

type DefinitionStatementElementAttributeList struct {
	Attributes []DefinitionStatementElementAttribute `json:"attributes"`
}

func (l DefinitionStatementElementAttributeList) AsValue() common.Value {
	var valueMap = make(map[string]common.Value)

	for _, a := range l.Attributes {
		if a.Value != nil {
			valueMap[a.Name] = *a.Value
		} else {
			valueMap[a.Name] = common.BooleanValue(true)
		}
	}

	return common.MapValue(valueMap)
}

type DefinitionStatementElementAttribute struct {
	Name  string        `json:"name"`
	Value *common.Value `json:"value"`
}

type DefinitionStatementElementArgumentList struct {
	Arguments []DefinitionStatementElementArgument `json:"arguments"`
}

func (l DefinitionStatementElementArgumentList) AsValue() common.Value {
	var valueMap = make(map[string]common.Value)

	for _, a := range l.Arguments {
		valueMap[a.Name] = common.StringValue(a.Type.ToDisplayName())
	}

	return common.MapValue(valueMap)
}

type DefinitionStatementElementArgument struct {
	Name string                `json:"name"`
	Type common.TypeDefinition `json:"typeDefinition"`
}

type DefinitionStatementElementParameterList struct {
	Parameters []common.Expression `json:"parameters"`
}

func (l DefinitionStatementElementParameterList) AsValue() common.Value {
	var arr []common.Value

	for _, e := range l.Parameters {
		arr = append(arr, e.AsValue())
	}

	return common.ArrayValue(arr...)
}
