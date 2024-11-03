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

%token BracketOpen BracketClose BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dot Arrow ParenOpen ParenClose Eol
%token IfKeyword ElseKeyword ReturnKeyword SwitchKeyword CaseKeyword VarKeyword
%token Plus Minus Star Slash Percent Exclamation And Or Xor

%type<node> type_definition
%type<node> definition definition_signature definition_body definition_statements definition_statement definition_statement_element
%type<node> definition_statement_element_identifier definition_statement_element_value definition_statement_element_attribute_list definition_statement_element_attribute_list_content definition_statement_element_attribute_list_item
%type<node> definition_statement_element_argument_list definition_statement_element_argument_list_content definition_statement_element_argument_list_item
%type<node> code_block code_block_statements code_block_statement expression_statement assignment_statement if_statement return_statement variable_declaration_statement
%type<node> expression literal variable binary_expression function_call
%type<node> function_params
%type<string> operator
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

definition_statement_element: definition_statement_element_identifier | definition_statement_element_value | definition_statement_element_attribute_list | definition_statement_element_argument_list | code_block;

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

definition_statement_element_argument_list: ParenOpen definition_statement_element_argument_list_content ParenClose
{
	$$ = $2
};

definition_statement_element_argument_list_content: definition_statement_element_argument_list_item
{
	$$ = appendNode(NodeOpArgumentList, $1)
}
| definition_statement_element_argument_list_content Comma definition_statement_element_argument_list_item
{
	$$ = appendNodeTo(&$1, $3)
};

definition_statement_element_argument_list_item: token_identifier type_definition
{
	$$ = newNode(NodeOpArgument, $1, $2)
};


type_definition: token_identifier
{
	$$ = newNode(NodeOpTypeDef, $1)
}
| token_identifier LessThan type_definition GreaterThan
{
	$$ = newNode(NodeOpTypeDef, $1, $3)
};


// ################################ CODE_BLOCK ############################

code_block: BraceOpen eol_allowed code_block_statements eol_allowed BraceClose
{
	$$ = appendNode(NodeOpCodeBlock, $3)
};

code_block_statements: code_block_statement eol_required
{
	$$ = appendNode(NodeOpStatements, $1)
}
| code_block_statements code_block_statement eol_required
{
	$$ = appendNodeTo(&$1, $2)
}
| // empty
{
	$$ = appendNode(NodeOpStatements)
}
;

code_block_statement: expression_statement | assignment_statement | if_statement | return_statement | variable_declaration_statement;

// Statements

expression_statement: expression
{
	$$ = appendNode(NodeOpExpression, $1)
};

assignment_statement: expression Equal expression
{
	$$ = appendNode(NodeOpAssignment, $1, $3)
};

if_statement: IfKeyword ParenOpen expression ParenClose code_block
{
	$$ = appendNode(NodeOpIf, $3, $5)
}
| IfKeyword ParenOpen expression ParenClose code_block ElseKeyword code_block
{
	$$ = appendNode(NodeOpIfElse, $3, $5, $7)
}
| IfKeyword expression code_block
{
	$$ = appendNode(NodeOpIf, $2, $3)
}
| IfKeyword expression code_block ElseKeyword code_block
{
	$$ = appendNode(NodeOpIfElse, $2, $3, $5)
};

return_statement: ReturnKeyword expression
{
	$$ = appendNode(NodeOpReturn, $2)
};

variable_declaration_statement: VarKeyword token_identifier type_definition
{
	$$ = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, $2), $3)
}
| VarKeyword token_identifier type_definition Equal expression
{
	$$ = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, $2), $3, $5)
};

// Expressions

expression: binary_expression | function_call | literal | variable;

literal: token_string
{
	$$ = newNode(NodeOpLiteral, $1)
}
| token_number
{
	$$ = newNode(NodeOpLiteral, $1)
}
| token_bool
{
	$$ = newNode(NodeOpLiteral, $1)
};

variable: token_identifier
{
	$$ = newNode(NodeOpVariable, $1)
};

operator: Plus {
	$$ = "+"
}| Minus
{
	$$ = "-"
}
| Star
{
	$$ = "*"
}
| Slash
{
	$$ = "/"
}
| Percent
{
	$$ = "%"
}
| Exclamation
{
	$$ = "!"
}
| And
{
	$$ = "&&"
}
| LessThan
{
	$$ = "<"
};
//| Or | Xor | Equal Equal | GreaterThan | LessThan | GreaterThan Equal | LessThan Equal;

binary_expression: expression operator expression
{
	$$ = appendNode(NodeOpBinaryExpression, $1, $3, newNode(NodeOpOperator, $2))
};

function_call: token_identifier ParenOpen function_params ParenClose
{
	$$ = newNode(NodeOpFunctionCall, $1)
};

function_params: expression
{
	$$ = appendNode(NodeOpFunctionParams, $1)
}
| function_params Comma expression
{
	$$ = appendNodeTo(&$1, $3)
}
| // empty
{
	$$ = appendNode(NodeOpFunctionParams)
}
;

// ################################################################




%%

// 	entity User {
//   		id int [primary, autoincrement]
//  		name string [required, default "John Doe"]
//   	}