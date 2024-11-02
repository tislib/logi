package logi

type NodeOp string

const (
	NodeOpFile          = "file"
	NodeOpSignature     = "signature"
	NodeOpMacro         = "macro"
	NodeOpName          = "name"
	NodeOpValue         = "value"
	NodeOpBody          = "body"
	NodeOpDefinition    = "definition"
	NodeOpStatements    = "statements"
	NodeOpStatement     = "statement"
	NodeOpIdentifier    = "identifier"
	NodeOpAttributeList = "attributeList"
	NodeOpAttribute     = "attribute"
)

type yaccNode struct {
	op       NodeOp
	children []yaccNode
	value    interface{}
}

func appendNode(nodeOp NodeOp, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, children: children}
}

func newNode(nodeOp NodeOp, value interface{}, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, value: value, children: children}
}

func appendNodeTo(node *yaccNode, child yaccNode) yaccNode {
	node.children = append(node.children, child)

	return *node
}

func registerRootNode(parser yyLexer, n yaccNode) {
	parser.(*yyLogiLexerProxy).Node.children = append(parser.(*yyLogiLexerProxy).Node.children, n)
}

func assertEqual(parser yyLexer, a, b interface{}, msg string) {
	if a != b {
		parser.Error(msg)
	}
}

func init() {
	yyErrorVerbose = true
	yyDebug = 5
}
