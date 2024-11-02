// Code generated by goyacc -o macro.go macro.y. DO NOT EDIT.

//line macro.y:2
package parser

import __yyfmt__ "fmt"

//line macro.y:2

import ()

//line macro.y:9
type yySymType struct {
	yys    int
	node   yaccMacroNode
	bool   bool
	number interface{}
	string string
}

const token_number = 57346
const token_string = 57347
const token_identifier = 57348
const token_bool = 57349
const DefinitionKeyword = 57350
const SyntaxKeyword = 57351
const MacroKeyword = 57352
const BraceOpen = 57353
const BraceClose = 57354
const Comma = 57355
const Colon = 57356
const Semicolon = 57357
const Equal = 57358
const GreaterThan = 57359
const LessThan = 57360
const Dash = 57361
const Dot = 57362
const Arrow = 57363
const ParenOpen = 57364
const ParenClose = 57365
const Eol = 57366

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"token_number",
	"token_string",
	"token_identifier",
	"token_bool",
	"DefinitionKeyword",
	"SyntaxKeyword",
	"MacroKeyword",
	"BraceOpen",
	"BraceClose",
	"Comma",
	"Colon",
	"Semicolon",
	"Equal",
	"GreaterThan",
	"LessThan",
	"Dash",
	"Dot",
	"Arrow",
	"ParenOpen",
	"ParenClose",
	"Eol",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line macro.y:147

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 77

var yyAct = [...]int8{
	37, 3, 35, 7, 9, 19, 11, 34, 52, 39,
	13, 42, 21, 15, 25, 10, 5, 16, 39, 20,
	10, 40, 40, 10, 24, 41, 10, 20, 30, 31,
	40, 29, 28, 23, 41, 44, 5, 47, 43, 6,
	46, 45, 50, 32, 51, 39, 17, 10, 47, 21,
	57, 53, 58, 5, 26, 55, 48, 40, 59, 18,
	56, 41, 12, 10, 10, 2, 1, 8, 49, 38,
	36, 54, 33, 27, 22, 14, 4,
}

var yyPact = [...]int16{
	29, 29, -8, -9, -8, -1000, 56, -9, -8, -9,
	-1000, 2, -1000, -9, -1000, -8, 40, 53, -5, 25,
	-1000, -1000, -8, 43, 23, -5, -8, -8, 43, -12,
	39, -1, -5, 12, 3, -1000, -1000, -1000, -1000, -1000,
	50, 4, -8, -12, -4, 3, -12, -1000, 49, 37,
	-1000, -9, -1000, -12, 35, -1000, -1000, 4, -1000, -1000,
}

var yyPgo = [...]int8{
	0, 65, 76, 75, 74, 73, 14, 72, 71, 7,
	2, 0, 70, 69, 68, 66, 1, 5,
}

var yyR1 = [...]int8{
	0, 16, 16, 16, 17, 17, 15, 15, 15, 15,
	1, 2, 3, 4, 4, 5, 5, 6, 7, 7,
	7, 9, 9, 10, 10, 10, 12, 11, 13, 14,
	14, 8,
}

var yyR2 = [...]int8{
	0, 1, 2, 0, 1, 2, 2, 2, 1, 3,
	3, 2, 11, 3, 0, 3, 0, 5, 2, 3,
	0, 1, 2, 1, 1, 1, 1, 4, 3, 1,
	3, 1,
}

var yyChk = [...]int16{
	-1000, -15, -1, -16, -2, 24, 10, -16, -1, -16,
	24, -16, 6, -16, -3, 11, -16, 6, 6, -17,
	24, 24, -4, 8, -16, -6, 11, -5, 9, -17,
	-16, -16, -6, -7, -9, -10, -12, -11, -13, 6,
	18, 22, 12, -17, -16, -9, -17, -10, 6, -14,
	-11, -16, 12, -17, -8, 6, 23, 13, 17, -11,
}

var yyDef = [...]int8{
	3, -2, 3, 8, 3, 1, 0, 7, 3, 6,
	2, 0, 11, 9, 10, 3, 0, 0, 0, 14,
	4, 5, 3, 0, 16, 0, 3, 3, 0, 13,
	20, 0, 0, 3, 0, 21, 23, 24, 25, 26,
	0, 0, 3, 15, 0, 0, 18, 22, 0, 0,
	29, 12, 17, 19, 0, 31, 28, 0, 27, 30,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:40
		{
			registerRootNode(yylex, yyDollar[1].node)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:45
		{
			registerRootNode(yylex, yyDollar[2].node)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:51
		{
			yyVAL.node = appendNode(NodeOpMacro, yyDollar[1].node, yyDollar[3].node)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:56
		{
			yyVAL.node = appendNode(NodeOpSignature, newNode(NodeOpName, yyDollar[2].string))
		}
	case 12:
		yyDollar = yyS[yypt-11 : yypt+1]
//line macro.y:68
		{
			assertEqual(yylex, yyDollar[3].string, "kind", "First identifier in macro body must be 'kind'")
			yyVAL.node = appendNode(NodeOpBody, newNode(NodeOpKind, yyDollar[4].string), yyDollar[6].node, yyDollar[8].node)
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:74
		{
			yyVAL.node = appendNode(NodeOpDefinition, yyDollar[2].node)
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:78
		{
			yyVAL.node = appendNode(NodeOpDefinition)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:83
		{
			yyVAL.node = appendNode(NodeOpSyntax, yyDollar[2].node)
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:87
		{
			yyVAL.node = appendNode(NodeOpSyntax)
		}
	case 17:
		yyDollar = yyS[yypt-5 : yypt+1]
//line macro.y:93
		{
			yyVAL.node = yyDollar[3].node
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:97
		{
			yyVAL.node = appendNode(NodeOpBody, yyDollar[1].node)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:100
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:103
		{
			yyVAL.node = appendNode(NodeOpBody)
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:108
		{
			yyVAL.node = appendNode(NodeOpSyntaxStatement, yyDollar[1].node)
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:112
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:119
		{
			yyVAL.node = newNode(NodeOpSyntaxKeywordElement, yyDollar[1].string)
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:124
		{
			yyVAL.node = appendNode(NodeOpSyntaxVariableKeywordElement, newNode(NodeOpName, yyDollar[2].string), yyDollar[3].node)
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:129
		{
			yyVAL.node = yyDollar[2].node
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:134
		{
			yyVAL.node = appendNode(NodeOpSyntaxParameterListElement, yyDollar[1].node)
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:138
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:143
		{
			yyVAL.node = newNode(NodeOpTypeDef, yyDollar[1].string)
		}
	}
	goto yystack /* stack new state and value */
}
