package macro

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
	"strings"
)

type converter struct {
	enableSourceMap bool
}

func (c *converter) convertNodeToMacroAst(node yaccNode) (*astMacro.Ast, error) {
	var res = new(astMacro.Ast)

	if node.op != NodeOpFile {
		return res, ErrUnexpectedNode
	}

	for _, child := range node.children {
		macro, err := c.convertMacro(child)

		if err != nil {
			return res, err
		}

		res.Macros = append(res.Macros, *macro)
	}

	return res, nil
}

func (c *converter) convertMacro(macroNode yaccNode) (*astMacro.Macro, error) {
	var signature = macroNode.children[0]
	var name = signature.children[0]
	var body = macroNode.children[1]
	var kind = body.children[0].value.(string)

	var result = new(astMacro.Macro)

	// source maps
	if c.enableSourceMap {
		result.SourceMap = make(map[string]common.SourceLocation)
		result.SourceMap["macro"] = signature.location.AsSourceLocation()
		result.SourceMap["name"] = name.location.AsSourceLocation()
	}

	if !NamePattern.MatchString(name.value.(string)) {
		return nil, fmt.Errorf("unexpected name value: %s", name.value)
	}

	result.Name = name.value.(string)
	switch kind {
	case "Syntax":
		result.Kind = astMacro.KindSyntax
	default:
		return result, c.newErrorFromNode(body.children[0], fmt.Sprintf("unexpected kind value: \"%s\", expecting \"Syntax\"", kind))
	}

	for _, child := range body.children {
		switch child.op {
		case NodeOpSyntax:
			if c.enableSourceMap {
				result.SourceMap["syntax"] = child.location.AsSourceLocation()
			}

			if result.Kind != astMacro.KindSyntax {
				return result, fmt.Errorf("syntax defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				syntaxBody, err := c.convertSyntaxBody(child.children[0])

				if err != nil {
					return result, err
				}

				result.Syntax = astMacro.Syntax{Statements: syntaxBody}
			}
		case NodeOpTypes:
			if c.enableSourceMap {
				result.SourceMap["types"] = child.location.AsSourceLocation()
			}
			if result.Kind != astMacro.KindSyntax {
				return result, fmt.Errorf("types defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				types, err := c.convertTypes(child.children[0])

				if err != nil {
					return result, err
				}

				result.Types = *types
			}
		case NodeOpScopes:
			if c.enableSourceMap {
				result.SourceMap["scopes"] = child.location.AsSourceLocation()
			}
			if result.Kind != astMacro.KindSyntax {
				return result, fmt.Errorf("scopes defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				scopes, err := c.convertScopes(child.children[0])

				if err != nil {
					return result, err
				}

				result.Scopes = *scopes
			}
		}
	}

	return result, nil
}

func (c *converter) newErrorFromNode(name yaccNode, msg string) error {
	return newError(name.location.Line, name.location.Column, name.token.Value, msg)
}

func (c *converter) convertSyntaxBody(syntaxNode yaccNode) ([]astMacro.SyntaxStatement, error) {
	if syntaxNode.children == nil {
		return nil, nil
	}

	var result []astMacro.SyntaxStatement

	for _, child := range syntaxNode.children {
		statement, err := c.convertSyntaxStatement(child)

		if err != nil {
			return nil, err
		}

		result = append(result, *statement)
	}

	return result, nil
}

func (c *converter) convertTypes(typesNode yaccNode) (*astMacro.Types, error) {
	if typesNode.children == nil {
		return nil, nil
	}

	var result []astMacro.TypeStatement

	for _, child := range typesNode.children {
		statement, err := c.convertTypeStatement(child)

		if err != nil {
			return nil, err
		}

		result = append(result, *statement)
	}

	return &astMacro.Types{Types: result}, nil
}

func (c *converter) convertScopes(scopeNode yaccNode) (*astMacro.Scopes, error) {
	if scopeNode.children == nil {
		return nil, nil
	}

	var result []astMacro.ScopeItem

	for _, child := range scopeNode.children {
		statement, err := c.convertScopeItem(child)

		if err != nil {
			return nil, err
		}

		result = append(result, *statement)
	}

	return &astMacro.Scopes{Scopes: result}, nil
}

func (c *converter) convertScopeItem(node yaccNode) (*astMacro.ScopeItem, error) {
	var result = new(astMacro.ScopeItem)
	var name = node.children[0].value.(string)

	body, err := c.convertSyntaxBody(node.children[1])

	if err != nil {
		return nil, fmt.Errorf("unexpected scope item op: %s", node.op)
	}

	result.Name = name
	result.Statements = body

	return result, nil
}

func (c *converter) convertTypeStatement(node yaccNode) (*astMacro.TypeStatement, error) {
	var result = new(astMacro.TypeStatement)
	var name = node.children[0].value.(string)
	var body = node.children[1].children

	result.Name = name

	for _, item := range body {
		switch item.op {
		case NodeOpSyntaxElements:
			for _, child := range item.children {
				element, err := c.convertSyntaxStatementElement(child)

				if err != nil {
					return nil, err
				}

				result.Elements = append(result.Elements, *element)
			}
		default:
			return nil, fmt.Errorf("unexpected type statement op: %s", item.op)
		}
	}

	return result, nil
}

func (c *converter) convertSyntaxStatement(body yaccNode) (*astMacro.SyntaxStatement, error) {
	var result = new(astMacro.SyntaxStatement)

	for _, item := range body.children {
		switch item.op {
		case NodeOpSyntaxElements:
			for _, child := range item.children {
				element, err := c.convertSyntaxStatementElement(child)

				if err != nil {
					return nil, err
				}

				result.Elements = append(result.Elements, *element)
			}
		case NodeOpSyntaxExamples:
			var examples []string

			for _, child := range item.children {
				examples = append(examples, c.convertValue(child))
			}

			result.Examples = examples
		default:
			return nil, fmt.Errorf("unexpected syntax statement op: %s", item.op)
		}
	}

	return result, nil
}

func (c *converter) convertValue(child yaccNode) string {
	var parts []string

	for _, item := range child.children {
		switch item.op {
		case NodeOpValueIdentifier:
			parts = append(parts, item.value.(string))
		case NodeOpValueString:
			parts = append(parts, fmt.Sprintf("\"%s\"", item.value.(string)))
		case NodeOpValueArray:
			var array []string

			for _, arrayItem := range item.children {
				array = append(array, c.convertValue(arrayItem))
			}

			parts = append(parts, fmt.Sprintf("[%s]", strings.Join(array, ", ")))
		default:
			parts = append(parts, fmt.Sprintf("%v", item.value))
		}
	}

	return strings.Join(parts, " ")
}

func (c *converter) convertSyntaxStatementElement(node yaccNode) (*astMacro.SyntaxStatementElement, error) {
	var result = new(astMacro.SyntaxStatementElement)

	switch node.op {
	case NodeOpSyntaxKeywordElement:
		result.Kind = astMacro.SyntaxStatementElementKindKeyword
		result.KeywordDef = &astMacro.SyntaxStatementElementKeywordDef{Name: node.value.(string)}
	case NodeOpSyntaxSymbolElement:
		result.Kind = astMacro.SyntaxStatementElementKindSymbol
		result.SymbolDef = &astMacro.SyntaxStatementElementSymbolDef{Name: node.value.(string)}
	case NodeOpSyntaxTypeReferenceElement:
		result.Kind = astMacro.SyntaxStatementElementKindTypeReference
		result.TypeReference = &astMacro.SyntaxStatementElementTypeReference{Name: node.value.(string)}
	case NodeOpSyntaxVariableKeywordElement:
		result.Kind = astMacro.SyntaxStatementElementKindVariableKeyword
		var varName = node.children[0].value.(string)
		typeDef, err := c.convertTypeDefinition(node.children[1])

		if err != nil {
			return nil, err
		}

		result.VariableKeyword = &astMacro.SyntaxStatementElementVariableKeyword{Name: varName, Type: *typeDef}
	case NodeOpSyntaxCombinationElement:
		result.Kind = astMacro.SyntaxStatementElementKindCombination

		var elements []astMacro.SyntaxStatementElement
		for _, elementNode := range node.children {
			element, err := c.convertSyntaxStatementElement(elementNode)

			if err != nil {
				return nil, err
			}

			elements = append(elements, *element)
		}

		result.Combination = &astMacro.SyntaxStatementElementCombination{Elements: elements}
	case NodeOpSyntaxParameterListElement:
		result.Kind = astMacro.SyntaxStatementElementKindParameterList

		if node.value == true {
			result.ParameterList = &astMacro.SyntaxStatementElementParameterList{Dynamic: true}
			break
		}

		var parameters []astMacro.SyntaxStatementElementParameter

		for _, parameterNode := range node.children {
			parameter, err := c.convertSyntaxStatementElementParameter(parameterNode)

			if err != nil {
				return nil, err
			}

			parameters = append(parameters, *parameter)
		}

		result.ParameterList = &astMacro.SyntaxStatementElementParameterList{Parameters: parameters}
	case NodeOpSyntaxArgumentListElement:
		result.Kind = astMacro.SyntaxStatementElementKindArgumentList

		var arguments []astMacro.SyntaxStatementElementArgument

		for _, argumentNode := range node.children {
			argument, err := c.convertSyntaxStatementElementArgument(argumentNode)

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, *argument)
		}

		result.ArgumentList = &astMacro.SyntaxStatementElementArgumentList{Arguments: arguments, VarArgs: true}
	case NodeOpSyntaxAttributeListElement:
		result.Kind = astMacro.SyntaxStatementElementKindAttributeList

		var attributes []astMacro.SyntaxStatementElementAttribute

		for _, attributeNode := range node.children {
			attribute, err := c.convertSyntaxStatementElementAttribute(attributeNode)

			if err != nil {
				return nil, err
			}

			attributes = append(attributes, *attribute)
		}

		result.AttributeList = &astMacro.SyntaxStatementElementAttributeList{Attributes: attributes}
	case NodeOpSyntaxScopeElement:
		result.Kind = astMacro.SyntaxStatementElementKindScope
		result.ScopeDef = &astMacro.SyntaxStatementElementScopeDef{}

		for _, scopeNode := range node.children {
			result.ScopeDef.Scopes = append(result.ScopeDef.Scopes, scopeNode.value.(string))
		}
	default:
		return nil, fmt.Errorf("unexpected syntax statement element op: %s", node.op)
	}

	return result, nil
}

func (c *converter) convertSyntaxStatementElementParameter(node yaccNode) (*astMacro.SyntaxStatementElementParameter, error) {
	var result = new(astMacro.SyntaxStatementElementParameter)
	var varName = node.children[0].value.(string)
	typeDef, err := c.convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, err
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func (c *converter) convertSyntaxStatementElementArgument(node yaccNode) (*astMacro.SyntaxStatementElementArgument, error) {
	var result = new(astMacro.SyntaxStatementElementArgument)
	var varName = node.children[0].value.(string)
	typeDef, err := c.convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, err
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func (c *converter) convertSyntaxStatementElementAttribute(node yaccNode) (*astMacro.SyntaxStatementElementAttribute, error) {
	var result = new(astMacro.SyntaxStatementElementAttribute)
	var varName = node.value.(string)

	result.Name = varName

	if len(node.children) > 0 {
		typeDef, err := c.convertTypeDefinition(node.children[0])

		if err != nil {
			return nil, err
		}

		result.Type = *typeDef
	}

	return result, nil
}

func (c *converter) convertTypeDefinition(node yaccNode) (*common.TypeDefinition, error) {
	var result = new(common.TypeDefinition)
	result.Name = node.value.(string)

	if len(node.children) > 0 {
		for _, child := range node.children {
			subType, err := c.convertTypeDefinition(child)

			if err != nil {
				return nil, err
			}

			result.SubTypes = append(result.SubTypes, *subType)
		}
	}

	return result, nil
}
