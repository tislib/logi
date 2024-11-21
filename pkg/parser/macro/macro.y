%{
package macro

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
%token TypesKeyword SyntaxKeyword MacroKeyword

// Braces
%token BracketOpen BracketClose BraceOpen BraceClose Comma Colon Semicolon ParenOpen ParenClose Eol CodeBlock ExpressionBlock

// Opeartors
%token Equal GreaterThan LessThan Dash Dot Arrow Or

%type<node> macro macro_signature macro_body syntax_definition syntax_body syntax_content type_definition types_definition_content
%type<node> types_definition types_definition_body types_definition_content types_definition_statement
%type<node> syntax_statement syntax_element syntax_element_combination syntax_element_structure syntax_element_structure_content syntax_element_type_reference syntax_element_combination_content syntax_element_variable_keyword syntax_element_keyword syntax_element_parameter_list syntax_element_parameter_list_content
%type<node> syntax_element_argument_list syntax_element_argument_list_content syntax_element_code_block syntax_element_attribute_list
%type<node> syntax_element_attribute_list_content syntax_element_attribute_list_item syntax_element_expression_block

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

three_dots: Dot Dot Dot
// End Helpers

file: macro eol_allowed {
	registerRootNode(yylex, $1)
}
| file eol_allowed
| eol_allowed
| file macro eol_allowed {
	registerRootNode(yylex, $2)
}
;

// Macro definition
macro: macro_signature eol_allowed macro_body {
	$$ = appendNode(NodeOpMacro, $1, $3)
};

macro_signature: MacroKeyword token_identifier
{
	$$ = newNode(NodeOpSignature, nil, yyDollar[1].token, yyDollar[1].location, newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location))
};

macro_body: BraceOpen eol_allowed
	token_identifier token_identifier eol_required

	types_definition eol_allowed

	syntax_definition eol_allowed

	BraceClose eol_allowed
{
	assertEqual(yylex, $3, "kind", "First identifier in macro body must be 'kind'")
	$$ = appendNode(NodeOpBody, newNode(NodeOpKind, $4, yyDollar[4].token, yyDollar[4].location), $6, $8)
};

types_definition: TypesKeyword types_definition_body eol_required
{
	$$ = newNode(NodeOpTypes, nil, yyDollar[1].token, yyDollar[1].location, $2)
}
| // empty
{
	$$ = newNode(NodeOpTypes, nil, emptyToken, emptyLocation)
};

types_definition_body: BraceOpen eol_allowed types_definition_content eol_allowed BraceClose
{
$$ = $3
};

types_definition_content: types_definition_statement eol_required {
        $$ = appendNode(NodeOpBody, $1)
}
| types_definition_content types_definition_statement eol_required {
	$$ = appendNodeTo(&$1, $2)
}
| // empty
{
	$$ = newNode(NodeOpBody, nil, emptyToken, emptyLocation)
};

types_definition_statement: token_identifier syntax_statement
{
	$$ = appendNode(NodeOpTypesStatement, newNode(NodeOpName, $1, yyDollar[1].token, yyDollar[1].location), $2)
};

syntax_definition: SyntaxKeyword syntax_body eol_required
{
	$$ = newNode(NodeOpSyntax, nil, yyDollar[1].token, yyDollar[1].location, $2)
}
| // empty
{
	$$ = newNode(NodeOpSyntax, nil, emptyToken, emptyLocation)
};

// Syntax definition
syntax_body: BraceOpen eol_allowed syntax_content eol_allowed BraceClose
{
	$$ = $3
};

syntax_content: syntax_statement eol_required {
        $$ = appendNode(NodeOpBody, $1)
}
| syntax_content syntax_statement eol_required {
	$$ = appendNodeTo(&$1, $2)
}
| {
	$$ = newNode(NodeOpBody, nil, emptyToken, emptyLocation)
};

syntax_statement: syntax_element
{
	$$ = appendNode(NodeOpSyntaxStatement, $1)
}
| syntax_statement syntax_element
{
	$$ = appendNodeTo(&$1, $2)
};

