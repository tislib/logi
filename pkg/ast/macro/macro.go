package macro

import "github.com/tislib/logi/pkg/ast/common"

type MacroKind string

const (
	KindSyntax MacroKind = "Syntax" // see SyntaxMacro
)

type Ast struct {
	// The Macros of the package
	Macros []Macro `json:"macros,omitempty"`
}

type Macro struct {
	// The Name of the macro, used to identify it, must be unique in the folder/package
	// required: true
	// example: "macro1"
	// pattern: ^[a-z][a-zA-Z0-9_]*$
	Name string `json:"name,omitempty"`

	// The Description of the macro, used to describe what the macro does, it must be written immediately before the macro
	Comment string `json:"comment,omitempty"`

	// The Kind of the macro, used to categorize the macro
	// For each kind of macro, a different struct will be used
	Kind MacroKind `json:"kind,omitempty"`

	// The Types of the macro, used to define the types of the macro, it will be used in the syntax section
	Types Types `json:"definition,omitempty"`

	// The Syntax of the macro, used to define the syntax of the macro, it will be used in the syntax section
	Syntax Syntax `json:"syntax,omitempty"`

	SourceMap map[string]common.SourceLocation `json:"sourceMap,omitempty"`
}

type Types struct {
	Types []TypeStatement `json:"types,omitempty"`
}

type Syntax struct {
	Statements []SyntaxStatement `json:"statements,omitempty"`
}

type SyntaxStatement struct {
	Elements []SyntaxStatementElement
}

type TypeStatement struct {
	Name     string                   `json:"name,omitempty"`
	Elements []SyntaxStatementElement `json:"elements,omitempty"`
}

type SyntaxStatementElementKind string

const (
	SyntaxStatementElementKindKeyword         SyntaxStatementElementKind = "Keyword"
	SyntaxStatementElementKindTypeReference   SyntaxStatementElementKind = "TypeReference"
	SyntaxStatementElementKindVariableKeyword SyntaxStatementElementKind = "VariableKeyword"
	SyntaxStatementElementKindCombination     SyntaxStatementElementKind = "Combination"
	SyntaxStatementElementKindStructure       SyntaxStatementElementKind = "Structure"
	SyntaxStatementElementKindParameterList   SyntaxStatementElementKind = "ParameterList"
	SyntaxStatementElementKindArgumentList    SyntaxStatementElementKind = "ArgumentList"
	SyntaxStatementElementKindCodeBlock       SyntaxStatementElementKind = "CodeBlock"
	SyntaxStatementElementKindExpressionBlock SyntaxStatementElementKind = "ExpressionBlock"
	SyntaxStatementElementKindAttributeList   SyntaxStatementElementKind = "AttributeList"
)

type SyntaxStatementElement struct {
	Kind SyntaxStatementElementKind `json:"kind,omitempty"`

	KeywordDef      *SyntaxStatementElementKeywordDef      `json:"keywordDef,omitempty"`
	TypeReference   *SyntaxStatementElementTypeReference   `json:"typeReference,omitempty"`
	VariableKeyword *SyntaxStatementElementVariableKeyword `json:"variableKeyword,omitempty"`
	Combination     *SyntaxStatementElementCombination     `json:"combination,omitempty"`
	Structure       *SyntaxStatementElementStructure       `json:"structure,omitempty"`
	ParameterList   *SyntaxStatementElementParameterList   `json:"parameterList,omitempty"`
	ArgumentList    *SyntaxStatementElementArgumentList    `json:"argumentList,omitempty"`
	CodeBlock       *SyntaxStatementElementCodeBlock       `json:"codeBlock,omitempty"`
	ExpressionBlock *SyntaxStatementElementExpressionBlock `json:"expressionBlock,omitempty"`
	AttributeList   *SyntaxStatementElementAttributeList   `json:"attributeList,omitempty"`
}

type SyntaxStatementElementCombination struct {
	Elements []SyntaxStatementElement `json:"elements,omitempty"`
}

type SyntaxStatementElementStructure struct {
	Statements []SyntaxStatement `json:"statements,omitempty"`
}

type SyntaxStatementElementParameterList struct {
	Parameters []SyntaxStatementElementParameter `json:"parameters,omitempty"`
}

type SyntaxStatementElementAttributeList struct {
	Attributes []SyntaxStatementElementAttribute `json:"attributes,omitempty"`
}

type SyntaxStatementElementArgumentList struct {
	VarArgs   bool                             `json:"varArgs,omitempty"`
	Arguments []SyntaxStatementElementArgument `json:"parameters,omitempty"`
}

type SyntaxStatementElementKeywordDef struct {
	Name string `json:"name,omitempty"`
}

type SyntaxStatementElementTypeReference struct {
	Name string `json:"name,omitempty"`
}

type SyntaxStatementElementVariableKeyword struct {
	Name string                `json:"name,omitempty"`
	Type common.TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementParameter struct {
	Name string                `json:"name,omitempty"`
	Type common.TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementArgument struct {
	Name string                `json:"name,omitempty"`
	Type common.TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementCodeBlock struct {
}

type SyntaxStatementElementExpressionBlock struct {
}

type SyntaxStatementElementAttribute struct {
	Name string                `json:"name,omitempty"`
	Type common.TypeDefinition `json:"type,omitempty"`
}
