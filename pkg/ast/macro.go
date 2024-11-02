package ast

type MacroKind string

const (
	MacroKindSyntax MacroKind = "Syntax" // see SyntaxMacro
)

type MacroAst struct {
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

	// The Definition of the macro, used to define subtypes / elements of the macro, it will be used in the syntax section
	Definition Definition `json:"definition,omitempty"`

	// The Syntax of the macro, used to define the syntax of the macro, it will be used in the syntax section
	Syntax Syntax `json:"syntax,omitempty"`
}

type Definition struct {
	Statements []SyntaxStatement `json:"statements,omitempty"`
}

type Syntax struct {
	Statements []SyntaxStatement `json:"statements,omitempty"`
}

type SyntaxStatement struct {
	Elements []SyntaxStatementElement
}

type SyntaxStatementElementKind string

const (
	SyntaxStatementElementKindKeyword         SyntaxStatementElementKind = "Keyword"
	SyntaxStatementElementKindVariableKeyword SyntaxStatementElementKind = "VariableKeyword"
	SyntaxStatementElementKindParameterList   SyntaxStatementElementKind = "ParameterList"
	SyntaxStatementElementKindArgumentList    SyntaxStatementElementKind = "ArgumentList"
	SyntaxStatementElementKindCodeBlock       SyntaxStatementElementKind = "CodeBlock"
	SyntaxStatementElementKindMacro           SyntaxStatementElementKind = "Macro"
)

type SyntaxStatementElement struct {
	Kind SyntaxStatementElementKind `json:"kind,omitempty"`

	DefinitionRef   *SyntaxStatementElementDefinitionRef   `json:"definitionRef,omitempty"`
	ParameterList   *SyntaxStatementElementParameterList   `json:"parameterList,omitempty"`
	ArgumentList    *SyntaxStatementElementArgumentList    `json:"argumentList,omitempty"`
	KeywordDef      *SyntaxStatementElementKeywordDef      `json:"keywordDef,omitempty"`
	VariableKeyword *SyntaxStatementElementVariableKeyword `json:"variableKeyword,omitempty"`
	CodeBlock       *SyntaxStatementElementCodeBlock       `json:"codeBlock,omitempty"`
	Macro           *SyntaxStatementElementMacro           `json:"macro,omitempty"`
}

type SyntaxStatementElementDefinitionRef struct {
	Name string `json:"name,omitempty"`
}

type SyntaxStatementElementParameterList struct {
	Parameters []SyntaxStatementElementParameter `json:"parameters,omitempty"`
}

type SyntaxStatementElementArgumentList struct {
	VarArgs   bool                             `json:"varArgs,omitempty"`
	Arguments []SyntaxStatementElementArgument `json:"parameters,omitempty"`
}

type SyntaxStatementElementKeywordDef struct {
	Name string `json:"name,omitempty"`
}

type SyntaxStatementElementVariableKeyword struct {
	Name string         `json:"name,omitempty"`
	Type TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementParameter struct {
	Name string         `json:"name,omitempty"`
	Type TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementArgument struct {
	Name string         `json:"name,omitempty"`
	Type TypeDefinition `json:"type,omitempty"`
}

type SyntaxStatementElementCodeBlock struct {
	ReturnType TypeDefinition `json:"returnType,omitempty"`
}

type SyntaxStatementElementMacro struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
