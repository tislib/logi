package logi

type ParseNodeVisitFunc func(node *ParseNode) error

type ParseNodeMode int

const (
	ModeToken ParseNodeMode = iota
	ModeSequence
	ModeOr
)

type ParseNode struct {
	TokenId    int
	TokenValue string
	Mode       ParseNodeMode

	Children []*ParseNode

	Value interface{}

	VisitFunc ParseNodeVisitFunc
}
