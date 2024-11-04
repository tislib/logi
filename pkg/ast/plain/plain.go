package plain

import (
	"github.com/tislib/logi/pkg/ast/common"
)

type Ast struct {
	// The Macros of the package
	Definitions []Definition `json:"definitions"`
	Functions   []Function   `json:"functions"`
}

type Definition struct {
	MacroName string `json:"macroName"`
	Name      string `json:"name"`

	Statements []DefinitionStatement `json:"elements"`
}

type Function struct {
	Name      string                               `json:"name"`
	Arguments []DefinitionStatementElementArgument `json:"arguments"`
	CodeBlock common.CodeBlock
}

type DefinitionStatement struct {
	Elements []DefinitionStatementElement `json:"elements"`
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
	DefinitionStatementElementKindAttributeList DefinitionStatementElementKind = "AttributeList"
	DefinitionStatementElementKindArgumentList  DefinitionStatementElementKind = "ArgumentList"
	DefinitionStatementElementKindParameterList DefinitionStatementElementKind = "ParameterList"
	DefinitionStatementElementKindCodeBlock     DefinitionStatementElementKind = "CodeBlock"
)

type DefinitionStatementElement struct {
	Kind DefinitionStatementElementKind `json:"kind"`

	Identifier    *DefinitionStatementElementIdentifier    `json:"identifier,omitempty"`
	Value         *DefinitionStatementElementValue         `json:"value,omitempty"`
	Array         *DefinitionStatementElementArray         `json:"array,omitempty"`
	AttributeList *DefinitionStatementElementAttributeList `json:"attributeList,omitempty"`
	ArgumentList  *DefinitionStatementElementArgumentList  `json:"argumentList,omitempty"`
	ParameterList *DefinitionStatementElementParameterList `json:"parameterList,omitempty"`
	CodeBlock     *DefinitionStatementElementCodeBlock     `json:"codeBlock,omitempty"`
}

func (e DefinitionStatementElement) AsValue() common.Value {
	switch e.Kind {
	case DefinitionStatementElementKindIdentifier:
		return e.Identifier.AsValue()
	case DefinitionStatementElementKindValue:
		return e.Value.Value
	case DefinitionStatementElementKindAttributeList:
		return e.AttributeList.AsValue()
	case DefinitionStatementElementKindArgumentList:
		return e.ArgumentList.AsValue()
	case DefinitionStatementElementKindParameterList:
		return e.ParameterList.AsValue()
	case DefinitionStatementElementKindCodeBlock:
		return common.NullValue()
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
		valueMap[a.Name] = *a.Value
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
	Parameters []DefinitionStatementElementParameter `json:"parameters"`
}

func (l DefinitionStatementElementParameterList) AsValue() common.Value {
	var arr []common.Value

	for _, a := range l.Parameters {
		arr = append(arr, a.Value)
	}

	return common.ArrayValue(arr...)
}

type DefinitionStatementElementParameter struct {
	Value common.Value `json:"value"`
}

type DefinitionStatementElementCodeBlock struct {
	CodeBlock common.CodeBlock `json:"codeBlock"`
}
