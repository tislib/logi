package common

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type AnalyseStatementResult struct {
	HasName          bool
	HasType          bool
	HasTypeParameter bool
	HasCodeBlock     bool
	HasArgumentList  bool
	NameParts        []string
	Parameters       map[string]common.TypeDefinition
	TypeDef          common.TypeDefinition
	HasAttributeList bool
	Attributes       map[string]common.TypeDefinition
	Arguments        map[string]common.TypeDefinition
	SubTypeAsr       *AnalyseStatementResult
}

type analyser struct {
}

func (a analyser) AnalyseStatement(macro macroAst.Macro, statement macroAst.SyntaxStatement) (AnalyseStatementResult, error) {
	var asr = AnalyseStatementResult{
		Parameters: make(map[string]common.TypeDefinition),
		Attributes: make(map[string]common.TypeDefinition),
		Arguments:  make(map[string]common.TypeDefinition),
	}

	err := a.analyseStatement(macro, statement, &asr)

	if err != nil {
		return asr, err
	}

	return asr, nil
}

func (a analyser) analyseStatement(macro macroAst.Macro, statement macroAst.SyntaxStatement, asr *AnalyseStatementResult) error {
	for _, element := range statement.Elements {
		switch element.Kind {
		case macroAst.SyntaxStatementElementKindKeyword:
			asr.HasName = true
			asr.NameParts = append(asr.NameParts, element.KeywordDef.Name)
		case macroAst.SyntaxStatementElementKindTypeReference:
			typeDef := a.locateTypeDefinition(macro, element.TypeReference.Name)

			if typeDef == nil {
				return fmt.Errorf("type definition not found: %s", element.TypeReference.Name)
			}

			err := a.analyseStatement(macro, macroAst.SyntaxStatement{Elements: typeDef.Elements}, asr)

			if err != nil {
				return fmt.Errorf("error analysing type definition: %v", err)
			}
		case macroAst.SyntaxStatementElementKindVariableKeyword:
			asr.Parameters[element.VariableKeyword.Name] = element.VariableKeyword.Type
			switch element.VariableKeyword.Type.Name {
			case "Type":
				asr.HasTypeParameter = true
			case "Name":
				asr.HasName = true
				asr.NameParts = append(asr.NameParts, element.VariableKeyword.Name)
			case "array":
				asr.HasType = true
				asr.TypeDef = element.VariableKeyword.Type

				typeDef := a.locateTypeDefinition(macro, element.VariableKeyword.Type.SubTypes[0].Name)

				if typeDef != nil {
					var subTypeAsr = AnalyseStatementResult{
						Parameters: make(map[string]common.TypeDefinition),
						Attributes: make(map[string]common.TypeDefinition),
						Arguments:  make(map[string]common.TypeDefinition),
					}
					err := a.analyseStatement(macro, macroAst.SyntaxStatement{Elements: typeDef.Elements}, &subTypeAsr)

					if err != nil {
						return fmt.Errorf("error analysing type definition: %v", err)
					}

					asr.SubTypeAsr = &subTypeAsr
				}
			default:
				asr.HasType = true
				asr.TypeDef = element.VariableKeyword.Type

				typeDef := a.locateTypeDefinition(macro, element.VariableKeyword.Type.Name)

				if typeDef != nil {
					var subTypeAsr = AnalyseStatementResult{
						Parameters: make(map[string]common.TypeDefinition),
						Attributes: make(map[string]common.TypeDefinition),
						Arguments:  make(map[string]common.TypeDefinition),
					}
					err := a.analyseStatement(macro, macroAst.SyntaxStatement{Elements: typeDef.Elements}, &subTypeAsr)

					if err != nil {
						return fmt.Errorf("error analysing type definition: %v", err)
					}

					asr.SubTypeAsr = &subTypeAsr
				}
			}
		case macroAst.SyntaxStatementElementKindAttributeList:
			asr.HasAttributeList = true
			for _, attribute := range element.AttributeList.Attributes {
				asr.Attributes[attribute.Name] = attribute.Type
			}
		case macroAst.SyntaxStatementElementKindCombination:
			continue
		case macroAst.SyntaxStatementElementKindCodeBlock:
			asr.HasCodeBlock = true
		case macroAst.SyntaxStatementElementKindArgumentList:
			asr.HasArgumentList = true
			for _, argument := range element.ArgumentList.Arguments {
				asr.Arguments[argument.Name] = argument.Type
			}
		case macroAst.SyntaxStatementElementKindStructure:
			for _, statement := range element.Structure.Statements {
				err := a.analyseStatement(macro, statement, asr)

				if err != nil {
					return fmt.Errorf("error analysing structure statement: %v", err)
				}
			}
		case macroAst.SyntaxStatementElementKindParameterList:
			for _, parameter := range element.ParameterList.Parameters {
				asr.Parameters[parameter.Name] = parameter.Type
			}
		default:
			return fmt.Errorf("unknown element kind: %s", element.Kind)
		}
	}

	return nil
}

func (a analyser) locateTypeDefinition(macro macroAst.Macro, name string) *macroAst.TypeStatement {
	for _, typeDef := range macro.Types.Types {
		if typeDef.Name == name {
			return &typeDef
		}
	}

	return nil
}

type Analyser interface {
	AnalyseStatement(macro macroAst.Macro, statement macroAst.SyntaxStatement) (AnalyseStatementResult, error)
}

func NewAnalyser() Analyser {
	return &analyser{}
}
