package golang

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	commonAst "github.com/tislib/logi/pkg/ast/common"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/generate/common"
	"go/format"
	"strings"
)

type generator struct {
}

func (g generator) Generate(ast *macroAst.Ast, pkg string) (string, error) {
	var code string

	code += "package " + pkg + "\n\n"

	for _, macro := range ast.Macros {
		err := g.writeMacro(&code, macro)

		if err != nil {
			return "", fmt.Errorf("error writing macro: %v", err)
		}
	}

	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		formattedCode = []byte(code)

		log.Warn(err)
	}

	return string(formattedCode), nil
}

func (g generator) writeMacro(code *string, macro macroAst.Macro) error {
	*code += fmt.Sprintf("type %s struct {\n", g.getStructName(macro))
	for _, statement := range macro.Syntax.Statements {
		err := g.writeStatementAsProperty(code, macro, statement)

		if err != nil {
			return fmt.Errorf("error writing statement: %v", err)
		}
	}
	*code += "}\n"

	return nil
}

func (g generator) writeStatementAsProperty(code *string, macro macroAst.Macro, statement macroAst.SyntaxStatement) error {
	analyser := common.NewAnalyser()

	asr, err := analyser.AnalyseStatement(macro, statement)

	if err != nil {
		return fmt.Errorf("error analysing statement: %v", err)
	}

	if !asr.HasName {
		return nil
	}

	var name = camelCaseFromNameParts(asr.NameParts...)

	var goType = g.getGoTypeForAsr(asr)

	*code += fmt.Sprintf("  %s %s\n", pascalCaseFromNameParts(name), goType)

	return nil
}

func (g generator) getGoTypeForAsr(asr common.AnalyseStatementResult) string {
	if asr.HasType {
		if asr.TypeDef.Name == "array" {
			if asr.SubTypeAsr != nil {
				return "[]" + g.getGoTypeForAsr(*asr.SubTypeAsr)
			} else {
				return "[]" + g.getGoType(asr.TypeDef.SubTypes[0])
			}
		} else {
			if asr.SubTypeAsr != nil {
				return g.getGoTypeForAsr(*asr.SubTypeAsr)
			}
			return g.getGoType(asr.TypeDef)
		}
	} else if asr.HasCodeBlock {
		return "func(vm *vm.VM, args ...interface{}) (interface{}, error)"
	} else if len(asr.Parameters) > 0 {
		var code = "struct {\n"
		for name, def := range asr.Parameters {
			code += fmt.Sprintf("  %s %s\n", pascalCaseFromNameParts(name), g.getGoType(def))
		}
		code += "}"
		return code
	}

	return "interface{}"
}

func (g generator) getStructName(macro macroAst.Macro) string {
	return ToTitle(macro.Name)
}

func (g generator) getGoType(def commonAst.TypeDefinition) string {
	switch def.Name {
	case "Name":
		return "string"
	case "number":
		return "int"
	}

	return def.Name
}

func ToTitle(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", string(s[0]+'A'-'a'), s[1:])
}

func NewGenerator() common.CodeGenerator {
	return &generator{}
}

func camelCaseFromNameParts(parts ...string) string {
	var result string

	for i, part := range parts {
		if i == 0 {
			result += strings.ToLower(part[:1]) + part[1:]
		} else {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return result
}

func pascalCaseFromNameParts(parts ...string) string {
	var result string

	for _, part := range parts {
		result += strings.ToUpper(part[:1]) + part[1:]
	}

	return result
}