syntax_element:syntax_element_code_block | syntax_element_expression_block | syntax_element_combination | syntax_element_structure | syntax_element_type_reference | syntax_element_keyword | syntax_element_variable_keyword | syntax_element_parameter_list | syntax_element_argument_list | syntax_element_attribute_list ;

syntax_element_type_reference: LessThan token_identifier GreaterThan
{
	$$ = newNode(NodeOpSyntaxTypeReferenceElement, $2, yyDollar[2].token, yyDollar[2].location)
};

syntax_element_combination: ParenOpen syntax_element_combination_content ParenClose
{
	$$ = $2
};

syntax_element_combination_content: syntax_element Or syntax_element
{
	$$ = appendNode(NodeOpSyntaxCombinationElement, $1, $3)
}
| syntax_element_combination_content Or syntax_element
{
	$$ = appendNodeTo(&$1, $3)
};

syntax_element_structure: BraceOpen eol_required syntax_element_structure_content BraceClose
{
	$$ = $3
};

syntax_element_structure_content: syntax_statement eol_required
{
	$$ = appendNode(NodeOpSyntaxStructureElement, $1)
}
| syntax_element_structure_content syntax_statement eol_required
{
	$$ = appendNodeTo(&$1, $2)
};

syntax_element_keyword: token_identifier
{
	$$ = newNode(NodeOpSyntaxKeywordElement, $1, yyDollar[1].token, yyDollar[1].location)
};

syntax_element_variable_keyword: LessThan token_identifier type_definition GreaterThan
{
	$$ = appendNode(NodeOpSyntaxVariableKeywordElement, newNode(NodeOpName, $2, yyDollar[2].token, yyDollar[2].location), $3)
};

syntax_element_parameter_list: ParenOpen syntax_element_parameter_list_content ParenClose
{
	$$ = $2
};

syntax_element_parameter_list_content: syntax_element_variable_keyword
{
	$$ = appendNode(NodeOpSyntaxParameterListElement, $1)
}
| syntax_element_parameter_list_content Comma eol_allowed syntax_element_variable_keyword
{
	$$ = appendNodeTo(&$1, $4);
};

syntax_element_attribute_list: BracketOpen eol_allowed syntax_element_attribute_list_content eol_allowed BracketClose
{
	$$ = $3
};

syntax_element_attribute_list_content: syntax_element_attribute_list_item
{
	$$ = appendNode(NodeOpSyntaxAttributeListElement, $1)
}
| syntax_element_attribute_list_content Comma eol_allowed syntax_element_attribute_list_item
{
	$$ = appendNodeTo(&$1, $4);
};

syntax_element_attribute_list_item: token_identifier
{
	$$ = newNode(NodeOpName, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_identifier type_definition
{
	$$ = newNode(NodeOpValue, $1, yyDollar[1].token, yyDollar[1].location, $2)
};

syntax_element_argument_list: ParenOpen three_dots BracketOpen eol_allowed syntax_element_argument_list_content eol_allowed BracketClose ParenClose
{
	$$ = $5
};

syntax_element_argument_list_content: syntax_element_variable_keyword
{
	$$ = appendNode(NodeOpSyntaxArgumentListElement, $1)
}
| syntax_element_argument_list_content Comma eol_allowed syntax_element_variable_keyword
{
	$$ = appendNodeTo(&$1, $4);
};

syntax_element_code_block: CodeBlock
{
	$$ = newNode(NodeOpSyntaxCodeBlockElement, nil, yyDollar[1].token, yyDollar[1].location)
};

syntax_element_expression_block: ExpressionBlock
{
	$$ = newNode(NodeOpSyntaxExpressionBlockElement, nil, yyDollar[1].token, yyDollar[1].location)
};

type_definition: token_identifier
{
	$$ = newNode(NodeOpTypeDef, $1, yyDollar[1].token, yyDollar[1].location)
}
| token_identifier LessThan type_definition GreaterThan
{
	$$ = newNode(NodeOpTypeDef, $1, yyDollar[1].token, yyDollar[1].location, $3)
};

%%
