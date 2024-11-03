package logi

import (
	"logi/pkg/ast/common"
	"logi/pkg/ast/plain"
)

type Ast struct {
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	MacroName       string                      `json:"macroName"`
	Name            string                      `json:"name"`
	PlainStatements []plain.DefinitionStatement `json:"plainStatements"`
	Properties      []Property                  `json:"properties"`
	MethodSignature []MethodSignature           `json:"methodSignature"`
	Methods         []Method                    `json:"methods"`
}

type Property struct {
	Name       string                `json:"name"`
	Type       common.TypeDefinition `json:"type"`
	Attributes []Attribute           `json:"attributes"`
	Parameters []Parameter           `json:"parameters"`
	CodeBlock  *common.CodeBlock     `json:"codeBlock"`
}

type MethodSignature struct {
	Name       string                `json:"name"`
	Type       common.TypeDefinition `json:"type"`
	Attributes []Attribute           `json:"attributes"`
	Parameters []Parameter           `json:"parameters"`
	Arguments  []Argument            `json:"arguments"`
}

type Method struct {
	Name       string                `json:"name"`
	Type       common.TypeDefinition `json:"type"`
	Attributes []Attribute           `json:"attributes"`
	Parameters []Parameter           `json:"parameters"`
	Arguments  []Argument            `json:"arguments"`
	CodeBlock  *common.CodeBlock     `json:"codeBlock"`
}

type Attribute struct {
	Name  string        `json:"name"`
	Value *common.Value `json:"value"`
}

type Parameter struct {
	Name  string        `json:"name"`
	Value *common.Value `json:"value"`
}

type Argument struct {
	Name string                `json:"name"`
	Type common.TypeDefinition `json:"type"`
}
