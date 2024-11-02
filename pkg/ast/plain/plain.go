package plain

type LogiPlainAst struct {
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

type DefinitionStatementElement struct {
	Kind MacroSyntaxStatementElementKind `json:"kind"`
}
