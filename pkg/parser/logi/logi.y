%{
package logi

import (
	"github.com/tislib/logi/pkg/parser/lexer"
)

%}

%union {
	node yaccNode
	bool bool
	number interface{}
	string string
	token    lexer.Token
	location lexer.Location
}


%token<number> token_number
%token<string> token_string
%token<string> token_identifier
%token<bool> token_bool

// Keywords
%token DefinitionKeyword SyntaxKeyword FuncKeyword

%token BracketOpen BracketClose BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dot Arrow ParenOpen ParenClose Eol
%token IfKeyword ElseKeyword ReturnKeyword SwitchKeyword CaseKeyword VarKeyword
%token Plus Minus Star Slash Percent Exclamation And Or Xor

%type<node> type_definition
%type<node> definition definition_signature definition_body definition_statements definition_statement definition_statement_element
%type<node> definition_statement_element_identifier definition_statement_element_array definition_statement_element_array_content definition_statement_element_struct definition_statement_element_value definition_statement_element_attribute_list definition_statement_element_attribute_list_content definition_statement_element_attribute_list_item
%type<node> definition_statement_element_argument_list definition_statement_element_argument_list_content definition_statement_element_argument_list_item
%type<node> definition_statement_element_parameter_list definition_statement_element_parameter_list_content definition_statement_element_parameter_list_item
%type<node> code_block code_block_statements code_block_statement assignment_statement if_statement return_statement variable_declaration_statement function_call_statement
%type<node> expression literal variable binary_expression function_call
%type<node> definition_statement_element_json json_object json_object_content json_object_item json_array json_value json_array_content
%type<node> function_params
%type<node> function_definition
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
| function_definition eol_allowed {
  	registerRootNode(yylex, $1)
}
| file eol_allowed
| eol_allowed
| file definition eol_allowed {
	registerRootNode(yylex, $2)
}
| file function_definition eol_allowed {
	registerRootNode(yylex, $2)
}
| // empty;

definition: definition_signature eol_allowed definition_body
{
	$$ = appendNode(NodeOpDefinition, $1, $3)
};

definition_signature: token_identifier token_identifier
{
	$$ = appendNode(NodeOpSignature, newNode(NodeOpMacro, $1, yyDollar[1].token, yyDollar[1].location), newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location))
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

definition_statement_element: definition_statement_element_identifier | definition_statement_element_value | definition_statement_element_array | definition_statement_element_struct | definition_statement_element_json | definition_statement_element_parameter_list | definition_statement_element_attribute_list | definition_statement_element_argument_list | code_block;

definition_statement_element_identifier: token_identifier
{
	$$ = newNode(NodeOpIdentifier, $1, yyDollar[1].token, yyDollar[1].location)
};

