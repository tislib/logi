package logi

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
)

type recursiveStatementParser struct {
	plainStatement  plain.DefinitionStatement
	syntaxStatement macroAst.SyntaxStatement
	macroDefinition *macroAst.Macro

	maxMatch      int
	bestMatch     *macroAst.SyntaxStatement
	mismatchCause string

	sei       int
	pei       int
	statement logiAst.Statement
}

func (p *recursiveStatementParser) parse(scope string) error {
	var maxMatch = -1
	var mismatchCause string

	for _, syntaxStatement := range p.macroDefinition.Syntax.Statements {
		p.sei = 0
		p.pei = 0
		p.maxMatch = -1
		p.mismatchCause = ""
		p.syntaxStatement = syntaxStatement
		p.statement = logiAst.Statement{
			Scope: scope,
		}

		p.match()

		if p.mismatchCause != "" {
			if p.maxMatch > maxMatch {
				maxMatch = p.maxMatch
				mismatchCause = p.mismatchCause
			}
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to match statement: %s", mismatchCause)
}

func (p *recursiveStatementParser) reportMismatch(reason string) {
	if p.pei >= p.maxMatch || p.mismatchCause == "" {
		p.maxMatch = p.pei
		p.bestMatch = &p.syntaxStatement
		p.mismatchCause = reason
	}
}

func (p *recursiveStatementParser) match() {
	for p.sei < len(p.syntaxStatement.Elements) {
		var syntaxStatementElement = p.syntaxStatement.Elements[p.sei]

		// check if the statement is always required
		if p.pei >= len(p.plainStatement.Elements) {
			if isSyntaxElementAlwaysRequired(syntaxStatementElement) {
				p.reportMismatch("statement is shorter than syntax")
				return
			}

			p.sei++
			return
		}

		var currentElement = p.plainStatement.Elements[p.pei]

		p.matchNextElement(syntaxStatementElement, currentElement)

		if p.mismatchCause != "" {
			return // stop matching if a mismatch is found
		}

		p.sei++
		p.pei++
	}

	if p.pei < len(p.plainStatement.Elements) {
		p.reportMismatch(fmt.Sprintf("plain statement has more elements than syntax statement: %s", p.plainStatement.Elements[p.pei].Kind))
	}
}

func (p *recursiveStatementParser) matchNextElement(syntaxStatementElement macroAst.SyntaxStatementElement, currentElement plain.DefinitionStatementElement) {
	log.Trace(fmt.Sprintf("matching %s with %s at: %s", syntaxStatementElement.Kind, currentElement.Kind, currentElement.SourceLocation))
	log.Trace(fmt.Sprintf("Current element: %v", currentElement.AsValue().AsInterface()))

	switch syntaxStatementElement.Kind {
	case macroAst.SyntaxStatementElementKindKeyword:
		if currentElement.Kind != plain.DefinitionStatementElementKindIdentifier {
			p.reportMismatch(fmt.Sprintf("expected keyword (%s), got %s (%v)", syntaxStatementElement.KeywordDef.Name, currentElement.Kind, currentElement))
			break
		} else if currentElement.Identifier.Identifier != syntaxStatementElement.KeywordDef.Name {
			p.reportMismatch(fmt.Sprintf("expected keyword (%s), got %s", syntaxStatementElement.KeywordDef.Name, currentElement.Identifier.Identifier))
			break
		}
		if p.statement.Command == "" && p.pei == 0 {
			p.statement.Command = syntaxStatementElement.KeywordDef.Name
		}
	case macroAst.SyntaxStatementElementKindSymbol:
		if currentElement.Kind != plain.DefinitionStatementElementKindSymbol {
			p.reportMismatch(fmt.Sprintf("expected keyword (%s), got %s (%v)", syntaxStatementElement.KeywordDef.Name, currentElement.Kind, currentElement))
			break
		} else if currentElement.Symbol.Symbol != syntaxStatementElement.SymbolDef.Name {
			p.reportMismatch(fmt.Sprintf("expected symbol (%s), got %s", syntaxStatementElement.KeywordDef.Name, currentElement.Identifier.Identifier))
			break
		}
	case macroAst.SyntaxStatementElementKindTypeReference:
		p.matchTypeReference(syntaxStatementElement, currentElement)
	case macroAst.SyntaxStatementElementKindVariableKeyword:
		switch currentElement.Kind {
		case plain.DefinitionStatementElementKindIdentifier:
			p.statement.Parameters = append(p.statement.Parameters, logiAst.Parameter{
				Name:  syntaxStatementElement.VariableKeyword.Name,
				Value: common.StringValue(currentElement.Identifier.Identifier),
			})
			if syntaxStatementElement.VariableKeyword.Type.Name != "Name" && syntaxStatementElement.VariableKeyword.Type.Name != "Type" {
				p.reportMismatch(fmt.Sprintf("expected variable keyword %s in type %s, got %s in type Name", syntaxStatementElement.VariableKeyword.Name, syntaxStatementElement.VariableKeyword.Type, currentElement.Identifier.Identifier))
				break
			}
		case plain.DefinitionStatementElementKindArray:
			p.matchArray(currentElement, syntaxStatementElement)
			return
		case plain.DefinitionStatementElementKindValue:
			p.matchValue(syntaxStatementElement)
		default:
			p.reportMismatch(fmt.Sprintf("expected variable keyword (%s), got %s", syntaxStatementElement.VariableKeyword.Name, currentElement.Kind))
		}
	case macroAst.SyntaxStatementElementKindAttributeList:
		if currentElement.Kind != plain.DefinitionStatementElementKindArray {
			p.reportMismatch(fmt.Sprintf("expected attribute list, got %s", currentElement.Kind))

			break
		}

		var attributeMap = make(map[string]macroAst.SyntaxStatementElementAttribute)
		for _, attr := range syntaxStatementElement.AttributeList.Attributes {
			attributeMap[attr.Name] = attr
		}

		for _, item := range currentElement.Array.Items {
			if len(item.Elements) == 1 {
				var elem0 = item.Elements[0]

				if elem0.Kind != plain.DefinitionStatementElementKindIdentifier {
					p.reportMismatch(fmt.Sprintf("expected attribute list, got %v", currentElement.AsValue().AsInterface()))
					break
				}

				attr, ok := attributeMap[elem0.Identifier.Identifier]

				if !ok {
					p.reportMismatch(fmt.Sprintf("attribute %s not found", elem0.Identifier.Identifier))
					break
				}

				if attr.Type.Name != "bool" {
					p.reportMismatch(fmt.Sprintf("expected bool attribute got %s", attr.Type.Name))
					break
				}

				p.statement.Attributes = append(p.statement.Attributes, logiAst.Attribute{
					Name: elem0.Identifier.Identifier,
				})
			} else if len(item.Elements) == 2 {
				var elem0 = item.Elements[0]
				var elem1 = item.Elements[1]

				if elem0.Kind != plain.DefinitionStatementElementKindIdentifier {
					p.reportMismatch(fmt.Sprintf("expected attribute list, got %v", currentElement.AsValue().AsInterface()))
					break
				}

				_, ok := attributeMap[elem0.Identifier.Identifier]

				if !ok {
					p.reportMismatch(fmt.Sprintf("attribute %s not found", elem0.Identifier.Identifier))
					break
				}

				p.statement.Attributes = append(p.statement.Attributes, logiAst.Attribute{
					Name:  elem0.Identifier.Identifier,
					Value: common.PointerValue(elem1.AsValue()),
				})
			} else {
				p.reportMismatch(fmt.Sprintf("expected attribute list, got %v", currentElement.AsValue().AsInterface()))
				break
			}
		}
	case macroAst.SyntaxStatementElementKindArgumentList:
		if currentElement.Kind != plain.DefinitionStatementElementKindArgumentList {
			p.reportMismatch(fmt.Sprintf("expected argument list, got %s", currentElement.Kind))

			break
		}

		// check if the argument list is valid
		err := isValidArgumentList(currentElement.ArgumentList, syntaxStatementElement.ArgumentList)

		if err != nil {
			p.reportMismatch(fmt.Sprintf(err.Error()))

			break
		}
		for _, argument := range currentElement.ArgumentList.Arguments {
			p.statement.Arguments = append(p.statement.Arguments, logiAst.Argument{
				Name: argument.Name,
				Type: argument.Type,
			})
		}
	case macroAst.SyntaxStatementElementKindParameterList:
		if currentElement.Kind != plain.DefinitionStatementElementKindParameterList {
			p.reportMismatch(fmt.Sprintf("expected parameter list, got %s", currentElement.Kind))

			break
		}

		if len(currentElement.ParameterList.Names) != 0 {
			parameterNameIdx := make(map[string]int)
			parameterChecked := make(map[string]bool)

			if !syntaxStatementElement.ParameterList.Dynamic {
				for idx, name := range currentElement.ParameterList.Names {
					parameterNameIdx[name] = idx
				}

				for _, syntaxStatementElementParameter := range syntaxStatementElement.ParameterList.Parameters {
					idx, ok := parameterNameIdx[syntaxStatementElementParameter.Name]

					if !ok {
						continue
					}

					var param = currentElement.ParameterList.Parameters[idx]

					p.statement.Parameters = append(p.statement.Parameters, logiAst.Parameter{
						Name:       syntaxStatementElementParameter.Name,
						Value:      param.AsValue(),
						Expression: &param,
					})

					parameterChecked[syntaxStatementElementParameter.Name] = true
				}

				for name := range parameterNameIdx {
					if !parameterChecked[name] {
						p.reportMismatch(fmt.Sprintf("parameter %s not found", name))
						break
					}
				}
			} else {
				for idx, name := range currentElement.ParameterList.Names {
					var param = currentElement.ParameterList.Parameters[idx]

					p.statement.Parameters = append(p.statement.Parameters, logiAst.Parameter{
						Name:       name,
						Value:      param.AsValue(),
						Expression: &param,
					})
				}
			}
		} else {
			if len(syntaxStatementElement.ParameterList.Parameters) > 0 && syntaxStatementElement.ParameterList.Dynamic {
				p.reportMismatch(fmt.Sprintf("positional parameters are not allowed in dynamic parameter list"))
				break
			}
			for idx, syntaxStatementElementParameter := range syntaxStatementElement.ParameterList.Parameters {
				if idx >= len(currentElement.ParameterList.Parameters) {
					break
				}
				var param = currentElement.ParameterList.Parameters[idx]
				p.statement.Parameters = append(p.statement.Parameters, logiAst.Parameter{
					Name:       syntaxStatementElementParameter.Name,
					Value:      param.AsValue(),
					Expression: &param,
				})
			}
		}
	case macroAst.SyntaxStatementElementKindCombination:
		p.matchCombination(syntaxStatementElement)
	case macroAst.SyntaxStatementElementKindScope:
		p.matchScope(currentElement, syntaxStatementElement)
	default:
		p.reportMismatch(fmt.Sprintf("unexpected syntax element kind: %s", syntaxStatementElement.Kind))
	}
}

func (p *recursiveStatementParser) matchValue(syntaxStatementElement macroAst.SyntaxStatementElement) {
	for _, typeStatement := range p.macroDefinition.Types.Types {
		if typeStatement.Name == syntaxStatementElement.VariableKeyword.Type.Name {
			// start matching the value
			p.matchTypeStatement(syntaxStatementElement.VariableKeyword.Name, typeStatement)
			return
		}
	}

	currentElement := p.plainStatement.Elements[p.pei]
	var value common.Value
	if currentElement.Kind == plain.DefinitionStatementElementKindValue {
		value = currentElement.Value.Value
	} else if currentElement.Kind == plain.DefinitionStatementElementKindIdentifier {
		value = common.StringValue(currentElement.Identifier.Identifier)
	} else {
		p.reportMismatch(fmt.Sprintf("expected value, got %s", currentElement.Kind))
		return

	}
	currentElementKind := value.Kind

	switch syntaxStatementElement.VariableKeyword.Type.Name {
	case "int":
		if currentElementKind != common.ValueKindInteger {
			p.reportMismatch(fmt.Sprintf("expected int got %s", currentElementKind))
			return
		}
	case "string":
		if currentElementKind != common.ValueKindString {
			p.reportMismatch(fmt.Sprintf("expected string got %s", currentElementKind))
			return
		}
	case "bool":
		if currentElementKind != common.ValueKindBoolean {
			p.reportMismatch(fmt.Sprintf("expected bool got %s", currentElementKind))
			return
		}
	case "float":
		if currentElementKind != common.ValueKindFloat {
			p.reportMismatch(fmt.Sprintf("expected float got %s", currentElementKind))
			return
		}
	}

	p.statement.Parameters = append(p.statement.Parameters, logiAst.Parameter{
		Name:  syntaxStatementElement.VariableKeyword.Name,
		Value: value,
	})
}

func (p *recursiveStatementParser) matchTypeStatement(name string, statement macroAst.TypeStatement) {
	for _, syntaxElement := range statement.Elements {
		if p.pei >= len(p.plainStatement.Elements) {
			p.reportMismatch("statement is shorter than syntax")
			return
		}

		var currentElement = p.plainStatement.Elements[p.pei]

		p.matchNextElement(syntaxElement, currentElement)

		if p.mismatchCause != "" {
			return // stop matching if a mismatch is found
		}

		p.pei++
	}

	p.pei-- // decrement the index to match the last element
}

func (p *recursiveStatementParser) matchScope(plainElement plain.DefinitionStatementElement, syntaxStatementElement macroAst.SyntaxStatementElement) {
	if plainElement.Kind != plain.DefinitionStatementElementKindStruct {
		p.reportMismatch(fmt.Sprintf("expected structure, got %s", plainElement.Kind))
		return
	}

	var scopeMap = make(map[string]macroAst.ScopeItem)

	for _, scope := range syntaxStatementElement.ScopeDef.Scopes {
		var scopeFound = false
		for _, item := range p.macroDefinition.Scopes.Scopes {
			if item.Name == scope {
				scopeMap[scope] = item
				scopeFound = true
			}
		}

		if !scopeFound {
			p.reportMismatch(fmt.Sprintf("scope %s not found", scope))
			return
		}
	}

	var result = make([]logiAst.Statement, 0)

MainLoop:
	for _, item := range plainElement.Struct.Statements {
		var maxMatch = -1
		var mismatchCause string
		var bestMatch *macroAst.SyntaxStatement

		for _, scopName := range syntaxStatementElement.ScopeDef.Scopes {
			scope := scopeMap[scopName]

			sp := recursiveStatementParser{
				macroDefinition: &macroAst.Macro{
					Types: p.macroDefinition.Types,
					Syntax: macroAst.Syntax{
						Statements: scope.Statements,
					},
					Scopes: p.macroDefinition.Scopes,
				},
			}

			sp.plainStatement = item

			err := sp.parse(scope.Name)

			if err != nil {
				if sp.maxMatch > maxMatch {
					maxMatch = sp.maxMatch
					bestMatch = sp.bestMatch
					mismatchCause = err.Error()
				}
				continue
			}

			result = append(result, sp.statement)

			continue MainLoop
		}

		p.maxMatch = maxMatch
		p.bestMatch = bestMatch
		p.reportMismatch(mismatchCause)
		return
	}

	p.statement.SubStatements = append(p.statement.SubStatements, result)
}

func (p *recursiveStatementParser) matchArray(plainElement plain.DefinitionStatementElement, syntaxElement macroAst.SyntaxStatementElement) {
	if syntaxElement.VariableKeyword.Type.Name != "array" {
		p.reportMismatch(fmt.Sprintf("expected %s got array", syntaxElement.VariableKeyword.Type.Name))
		return
	}
	itemSyntaxStatement := macroAst.SyntaxStatement{
		Elements: []macroAst.SyntaxStatementElement{
			{
				Kind: macroAst.SyntaxStatementElementKindVariableKeyword,

				VariableKeyword: &macroAst.SyntaxStatementElementVariableKeyword{
					Name: syntaxElement.VariableKeyword.Name,
					Type: syntaxElement.VariableKeyword.Type.SubTypes[0],
				},
			},
		},
	}

	var result = make([]logiAst.Statement, 0)
	for _, item := range plainElement.Array.Items {
		sp := recursiveStatementParser{
			plainStatement:  item,
			syntaxStatement: itemSyntaxStatement,
			macroDefinition: p.macroDefinition,
		}
		sp.matchValue(itemSyntaxStatement.Elements[0])

		if sp.mismatchCause != "" {
			p.reportMismatch(sp.mismatchCause)
			break
		}

		result = append(result, sp.statement)

		if sp.pei+1 < len(sp.plainStatement.Elements) {
			p.reportMismatch("type has more elements than syntax")
			break
		}
	}

	p.statement.SubStatements = append(p.statement.SubStatements, result)
	return
}

func (p *recursiveStatementParser) matchCombination(element macroAst.SyntaxStatementElement) {
	var currentElement = p.plainStatement.Elements[p.pei]

	var recordedPei = p.pei
	var recordedStatement = p.statement
	var cause string
	var maxMatch = -1

	if p.mismatchCause != "" {
		panic("mismatch cause is not empty")
	}

	for _, syntaxElement := range element.Combination.Elements {
		p.pei = recordedPei
		p.statement = recordedStatement
		p.matchNextElement(syntaxElement, currentElement)

		if p.mismatchCause == "" { // If it matches return
			return // stop matching, we found a match
		} else {
			if p.pei > maxMatch {
				maxMatch = p.pei
				cause = p.mismatchCause
			}
			p.mismatchCause = "" // reset the mismatch cause
		}
	}

	p.reportMismatch(fmt.Sprintf("no combination matched for: %v at %s; %s", currentElement.AsValue().AsInterface(), currentElement.SourceLocation, cause))
}

func (p *recursiveStatementParser) matchTypeReference(syntaxElement macroAst.SyntaxStatementElement, currentElement plain.DefinitionStatementElement) {
	// locate type

	for _, typeStatement := range p.macroDefinition.Types.Types {
		if typeStatement.Name == syntaxElement.TypeReference.Name {

			for _, syntaxElement := range typeStatement.Elements {
				if p.pei >= len(p.plainStatement.Elements) {
					p.reportMismatch("statement is shorter than syntax")
					return
				}

				var currentElement = p.plainStatement.Elements[p.pei]

				p.matchNextElement(syntaxElement, currentElement)

				if p.mismatchCause != "" {
					return // stop matching if a mismatch is found
				}

				p.pei++
			}

			return
		}
	}

	p.pei--
}

func isSyntaxElementAlwaysRequired(element macroAst.SyntaxStatementElement) bool {
	switch element.Kind {
	case macroAst.SyntaxStatementElementKindAttributeList:
		return false
	}

	return true
}
