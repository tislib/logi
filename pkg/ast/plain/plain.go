package plain

import (
	"logi/pkg/ast/common"
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

type DefinitionStatementElementKind string

const (
	DefinitionStatementElementKindIdentifier    DefinitionStatementElementKind = "Identifier"
	DefinitionStatementElementKindValue         DefinitionStatementElementKind = "Value"
	DefinitionStatementElementKindAttributeList DefinitionStatementElementKind = "AttributeList"
	DefinitionStatementElementKindArgumentList  DefinitionStatementElementKind = "ArgumentList"
	DefinitionStatementElementKindCodeBlock     DefinitionStatementElementKind = "CodeBlock"
)

type DefinitionStatementElement struct {
	Kind DefinitionStatementElementKind `json:"kind"`

	Identifier    *DefinitionStatementElementIdentifier    `json:"identifier,omitempty"`
	Value         *DefinitionStatementElementValue         `json:"value,omitempty"`
	AttributeList *DefinitionStatementElementAttributeList `json:"attributeList,omitempty"`
	ArgumentList  *DefinitionStatementElementArgumentList  `json:"argumentList,omitempty"`
	CodeBlock     *DefinitionStatementElementCodeBlock     `json:"codeBlock,omitempty"`
}

type DefinitionStatementElementIdentifier struct {
	Identifier string `json:"identifier"`
}

type DefinitionStatementElementValue struct {
	Value common.Value `json:"value"`
}

type DefinitionStatementElementAttributeList struct {
	Attributes []DefinitionStatementElementAttribute `json:"attributes"`
}

type DefinitionStatementElementAttribute struct {
	Name  string        `json:"name"`
	Value *common.Value `json:"value"`
}

type DefinitionStatementElementArgumentList struct {
	Arguments []DefinitionStatementElementArgument `json:"arguments"`
}

type DefinitionStatementElementArgument struct {
	Name string                `json:"name"`
	Type common.TypeDefinition `json:"typeDefinition"`
}

type DefinitionStatementElementCodeBlock struct {
	CodeBlock common.CodeBlock `json:"codeBlock"`
}
