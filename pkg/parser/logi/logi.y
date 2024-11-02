%{
package parser

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
%token DefinitionKeyword SyntaxKeyword MacroKeyword

%token BracketOpen BracketClose BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dash Dot Arrow ParenOpen ParenClose Eol

%type<node> macro macro_signature macro_body definition_definition syntax_definition syntax_body syntax_content type_definition
%type<node> syntax_statement syntax_element syntax_element_variable_keyword syntax_element_keyword syntax_element_parameter_list syntax_element_parameter_list_content
%type<node> syntax_element_argument_list syntax_element_argument_list_content syntax_element_code_block syntax_element_attribute_list
%type<node> syntax_element_attribute_list_content syntax_element_attribute_list_item

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
	$$ = appendNode(NodeOpSignature, newNode(NodeOpName, $2))
};

macro_body: BraceOpen eol_allowed
	token_identifier token_identifier eol_required

	definition_definition eol_allowed

	syntax_definition eol_allowed

	BraceClose eol_allowed
{
	assertEqual(yylex, $3, "kind", "First identifier in macro body must be 'kind'")
	$$ = appendNode(NodeOpBody, newNode(NodeOpKind, $4), $6, $8)
};

definition_definition: DefinitionKeyword syntax_body eol_required
{
	$$ = appendNode(NodeOpDefinition, $2)
}
| // empty
{
	$$ = appendNode(NodeOpDefinition)
};

syntax_definition: SyntaxKeyword syntax_body eol_required
{
	$$ = appendNode(NodeOpSyntax, $2)
}
| // empty
{
	$$ = appendNode(NodeOpSyntax)
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
	$$ = appendNode(NodeOpBody)
};

syntax_statement: syntax_element
{
	$$ = appendNode(NodeOpSyntaxStatement, $1)
}
| syntax_statement syntax_element
{
	$$ = appendNodeTo(&$1, $2)
};

syntax_element:syntax_element_code_block | syntax_element_keyword | syntax_element_variable_keyword | syntax_element_parameter_list | syntax_element_argument_list | syntax_element_attribute_list ;

syntax_element_keyword: token_identifier
{
	$$ = newNode(NodeOpSyntaxKeywordElement, $1)
};

syntax_element_variable_keyword: LessThan token_identifier type_definition GreaterThan
{
	$$ = appendNode(NodeOpSyntaxVariableKeywordElement, newNode(NodeOpName, $2), $3)
};

syntax_element_parameter_list: ParenOpen syntax_element_parameter_list_content ParenClose
{
	$$ = $2
};

syntax_element_parameter_list_content: syntax_element_variable_keyword
{
	$$ = appendNode(NodeOpSyntaxParameterListElement, $1)
}
| syntax_element_parameter_list_content Comma syntax_element_variable_keyword
{
	$$ = appendNodeTo(&$1, $3);
};

syntax_element_attribute_list: BracketOpen syntax_element_attribute_list_content BracketClose
{
	$$ = $2
};

syntax_element_attribute_list_content: syntax_element_attribute_list_item
{
	$$ = appendNode(NodeOpSyntaxAttributeListElement, $1)
}
| syntax_element_attribute_list_content Comma syntax_element_attribute_list_item
{
	$$ = appendNodeTo(&$1, $3);
};

syntax_element_attribute_list_item: token_identifier
{
	$$ = newNode(NodeOpName, $1)
}
| token_identifier type_definition
{
	$$ = newNode(NodeOpValue, $1, $2)
};

syntax_element_argument_list: ParenOpen three_dots BracketOpen syntax_element_argument_list_content BracketClose ParenClose
{
	$$ = $4
};

syntax_element_argument_list_content: syntax_element_variable_keyword
{
	$$ = appendNode(NodeOpSyntaxArgumentListElement, $1)
}
| syntax_element_argument_list_content Comma syntax_element_variable_keyword
{
	$$ = appendNodeTo(&$1, $3);
};

syntax_element_code_block: BraceOpen BraceClose
{
	$$ = newNode(NodeOpSyntaxCodeBlockElement, nil)
}
| BraceOpen type_definition BraceClose
{
	$$ = newNode(NodeOpSyntaxCodeBlockElement, nil, $2)
};

type_definition: token_identifier
{
	$$ = newNode(NodeOpTypeDef, $1)
}
| token_identifier LessThan type_definition GreaterThan
{
	$$ = newNode(NodeOpTypeDef, $1, $3)
};

%%
