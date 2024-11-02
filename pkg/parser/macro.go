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

//line macro.y:111

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 50

var yyAct = [...]int8{
	5, 31, 18, 8, 9, 33, 24, 12, 38, 34,
	27, 22, 16, 11, 4, 15, 20, 7, 14, 19,
	11, 20, 25, 11, 23, 20, 28, 20, 7, 29,
	11, 11, 35, 36, 30, 37, 4, 17, 39, 10,
	41, 40, 2, 1, 6, 32, 26, 21, 13, 3,
}

var yyPact = [...]int16{
	26, 4, -7, -7, 33, -11, -7, -1000, -11, 7,
	-1000, -1000, -11, -1000, -7, 6, 31, -5, 3, -1000,
	-1000, -5, 11, 1, -1000, -7, -5, 11, -1, -3,
	-1000, -7, -7, -1000, -7, -4, -1, -11, -5, -7,
	-8, -11,
}

var yyPgo = [...]int8{
	0, 42, 49, 48, 47, 46, 6, 1, 45, 43,
	0, 2,
}

var yyR1 = [...]int8{
	0, 10, 10, 10, 11, 11, 9, 9, 9, 1,
	2, 3, 4, 4, 5, 5, 6, 7, 7, 7,
	8,
}

var yyR2 = [...]int8{
	0, 1, 2, 0, 1, 2, 2, 2, 3, 3,
	2, 11, 2, 0, 2, 0, 6, 1, 4, 0,
	1,
}

var yyChk = [...]int16{
	-1000, -9, -1, -2, 10, -10, -1, 24, -10, -10,
	6, 24, -10, -3, 11, -10, 6, 6, -11, 24,
	24, -4, 8, -11, -6, 11, -5, 9, -10, -11,
	-6, -7, -8, 6, 12, -10, -10, -10, 12, -7,
	-11, -10,
}

var yyDef = [...]int8{
	0, -2, 3, 3, 0, 7, 3, 1, 6, 0,
	10, 2, 8, 9, 3, 0, 0, 0, 0, 4,
	5, 0, 0, 0, 12, 3, 0, 0, 19, 0,
	14, 3, 3, 20, 3, 0, 19, 11, 0, 3,
	0, 18,
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
//line macro.y:39
		{
			registerRootNode(yylex, yyDollar[1].node)
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:43
		{
			registerRootNode(yylex, yyDollar[2].node)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:49
		{
			yyVAL.node = appendNode(NodeOpMacro, yyDollar[1].node, yyDollar[3].node)
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:54
		{
			yyVAL.node = appendNode(NodeOpSignature, newNode(NodeOpName, yyDollar[2].string))
		}
	case 11:
		yyDollar = yyS[yypt-11 : yypt+1]
//line macro.y:66
		{
			assertEqual(yylex, yyDollar[3].string, "kind", "First identifier in macro body must be 'kind'")
			yyVAL.node = appendNode(NodeOpBody, newNode(NodeOpKind, yyDollar[4].string), yyDollar[6].node, yyDollar[8].node)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:72
		{
			yyVAL.node = appendNode(NodeOpSyntax, yyDollar[2].node)
		}
	case 13:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:76
		{
			yyVAL.node = appendNode(NodeOpDefinition)
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:81
		{
			yyVAL.node = appendNode(NodeOpSyntax, yyDollar[2].node)
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:85
		{
			yyVAL.node = appendNode(NodeOpSyntax)
		}
	case 16:
		yyDollar = yyS[yypt-6 : yypt+1]
//line macro.y:91
		{
			yyVAL.node = appendNode(NodeOpDefinition)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:95
		{
			yyVAL.node = appendNode(NodeOpDefinition, newNode(NodeOpSyntaxElement, yyDollar[1].node))
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:98
		{
			yyVAL.node = appendNodeTo(&yyDollar[3].node, yyDollar[1].node)
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:101
		{
			yyVAL.node = appendNode(NodeOpDefinition)
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:107
		{
			yyVAL.node = appendNode(NodeOpSyntax, newNode(NodeOpName, yyDollar[1].string))
		}
	}
	goto yystack /* stack new state and value */
}