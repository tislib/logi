package logi

import macroAst "github.com/tislib/logi/pkg/ast/macro"

type MacroParseGraph struct {
	macroAst *macroAst.Macro
}

func (g MacroParseGraph) Prepare() *ParseNode {
	return g.sequence(
		g.eolAllowed(),

		// signature
		g.tokenNode(token_identifier, g.macroAst.Name, nil),
		g.tokenNode(token_identifier, "", nil),

		// body
		g.tokenNode(BraceOpen, "", nil),
		g.eolRequired(),
		g.statements(),
		g.tokenNode(BraceClose, "", nil),

		g.eolAllowed(),
	)
}

func (g MacroParseGraph) tokenNode(identifier int, name string, visitFunc ParseNodeVisitFunc) *ParseNode {
	return &ParseNode{
		TokenId:    identifier,
		TokenValue: name,
		VisitFunc:  visitFunc,
		Mode:       ModeToken,
	}
}

func (g MacroParseGraph) tokenValue() func(node *ParseNode) error {
	return func(node *ParseNode) error {
		node.Value = node.TokenValue

		return nil
	}
}

func (g MacroParseGraph) sequence(nodes ...*ParseNode) *ParseNode {
	var node = &ParseNode{
		Mode: ModeSequence,
	}

	for _, n := range nodes {
		node.Children = append(node.Children, n)
	}

	return node
}

func (g MacroParseGraph) or(nodes ...*ParseNode) *ParseNode {
	var node = &ParseNode{
		Mode: ModeOr,
	}

	for _, n := range nodes {
		node.Children = append(node.Children, n)
	}

	return node
}

func (g MacroParseGraph) multiple(node *ParseNode) *ParseNode {
	var multipleNode = g.or(node)

	multipleNode.Children = append(multipleNode.Children, multipleNode)

	return multipleNode
}

func (g MacroParseGraph) optional(node *ParseNode) *ParseNode {
	return g.or(node, nil)
}

func (g MacroParseGraph) eolRequired() *ParseNode {
	return g.multiple(g.eol())
}

func (g MacroParseGraph) eolAllowed() *ParseNode {
	return g.optional(g.multiple(g.eol()))
}

func (g MacroParseGraph) eol() *ParseNode {
	return g.tokenNode(Eol, "", nil)
}

func (g MacroParseGraph) statements() *ParseNode {
	return g.multiple(g.statementNode())
}

func (g MacroParseGraph) statementNode() *ParseNode {
	var items []*ParseNode

	for _, s := range g.macroAst.Syntax.Statements {
		items = append(items, g.statementVariantNode(s))
	}

	return g.or(items...)
}

func (g MacroParseGraph) statementVariantNode(s macroAst.SyntaxStatement) *ParseNode {
	var items []*ParseNode

	for _, e := range s.Elements {
		items = append(items, g.statementVariantElementNode(e))
	}

	return g.sequence(items...)
}

func (g MacroParseGraph) statementVariantElementNode(e macroAst.SyntaxStatementElement) *ParseNode {
	// 	SyntaxStatementElementKindKeyword         SyntaxStatementElementKind = "Keyword"
	//	SyntaxStatementElementKindTypeReference   SyntaxStatementElementKind = "TypeReference"
	//	SyntaxStatementElementKindVariableKeyword SyntaxStatementElementKind = "VariableKeyword"
	//	SyntaxStatementElementKindCombination     SyntaxStatementElementKind = "Combination"
	//	SyntaxStatementElementKindStructure       SyntaxStatementElementKind = "Combination"
	//	SyntaxStatementElementKindParameterList   SyntaxStatementElementKind = "ParameterList"
	//	SyntaxStatementElementKindArgumentList    SyntaxStatementElementKind = "ArgumentList"
	//	SyntaxStatementElementKindCodeBlock       SyntaxStatementElementKind = "CodeBlock"
	//	SyntaxStatementElementKindExpressionBlock SyntaxStatementElementKind = "ExpressionBlock"
	//	SyntaxStatementElementKindAttributeList   SyntaxStatementElementKind = "AttributeList"

	switch e.Kind {
	case macroAst.SyntaxStatementElementKindKeyword:
		return g.tokenNode(token_identifier, e.KeywordDef.Name, nil)
	case macroAst.SyntaxStatementElementKindTypeReference:
		return g.typeReferenceNode(e.TypeReference)
	case macroAst.SyntaxStatementElementKindVariableKeyword:
		return g.tokenNode(token_identifier, e.VariableKeyword.Name, nil)
	default:
		panic("unknown syntax statement element kind")
	}
}

func (g MacroParseGraph) typeReferenceNode(reference *macroAst.SyntaxStatementElementTypeReference) *ParseNode {
	return nil
}

func NewMacroParseGraph(mAst *macroAst.Macro) *MacroParseGraph {
	return &MacroParseGraph{macroAst: mAst}
}
