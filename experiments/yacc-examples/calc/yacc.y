%{
package main

import (
	"math"
)

%}

%union {
	val interface{}
}

%token<val> NUMBER
%token PLUS MINUS TIMES DIVIDE POW LP RP
%type<val> expr
%start statement

// Define precedence and associativity
%left PLUS MINUS
%left TIMES DIVIDE
%right POW

%%

statement: expr
    {
    	// Print the result
    	statement($1)
    };

expr: LP expr RP
         {
         	$$ = ($2)
         }
| expr POW expr
     {
     	$$ = math.Pow(toFloat64($1), toFloat64($3))
     }
| expr TIMES expr
    {
    	$$ = toFloat64($1) * toFloat64($3)
    }
| expr DIVIDE expr
    {
    	$$ = toFloat64($1) / toFloat64($3)
    }
| expr PLUS expr
    {
	$$ = toFloat64($1) + toFloat64($3)
    }
| expr MINUS expr
    {
    	$$ = toFloat64($1) - toFloat64($3)
    }
|NUMBER
    {
    	$$ = $1
    }
;

%%
