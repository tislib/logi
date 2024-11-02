package plain

import "logi/pkg/ast/common"

type Ast struct {
	// The Macros of the package
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	MacroName string `json:"macroName"`
	Name      string `json:"name"`

	Statements []DefinitionStatement `json:"elements"`
}

type DefinitionStatement struct {
	Elements []DefinitionStatementElement `json:"elements"`
}

type DefinitionStatementElementKind string

const (
	DefinitionStatementElementKindIdentifier    DefinitionStatementElementKind = "Identifier"
	DefinitionStatementElementKindValue         DefinitionStatementElementKind = "Value"
	DefinitionStatementElementKindAttributeList DefinitionStatementElementKind = "AttributeList"
)

type DefinitionStatementElement struct {
	Kind DefinitionStatementElementKind `json:"kind"`

	Identifier    *DefinitionStatementElementIdentifier    `json:"identifier,omitempty"`
	Value         *DefinitionStatementElementValue         `json:"value,omitempty"`
	AttributeList *DefinitionStatementElementAttributeList `json:"attributeList,omitempty"`
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
