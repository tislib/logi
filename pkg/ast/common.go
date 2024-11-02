package ast

type TypeDefinition struct {
	Name     string           `json:"typeName,omitempty"`
	SubTypes []TypeDefinition `json:"subTypes,omitempty"`
}
