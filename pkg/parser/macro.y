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

%token macro_keyword BraceOpen BraceClose Comma Colon Semicolon Equal GreaterThan LessThan Dash Dot Arrow ParenOpen ParenClose

%type<node> syntax_macro_definition syntax_macro_signature syntax_macro_body

%start file

%%

file: syntax_macro_definition {
	registerRootNode(yylex, $1)
}
| syntax_macro_definition file{
	registerRootNode(yylex, $1)
}
;

syntax_macro_definition: syntax_macro_signature {
	$$ = appendNode(NodeOpMacro, $1)
};

syntax_macro_signature: macro_keyword token_identifier syntax_macro_body
{
	$$ = appendNode(NodeOpSignature, newNode(NodeOpName, $2), $3)
};

syntax_macro_body: BraceOpen token_identifier token_identifier BraceClose
{
	assertEqual(yylex, $2, "kind", "First identifier in macro body must be 'kind'")
	$$ = appendNode(NodeOpBody, newNode(NodeOpKind, $3))
};

%%
