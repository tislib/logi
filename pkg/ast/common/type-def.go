package common

type TypeDefinition struct {
	Name     string           `json:"typeName,omitempty"`
	SubTypes []TypeDefinition `json:"subTypes,omitempty"`
}

func (d TypeDefinition) ToDisplayName() string {
	if len(d.SubTypes) == 0 {
		return d.Name
	}

	var result = d.Name + "<"
	for i, subType := range d.SubTypes {
		if i > 0 {
			result += ","
		}
		result += subType.ToDisplayName()
	}

	return result + ">"
}
