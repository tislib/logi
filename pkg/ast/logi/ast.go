package logi

import (
	"github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
)

type Ast struct {
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	MacroName       string                      `json:"macroName"`
	Name            string                      `json:"name"`
	PlainStatements []plain.DefinitionStatement `json:"plainStatements"`
	Statements      []Statement                 `json:"statements"`
}

type Statement struct {
	Command       string      `json:"command"`
	Arguments     []Argument  `json:"arguments"`
	Parameters    []Parameter `json:"parameters"`
	Attributes    []Attribute `json:"attributes"`
	SubStatements []Statement `json:"subStatements"`
}

type Attribute struct {
	Name  string        `json:"name"`
	Value *common.Value `json:"value"`
}

type Parameter struct {
	Name       string             `json:"name"`
	Value      common.Value       `json:"value"`
	Expression *common.Expression `json:"expression"`
}

type Argument struct {
	Name string                `json:"name"`
	Type common.TypeDefinition `json:"type"`
}
