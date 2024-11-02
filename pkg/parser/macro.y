%{
package parser

import (
)

%}

%union {
	node yaccMacroNode
	bool bool
	number interface{}
	string string
}

%token<number> token_number
%token<string> token_string
%token<string> token_identifier
%token<bool> token_bool

%token DefinitionKeyword SyntaxKeyword MacroKeyword BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dash Dot Arrow ParenOpen ParenClose Eol

%type<node> macro macro_signature macro_body definition_definition syntax_definition syntax_body syntax_content syntax_statement

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
// End Helpers

file: macro eol_allowed {
	registerRootNode(yylex, $1)
}
| file eol_allowed
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

	definition_definition eol_required

	syntax_definition eol_required

	BraceClose eol_allowed
{
	assertEqual(yylex, $3, "kind", "First identifier in macro body must be 'kind'")
	$$ = appendNode(NodeOpBody, newNode(NodeOpKind, $4), $6, $8)
};

definition_definition: DefinitionKeyword syntax_body
{
	$$ = appendNode(NodeOpSyntax, $2)
}
| // empty
{
	$$ = appendNode(NodeOpDefinition)
};

syntax_definition: SyntaxKeyword syntax_body
{
	$$ = appendNode(NodeOpSyntax, $2)
}
| // empty
{
	$$ = appendNode(NodeOpSyntax)
}

// Syntax definition
syntax_body: BraceOpen eol_allowed syntax_content eol_allowed BraceClose eol_required
{
	$$ = appendNode(NodeOpDefinition)
}

syntax_content: syntax_statement {
        $$ = appendNode(NodeOpDefinition, newNode(NodeOpSyntaxElement, $1))
}
| syntax_statement eol_allowed syntax_content eol_allowed {
	$$ = appendNodeTo(&$3, $1)
}
| {
	$$ = appendNode(NodeOpDefinition)
}
;

syntax_statement: token_identifier
{
	$$ = appendNode(NodeOpSyntax, newNode(NodeOpName, $1))
}

%%