definition_statement_element_value: token_string
{
	$$ = newNode(NodeOpValue, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_number
{
	$$ = newNode(NodeOpValue, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_bool
{
	$$ = newNode(NodeOpValue, $1, yyDollar[1].token, yyDollar[1].location)
};

definition_statement_element_array: BracketOpen eol_allowed definition_statement_element_array_content eol_allowed BracketClose
{
	$$ = $3
};

definition_statement_element_array_content: definition_statement
{
	$$ = appendNode(NodeOpArray, $1)
}
| definition_statement_element_array_content Comma eol_allowed definition_statement
{
	$$ = appendNodeTo(&$1, $4)
}
| // empty
{
	$$ = appendNode(NodeOpArray)
};

definition_statement_element_struct: definition_body
{
	$$ = appendNode(NodeOpStruct, $1)
};

definition_statement_element_json: json_object
{
	$$ = $1
};

json_object: BraceOpen eol_allowed json_object_content eol_allowed BraceClose
{
	$$ = $3
};

json_object_content: json_object_item
{
	$$ = appendNode(NodeOpJsonObject, $1)
} | json_object_content Comma eol_allowed json_object_item
{
	$$ = appendNodeTo(&$1, $4)
} | // empty
{
	$$ = appendNode(NodeOpJsonObject)
};

json_object_item: token_string Colon json_value
{
	$$ = newNode(NodeOpJsonObjectItem, $1, yyDollar[1].token, yyDollar[1].location, $3)
};

json_value: token_string
{
	$$ = newNode(NodeOpJsonObjectItemValue, $1, yyDollar[1].token, yyDollar[1].location)
} | token_number
{
	$$ = newNode(NodeOpJsonObjectItemValue, $1, yyDollar[1].token, yyDollar[1].location)
} | token_bool
{
	$$ = newNode(NodeOpJsonObjectItemValue, $1, yyDollar[1].token, yyDollar[1].location)
} | json_object
{
	$$ = $1
} | json_array
{
	$$ = $1
};

json_array: BracketOpen eol_allowed json_array_content eol_allowed BracketClose
{
	$$ = $3
};

json_array_content: json_value
{
	$$ = appendNode(NodeOpJsonArray, $1)
} | json_array_content Comma eol_allowed json_value
{
	$$ = appendNodeTo(&$1, $4)
} | // empty
{
	$$ = appendNode(NodeOpJsonArray)
};


definition_statement_element_attribute_list: LessThan BracketOpen definition_statement_element_attribute_list_content BracketClose GreaterThan
{
	$$ = $3
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
	$$ = newNode(NodeOpAttribute, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_identifier definition_statement_element_value
{
	$$ = newNode(NodeOpAttribute, $1, yyDollar[1].token, yyDollar[1].location, $2)
};

definition_statement_element_argument_list: ParenOpen definition_statement_element_argument_list_content ParenClose
{
	$$ = $2
}
| ParenOpen ParenClose
{
	$$ = appendNode(NodeOpArgumentList)
};

definition_statement_element_argument_list_content: definition_statement_element_argument_list_item
{
	$$ = appendNode(NodeOpArgumentList, $1)
}
| definition_statement_element_argument_list_content Comma eol_allowed definition_statement_element_argument_list_item
{
	$$ = appendNodeTo(&$1, $4)
};

definition_statement_element_argument_list_item: token_identifier type_definition
{
	$$ = newNode(NodeOpArgument, $1, yyDollar[1].token, yyDollar[1].location, $2)
};

definition_statement_element_parameter_list: ParenOpen definition_statement_element_parameter_list_content ParenClose
{
	$$ = $2
}
| ParenOpen ParenClose
{
	$$ = appendNode(NodeOpParameterList)
};

definition_statement_element_parameter_list_content: definition_statement_element_parameter_list_item
{
	$$ = appendNode(NodeOpParameterList, $1)
}
| definition_statement_element_parameter_list_content Comma eol_allowed definition_statement_element_parameter_list_item
{
	$$ = appendNodeTo(&$1, $4)
};

definition_statement_element_parameter_list_item: definition_statement
{
	$$ = newNode(NodeOpParameter, $1, yyDollar[1].token, yyDollar[1].location)
};


type_definition: token_identifier
{
	$$ = newNode(NodeOpTypeDef, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_identifier LessThan type_definition GreaterThan
{
	$$ = newNode(NodeOpTypeDef, $1, yyDollar[1].token, yyDollar[1].location, $3)
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

code_block_statement: function_call_statement | assignment_statement | if_statement | return_statement | variable_declaration_statement;

// Statements
function_call_statement: function_call
{
	$$ = appendNode(NodeOpFunctionCall, $1)
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
	$$ = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location), $3)
}
| VarKeyword token_identifier type_definition Equal expression
{
	$$ = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location), $3, $5)
};

// Expressions

expression: binary_expression | function_call | literal | variable;

literal: token_string
{
	$$ = newNode(NodeOpLiteral, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_number
{
	$$ = newNode(NodeOpLiteral, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_bool
{
	$$ = newNode(NodeOpLiteral, $1, yyDollar[1].token, yyDollar[1].location)
};

variable: token_identifier
{
	$$ = newNode(NodeOpVariable, $1, yyDollar[1].token, yyDollar[1].location)
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
}
| GreaterThan
{
	$$ = ">"
}
| Or
{
	$$ = "||"
}
| Xor
{
	$$ = "^"
}
| Equal Equal
{
	$$ = "=="
}
| GreaterThan Equal
{
	$$ = ">="
}
| LessThan Equal
{
	$$ = "<="
};

binary_expression: expression operator expression
{
	$$ = appendNode(NodeOpBinaryExpression, $1, $3, newNode(NodeOpOperator, $2, yyDollar[2].token, yyDollar[2].location))
};

function_call: token_identifier ParenOpen function_params ParenClose
{
	$$ = newNode(NodeOpFunctionCall, $1, yyDollar[1].token, yyDollar[1].location, $3)
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

function_definition: FuncKeyword token_identifier definition_statement_element_argument_list code_block
{
	$$ = appendNode(NodeOpFunctionDefinition, newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location), $3, $4)
};

%%

// 	entity User {
//   		id int [primary, autoincrement]
//  		name string [required, default "John Doe"]
//   	}