package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
	"sort"
)

type recursiveStatementParser struct {
	plainStatement  plain.DefinitionStatement
	syntaxStatement macroAst.SyntaxStatement
	macroDefinition *macroAst.Macro

	maxMatch      int
	bestMatch     *macroAst.SyntaxStatement
	mismatchCause string

	sei        int
	pei        int
	asr        analyseStatementResult
	definition *logiAst.Definition
}

type analyseStatementResult struct {
	hasName          bool
	stopMatchingName bool
	hasType          bool
	hasCodeBlock     bool
	hasArgumentList  bool
	nameParts        []string
	parameters       map[string]common.Value
	typeDef          common.TypeDefinition
	hasAttributeList bool
	attributes       []plain.DefinitionStatementElementAttribute
	arguments        []plain.DefinitionStatementElementArgument
	codeBlock        common.CodeBlock
}

func (p *recursiveStatementParser) parse() error {
	for _, syntaxStatement := range p.macroDefinition.Syntax.Statements {
		p.asr = analyseStatementResult{
			parameters: make(map[string]common.Value),
		}
		p.sei = 0
		p.pei = 0
		p.mismatchCause = ""
		p.syntaxStatement = syntaxStatement

		p.match()

		if p.mismatchCause != "" {
			continue
		}

		return p.apply()
	}

	return fmt.Errorf("failed to match statement: %s", p.mismatchCause)
}

func (p *recursiveStatementParser) apply() error {
	var asr = p.asr

	var name = camelCaseFromNameParts(asr.nameParts)
	var parameters []logiAst.Parameter
	var attributes []logiAst.Attribute
	var arguments []logiAst.Argument

	var data = make(map[string]interface{})

	var keys []string
	for key := range asr.parameters {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		parameters = append(parameters, logiAst.Parameter{
			Name:  key,
			Value: asr.parameters[key],
		})

		data[key] = asr.parameters[key].AsInterface()
	}

	for _, attribute := range asr.attributes {
		attributes = append(attributes, logiAst.Attribute{
			Name:  attribute.Name,
			Value: attribute.Value,
		})

		if attribute.Value != nil {
			data[attribute.Name] = attribute.Value.AsInterface()
		} else {
			data[attribute.Name] = true
		}
	}

	for _, argument := range asr.arguments {
		arguments = append(arguments, logiAst.Argument{
			Name: argument.Name,
			Type: argument.Type,
		})

		data[argument.Name] = map[string]interface{}{
			"name": argument.Name,
			"type": argument.Type.AsValue().AsInterface(),
		}
	}

	if asr.hasCodeBlock {
		data["code"] = asr.codeBlock
	}

	if asr.hasType {
		data["type"] = asr.typeDef.AsValue().AsInterface()
	}

	if asr.hasName {
		if asr.hasArgumentList {
			if asr.hasCodeBlock { // method
				p.definition.Methods = append(p.definition.Methods, logiAst.Method{
					Name:       name,
					Type:       asr.typeDef,
					Attributes: attributes,
					Parameters: parameters,
					Arguments:  arguments,
					CodeBlock:  asr.codeBlock,
				})
			} else { // method signature
				p.definition.MethodSignature = append(p.definition.MethodSignature, logiAst.MethodSignature{
					Name:       name,
					Type:       asr.typeDef,
					Attributes: attributes,
					Parameters: parameters,
					Arguments:  arguments,
				})
			}
		} else {
			if asr.hasType { // property
				p.definition.Properties = append(p.definition.Properties, logiAst.Property{
					Name:       name,
					Type:       asr.typeDef,
					Attributes: attributes,
					Parameters: parameters,
				})
			} else { // parameter
				var parameter = logiAst.DefinitionParameter{
					Name:       name,
					Attributes: attributes,
					Parameters: parameters,
				}
				if asr.hasCodeBlock {
					parameter.CodeBlock = &asr.codeBlock
				}
				p.definition.Parameters = append(p.definition.Parameters, parameter)
			}
		}
	}

	p.definition.Dynamic[name] = data

	return nil
}

