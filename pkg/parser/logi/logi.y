%{
package logi

import (
)

%}

%union {
	node yaccNode
	bool bool
	number interface{}
	string string
}

%token<number> token_number
%token<string> token_string
%token<string> token_identifier
%token<bool> token_bool

// Keywords
%token DefinitionKeyword SyntaxKeyword LogiKeyword

%token BracketOpen BracketClose BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dash Dot Arrow ParenOpen ParenClose Eol

%type<node> definition definition_signature definition_body definition_statements definition_statement definition_statement_element
%type<node> definition_statement_element_identifier definition_statement_element_value definition_statement_element_attribute_list definition_statement_element_attribute_list_content definition_statement_element_attribute_list_item
%start file

%%

// Helpers
eol_allowed: Eol
| eol_allowed Eol
| // empty
;

eol_required: Eol
| eol_required Eol
;

file: definition eol_allowed {
	registerRootNode(yylex, $1)
}
| file eol_allowed
| eol_allowed
| file definition eol_allowed {
	registerRootNode(yylex, $2)
};

definition: definition_signature eol_allowed definition_body
{
	$$ = appendNode(NodeOpDefinition, $1, $3)
};

definition_signature: token_identifier token_identifier
{
	$$ = appendNode(NodeOpSignature, newNode(NodeOpMacro, $1), newNode(NodeOpName, $2))
};

definition_body: BraceOpen eol_allowed
definition_statements eol_allowed
BraceClose
{
	$$ = appendNode(NodeOpBody, $3)
};

definition_statements: definition_statement eol_required
{
	$$ = appendNode(NodeOpStatements, $1)
}
| definition_statements definition_statement eol_required
{
	$$ = appendNodeTo(&$1, $2)
};

definition_statement: definition_statement_element
{
	$$ = appendNode(NodeOpStatement, $1)
}
| definition_statement definition_statement_element
{
	$$ = appendNodeTo(&$1, $2)
};

definition_statement_element: definition_statement_element_identifier | definition_statement_element_value | definition_statement_element_attribute_list;

definition_statement_element_identifier: token_identifier
{
	$$ = newNode(NodeOpIdentifier, $1)
};

definition_statement_element_value: token_string
{
	$$ = newNode(NodeOpValue, $1)
}
| token_number
{
	$$ = newNode(NodeOpValue, $1)
}
| token_bool
{
	$$ = newNode(NodeOpValue, $1)
};

definition_statement_element_attribute_list: BracketOpen definition_statement_element_attribute_list_content BracketClose
{
	$$ = $2
};

definition_statement_element_attribute_list_content: definition_statement_element_attribute_list_item
{
	$$ = appendNode(NodeOpAttributeList, $1)
}
| definition_statement_element_attribute_list_content Comma definition_statement_element_attribute_list_item
{
	$$ = appendNodeTo(&$1, $3)
};

definition_statement_element_attribute_list_item: token_identifier
{
	$$ = newNode(NodeOpAttribute, $1)
}
| token_identifier definition_statement_element_value
{
	$$ = newNode(NodeOpAttribute, $1, $2)
};

%%

// 	entity User {
//   		id int [primary, autoincrement]
//  		name string [required, default "John Doe"]
//   	}