func (p *recursiveStatementParser) reportMismatch(reason string) {
	if p.pei >= p.maxMatch {
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
	switch syntaxStatementElement.Kind {
	case macroAst.SyntaxStatementElementKindKeyword:
		if currentElement.Kind != plain.DefinitionStatementElementKindIdentifier {
			p.reportMismatch(fmt.Sprintf("expected keyword (%s), got %s (%v)", syntaxStatementElement.KeywordDef.Name, currentElement.Kind, currentElement))
			break
		} else if currentElement.Identifier.Identifier != syntaxStatementElement.KeywordDef.Name {
			p.reportMismatch(fmt.Sprintf("expected keyword (%s), got %s", syntaxStatementElement.KeywordDef.Name, currentElement.Identifier.Identifier))
			break
		}
		if !p.asr.hasName || !p.asr.stopMatchingName {
			p.asr.hasName = true
			p.asr.nameParts = append(p.asr.nameParts, currentElement.Identifier.Identifier)
		}
	case macroAst.SyntaxStatementElementKindTypeReference:
		p.matchTypeReference(syntaxStatementElement, currentElement)
	case macroAst.SyntaxStatementElementKindVariableKeyword:
		switch currentElement.Kind {
		case plain.DefinitionStatementElementKindIdentifier:
			p.asr.parameters[syntaxStatementElement.VariableKeyword.Name] = common.StringValue(currentElement.Identifier.Identifier)
			if syntaxStatementElement.VariableKeyword.Type.Name == "Type" {
				p.asr.hasType = true
				p.asr.stopMatchingName = true
				p.asr.typeDef = common.TypeDefinition{
					Name: currentElement.Identifier.Identifier,
				}
				return
			} else if syntaxStatementElement.VariableKeyword.Type.Name != "Name" {
				p.reportMismatch(fmt.Sprintf("expected variable keyword %s in type %s, got %s in type Name", syntaxStatementElement.VariableKeyword.Name, syntaxStatementElement.VariableKeyword.Type, currentElement.Identifier.Identifier))
				break
			} else {
				p.matchIdentifierType(syntaxStatementElement, currentElement)
			}
			if syntaxStatementElement.VariableKeyword.Type.Name == "Name" {
				if !p.asr.hasName || !p.asr.stopMatchingName {
					p.asr.hasName = true
					p.asr.nameParts = append(p.asr.nameParts, currentElement.Identifier.Identifier)
				}
			} else {
				p.asr.stopMatchingName = true
			}
		case plain.DefinitionStatementElementKindArray:
			p.asr.stopMatchingName = true
			p.matchArray(currentElement, syntaxStatementElement)
			return
		case plain.DefinitionStatementElementKindValue:
			p.asr.stopMatchingName = true
			p.matchValue(syntaxStatementElement)
		default:
			p.asr.stopMatchingName = true
			p.reportMismatch(fmt.Sprintf("expected variable keyword (%s), got %s", syntaxStatementElement.VariableKeyword.Name, currentElement.Kind))
		}
	case macroAst.SyntaxStatementElementKindAttributeList:
		p.asr.stopMatchingName = true
		if currentElement.Kind != plain.DefinitionStatementElementKindAttributeList {
			p.reportMismatch(fmt.Sprintf("expected attribute list, got %s", currentElement.Kind))

			break
		}

		// check if the attribute list is valid
		err := isValidAttributeList(currentElement.AttributeList, syntaxStatementElement.AttributeList)

		if err != nil {
			p.reportMismatch(fmt.Sprintf(err.Error()))
			break
		}

		p.asr.hasAttributeList = true
		p.asr.attributes = currentElement.AttributeList.Attributes
	case macroAst.SyntaxStatementElementKindArgumentList:
		p.asr.stopMatchingName = true
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
		p.asr.hasArgumentList = true
		p.asr.arguments = currentElement.ArgumentList.Arguments
	case macroAst.SyntaxStatementElementKindCodeBlock:
		p.asr.stopMatchingName = true
		if currentElement.Kind != plain.DefinitionStatementElementKindCodeBlock {
			p.reportMismatch(fmt.Sprintf("expected code block, got %s", currentElement.Kind))

			break
		}
		p.asr.hasCodeBlock = true
		p.asr.codeBlock = currentElement.CodeBlock.CodeBlock
	case macroAst.SyntaxStatementElementKindParameterList:
		p.asr.stopMatchingName = true
		if currentElement.Kind != plain.DefinitionStatementElementKindParameterList {
			p.reportMismatch(fmt.Sprintf("expected parameter list, got %s", currentElement.Kind))

			break
		}

		// check if the parameter list is valid
		for idx, syntaxStatementElementParameter := range syntaxStatementElement.ParameterList.Parameters {
			p.asr.parameters[syntaxStatementElementParameter.Name] = currentElement.ParameterList.Parameters[idx].Value
		}
	case macroAst.SyntaxStatementElementKindCombination:
		p.asr.stopMatchingName = true
		p.matchCombination(syntaxStatementElement)
	case macroAst.SyntaxStatementElementKindStructure:
		p.asr.stopMatchingName = true
		p.matchStructure(syntaxStatementElement, currentElement)
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
	currentElementKind := currentElement.Value.Value.Kind

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

	p.asr.parameters[syntaxStatementElement.VariableKeyword.Name] = currentElement.Value.Value
}

func (p *recursiveStatementParser) matchTypeStatement(name string, statement macroAst.TypeStatement) {
	var asrBck = p.asr
	p.asr = analyseStatementResult{
		parameters: make(map[string]common.Value),
	}
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

	asrCopy := p.asr
	p.asr = asrBck

	p.asr.parameters[name] = common.Value{
		Kind: common.ValueKindMap,
		Map:  asrCopy.parameters,
	}

	p.pei-- // decrement the index to match the last element
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

	var valueArr []common.Value
	for _, item := range plainElement.Array.Items {
		sp := recursiveStatementParser{
			plainStatement:  item,
			syntaxStatement: itemSyntaxStatement,
			macroDefinition: p.macroDefinition,
			asr: analyseStatementResult{
				parameters: make(map[string]common.Value),
			},
		}
		sp.matchValue(itemSyntaxStatement.Elements[0])

		if sp.mismatchCause != "" {
			p.reportMismatch(sp.mismatchCause)
			break
		}

		valueArr = append(valueArr, sp.asr.parameters[syntaxElement.VariableKeyword.Name])

		if sp.pei+1 < len(sp.plainStatement.Elements) {
			p.reportMismatch("type has more elements than syntax")
			break
		}
	}

	p.asr.parameters[syntaxElement.VariableKeyword.Name] = common.Value{
		Kind:  common.ValueKindArray,
		Array: valueArr,
	}
	return
}

func (p *recursiveStatementParser) matchCombination(element macroAst.SyntaxStatementElement) {
	var currentElement = p.plainStatement.Elements[p.pei]

	var recordedPei = p.pei
	var recordedAsr = p.asr

	if p.mismatchCause != "" {
		panic("mismatch cause is not empty")
	}

	for _, syntaxElement := range element.Combination.Elements {
		p.pei = recordedPei
		p.asr = recordedAsr
		p.matchNextElement(syntaxElement, currentElement)

		if p.mismatchCause == "" { // If it matches return
			return // stop matching, we found a match
		} else {
			p.mismatchCause = "" // reset the mismatch cause
		}
	}
}

func (p *recursiveStatementParser) matchStructure(rootSyntaxElement macroAst.SyntaxStatementElement, plainElement plain.DefinitionStatementElement) {
	if plainElement.Kind != plain.DefinitionStatementElementKindStruct {
		p.reportMismatch(fmt.Sprintf("expected structure, got %s", plainElement.Kind))
		return
	}

	sp := recursiveStatementParser{
		macroDefinition: &macroAst.Macro{
			Types: p.macroDefinition.Types,
			Syntax: macroAst.Syntax{
				Statements: rootSyntaxElement.Structure.Statements,
			},
		},
		definition: &logiAst.Definition{
			PlainStatements: plainElement.Struct.Statements,
			Dynamic:         make(map[string]map[string]interface{}),
		},
		asr: analyseStatementResult{
			parameters: make(map[string]common.Value),
		},
	}

	var parameters map[string]common.Value
	for _, item := range plainElement.Struct.Statements {
		sp.plainStatement = item

		err := sp.parse()

		if err != nil {
			p.reportMismatch(err.Error())
			return // stop matching if a mismatch is found
		}

		err = sp.apply()

		if err != nil {
			p.reportMismatch(err.Error())
			return // stop matching if a mismatch is found
		}

		if len(sp.asr.parameters) > 0 {
			if parameters == nil {
				parameters = make(map[string]common.Value)
			}

			for _, param := range sp.definition.Parameters {
				var result = map[string]common.Value{}
				for _, subParam := range param.Parameters {
					result[subParam.Name] = subParam.Value
				}
				parameters[param.Name] = common.Value{
					Kind: common.ValueKindMap,
					Map:  result,
				}
			}
		}
	}

	// migrate values
	if len(parameters) > 0 {
		var name = camelCaseFromNameParts(p.asr.nameParts)
		p.asr.parameters[name] = common.Value{
			Kind: common.ValueKindMap,
			Map:  parameters,
		}
	}

	return
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

func (p *recursiveStatementParser) matchIdentifierType(element macroAst.SyntaxStatementElement, element2 plain.DefinitionStatementElement) {

}

func isSyntaxElementAlwaysRequired(element macroAst.SyntaxStatementElement) bool {
	switch element.Kind {
	case macroAst.SyntaxStatementElementKindAttributeList:
		return false
	}

	return true
}
