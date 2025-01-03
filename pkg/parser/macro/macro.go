// Code generated by goyacc -o macro.go macro.y. DO NOT EDIT.

//line macro.y:2
package macro

import __yyfmt__ "fmt"

//line macro.y:2

import (
	"github.com/tislib/logi/pkg/parser/lexer"
)

//line macro.y:10
type yySymType struct {
	yys      int
	node     yaccNode
	bool     bool
	number   interface{}
	string   string
	token    lexer.Token
	location lexer.Location
}

const token_number = 57346
const token_string = 57347
const token_identifier = 57348
const token_bool = 57349
const TypesKeyword = 57350
const SyntaxKeyword = 57351
const MacroKeyword = 57352
const ScopesKeyword = 57353
const BracketOpen = 57354
const BracketClose = 57355
const BraceOpen = 57356
const BraceClose = 57357
const Comma = 57358
const Colon = 57359
const Semicolon = 57360
const ParenOpen = 57361
const ParenClose = 57362
const Eol = 57363
const Equal = 57364
const GreaterThan = 57365
const LessThan = 57366
const Dash = 57367
const Dot = 57368
const Arrow = 57369
const Or = 57370
const Hash = 57371

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"token_number",
	"token_string",
	"token_identifier",
	"token_bool",
	"TypesKeyword",
	"SyntaxKeyword",
	"MacroKeyword",
	"ScopesKeyword",
	"BracketOpen",
	"BracketClose",
	"BraceOpen",
	"BraceClose",
	"Comma",
	"Colon",
	"Semicolon",
	"ParenOpen",
	"ParenClose",
	"Eol",
	"Equal",
	"GreaterThan",
	"LessThan",
	"Dash",
	"Dot",
	"Arrow",
	"Or",
	"Hash",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line macro.y:379

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 80,
	28, 56,
	-2, 77,
}

const yyPrivate = 57344

const yyLast = 229

var yyAct = [...]uint8{
	3, 51, 7, 9, 115, 11, 19, 94, 126, 13,
	112, 93, 119, 44, 103, 32, 16, 35, 101, 46,
	109, 133, 104, 24, 110, 135, 10, 30, 31, 148,
	113, 102, 29, 160, 40, 41, 134, 85, 65, 39,
	84, 116, 43, 162, 5, 159, 21, 111, 137, 71,
	5, 10, 42, 5, 69, 120, 10, 154, 152, 80,
	15, 83, 20, 10, 120, 73, 10, 10, 88, 89,
	5, 10, 59, 87, 38, 121, 91, 79, 60, 10,
	56, 10, 90, 64, 10, 57, 86, 117, 62, 61,
	58, 63, 10, 70, 28, 72, 36, 122, 36, 10,
	106, 124, 17, 23, 105, 123, 10, 130, 127, 131,
	6, 5, 113, 10, 108, 136, 21, 10, 164, 139,
	144, 5, 107, 145, 129, 33, 141, 138, 67, 132,
	26, 140, 147, 150, 146, 143, 142, 157, 153, 128,
	96, 97, 95, 98, 82, 75, 151, 155, 100, 18,
	158, 124, 12, 127, 156, 59, 78, 2, 161, 8,
	163, 60, 1, 56, 146, 165, 64, 125, 57, 99,
	59, 62, 61, 58, 63, 81, 60, 118, 56, 66,
	37, 64, 114, 57, 59, 5, 62, 61, 58, 63,
	60, 54, 56, 149, 53, 64, 74, 57, 59, 10,
	62, 61, 58, 63, 60, 47, 56, 77, 52, 64,
	55, 57, 50, 76, 62, 61, 58, 63, 49, 48,
	45, 92, 25, 22, 34, 68, 27, 14, 4,
}

var yyPact = [...]int16{
	100, 100, 23, 50, 23, -1000, 146, 50, 23, 50,
	-1000, 46, -1000, 50, -1000, 23, 96, 143, 41, 95,
	-1000, -1000, 23, 116, 85, 41, 23, 23, 111, 25,
	92, 63, 41, 23, 90, 41, 192, 23, 114, 25,
	178, 78, 41, 25, -1000, 66, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 139, 149, 138, -1000,
	23, -1000, 17, 14, -1000, 71, 41, 23, 164, 41,
	-1000, 25, 136, -1000, 3, -1000, -6, 84, 102, -8,
	-1000, -2, 24, 35, -1000, -1000, 23, 25, 58, 60,
	41, 25, 89, 136, -1000, -1000, -1000, -1000, -1000, -1000,
	136, -1000, 133, -1000, 192, -1000, 23, -1000, 23, 192,
	-5, -1000, 13, 1, 32, -1000, 106, 50, 49, 41,
	111, -1000, 25, 136, -1000, 107, 136, -1000, -1000, -1000,
	5, 5, -1000, -1000, -1000, 106, 45, 23, -1000, 42,
	41, 25, -1000, 136, -1000, 136, -1000, -1000, 131, 29,
	-1000, 10, -1000, 35, -1000, 25, 136, 106, 30, 23,
	-1000, -1000, 98, 5, -1000, -1000,
}

var yyPgo = [...]uint8{
	0, 157, 228, 227, 226, 15, 225, 10, 224, 223,
	222, 17, 13, 11, 221, 220, 19, 219, 218, 213,
	1, 212, 210, 208, 207, 205, 196, 194, 193, 191,
	182, 4, 7, 180, 179, 177, 12, 169, 167, 8,
	162, 0, 6, 156,
}

var yyR1 = [...]int8{
	0, 41, 41, 41, 42, 42, 43, 40, 40, 40,
	40, 1, 2, 3, 33, 33, 34, 35, 35, 35,
	36, 9, 9, 10, 8, 8, 8, 11, 4, 4,
	5, 6, 6, 6, 12, 12, 14, 14, 13, 13,
	32, 32, 32, 32, 32, 37, 39, 39, 38, 38,
	15, 15, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 25, 26, 26, 18, 17, 19, 19, 21, 22,
	22, 22, 22, 22, 20, 23, 23, 24, 24, 24,
	29, 30, 30, 31, 31, 27, 28, 28, 7, 7,
}

var yyR2 = [...]int8{
	0, 1, 2, 0, 1, 2, 3, 2, 2, 1,
	3, 3, 2, 13, 3, 0, 5, 2, 3, 0,
	2, 3, 0, 5, 2, 3, 0, 2, 3, 0,
	5, 2, 3, 0, 1, 3, 1, 3, 1, 2,
	1, 1, 1, 1, 1, 3, 1, 2, 1, 3,
	1, 2, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 1, 3, 3, 3, 3, 3, 1, 1,
	1, 2, 2, 1, 4, 3, 3, 1, 4, 0,
	5, 1, 4, 1, 2, 8, 1, 4, 1, 4,
}

var yyChk = [...]int16{
	-1000, -40, -1, -41, -2, 21, 10, -41, -1, -41,
	21, -41, 6, -41, -3, 14, -41, 6, 6, -42,
	21, 21, -9, 8, -41, -10, 14, -4, 9, -42,
	-41, -41, -5, 14, -8, -11, 6, -33, 11, -42,
	-41, -41, -11, -42, -12, -15, -16, -25, -17, -18,
	-21, -20, -23, -27, -29, -22, 14, 19, 24, 6,
	12, 23, 22, 25, 17, -41, -34, 14, -6, -12,
	15, -42, 29, -16, -26, 6, -19, -24, -43, -16,
	-20, 26, 6, -41, 23, 23, 15, -42, -41, -41,
	-12, -42, -14, -13, -32, 6, 4, 5, 7, -37,
	12, 15, 28, 20, 28, 20, 16, 20, 12, 28,
	26, 23, -7, 6, -30, -31, 6, -41, -35, -36,
	6, 15, -42, 16, -32, -38, -39, -32, 6, -16,
	-41, -41, -16, 26, 23, 24, -41, 16, -7, -41,
	-36, -42, -5, -13, 13, 16, -32, -20, 24, -28,
	-20, -7, 13, -41, 15, -42, -39, 6, -41, 16,
	23, -31, 13, -41, 20, -20,
}

var yyDef = [...]int8{
	3, -2, 3, 9, 3, 1, 0, 8, 3, 7,
	2, 0, 12, 10, 11, 3, 0, 0, 0, 22,
	4, 5, 3, 0, 29, 0, 3, 3, 0, 21,
	26, 15, 0, 3, 3, 0, 0, 3, 0, 28,
	33, 0, 0, 24, 27, 34, 50, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 0, 79, 70, 68,
	3, 69, 0, 0, 73, 0, 0, 3, 3, 0,
	23, 25, 0, 51, 0, 62, 0, 0, 0, 0,
	-2, 0, 0, 0, 71, 72, 3, 14, 19, 0,
	0, 31, 35, 36, 38, 40, 41, 42, 43, 44,
	0, 61, 0, 65, 0, 75, 3, 76, 3, 0,
	0, 64, 0, 88, 3, 81, 83, 13, 3, 0,
	0, 30, 32, 0, 39, 0, 48, 46, 63, 67,
	0, 0, 66, 6, 74, 0, 0, 3, 84, 0,
	0, 17, 20, 37, 45, 0, 47, 78, 0, 3,
	86, 0, 80, 0, 16, 18, 49, 0, 0, 3,
	89, 82, 0, 0, 85, 87,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29,
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

	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:58
		{
			registerRootNode(yylex, yyDollar[1].node)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:63
		{
			registerRootNode(yylex, yyDollar[2].node)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:69
		{
			yyVAL.node = appendNode(NodeOpMacro, yyDollar[1].node, yyDollar[3].node)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:74
		{
			yyVAL.node = newNode(NodeOpSignature, nil, yyDollar[1].token, yyDollar[1].location, newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location))
		}
	case 13:
		yyDollar = yyS[yypt-13 : yypt+1]
//line macro.y:88
		{
			assertEqual(yylex, yyDollar[3].string, "kind", "First identifier in macro body must be 'kind'")
			yyVAL.node = appendNode(NodeOpBody, newNode(NodeOpKind, yyDollar[4].string, yyDollar[4].token, yyDollar[4].location), yyDollar[6].node, yyDollar[8].node, yyDollar[10].node)
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:94
		{
			yyVAL.node = newNode(NodeOpScopes, nil, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:98
		{
			yyVAL.node = newNode(NodeOpScopes, nil, emptyToken, emptyLocation)
		}
	case 16:
		yyDollar = yyS[yypt-5 : yypt+1]
//line macro.y:103
		{
			yyVAL.node = yyDollar[3].node
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:107
		{
			yyVAL.node = appendNodeX(NodeOpBody, yyDollar[1].node)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:109
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:112
		{
			yyVAL.node = newNode(NodeOpBody, nil, emptyToken, emptyLocation)
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:117
		{
			yyVAL.node = appendNode(NodeOpScopesItem, newNode(NodeOpName, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location), yyDollar[2].node)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:122
		{
			yyVAL.node = newNode(NodeOpTypes, nil, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 22:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:126
		{
			yyVAL.node = newNode(NodeOpTypes, nil, emptyToken, emptyLocation)
		}
	case 23:
		yyDollar = yyS[yypt-5 : yypt+1]
//line macro.y:131
		{
			yyVAL.node = yyDollar[3].node
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:135
		{
			yyVAL.node = appendNode(NodeOpBody, yyDollar[1].node)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:138
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 26:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:142
		{
			yyVAL.node = newNode(NodeOpBody, nil, emptyToken, emptyLocation)
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:147
		{
			yyVAL.node = appendNode(NodeOpTypesStatement, newNode(NodeOpName, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location), yyDollar[2].node)
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:152
		{
			yyVAL.node = newNode(NodeOpSyntax, nil, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:156
		{
			yyVAL.node = newNode(NodeOpSyntax, nil, emptyToken, emptyLocation)
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
//line macro.y:162
		{
			yyVAL.node = yyDollar[3].node
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:166
		{
			yyVAL.node = appendNode(NodeOpBody, yyDollar[1].node)
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:169
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:172
		{
			yyVAL.node = newNode(NodeOpBody, nil, emptyToken, emptyLocation)
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:177
		{
			yyVAL.node = appendNode(NodeOpSyntaxStatement, yyDollar[1].node)
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:180
		{
			yyVAL.node = appendNode(NodeOpSyntaxStatement, yyDollar[1].node, yyDollar[3].node)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:185
		{
			yyVAL.node = appendNode(NodeOpSyntaxExamples, yyDollar[1].node)
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:188
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:193
		{
			yyVAL.node = appendNode(NodeOpSyntaxExample, yyDollar[1].node)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:196
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:201
		{
			yyVAL.node = newNode(NodeOpValueIdentifier, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:205
		{
			yyVAL.node = newNode(NodeOpValueNumber, yyDollar[1].number, yyDollar[1].token, yyDollar[1].location)
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:208
		{
			yyVAL.node = newNode(NodeOpValueString, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:211
		{
			yyVAL.node = newNode(NodeOpValueBool, yyDollar[1].bool, yyDollar[1].token, yyDollar[1].location)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:214
		{
			yyVAL.node = yyDollar[1].node
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:219
		{
			yyVAL.node = yyDollar[2].node
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:224
		{
			yyVAL.node = appendNode(NodeOpValueArrayItem, yyDollar[1].node)
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:227
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:232
		{
			yyVAL.node = appendNode(NodeOpValueArray, yyDollar[1].node)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:235
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:240
		{
			yyVAL.node = appendNode(NodeOpSyntaxElements, yyDollar[1].node)
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:244
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:251
		{
			yyVAL.node = yyDollar[2].node
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:256
		{
			yyVAL.node = appendNode(NodeOpSyntaxScopeElement, newNode(NodeOpName, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location))
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:259
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, newNode(NodeOpName, yyDollar[3].string, yyDollar[3].token, yyDollar[3].location))
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:265
		{
			yyVAL.node = newNode(NodeOpSyntaxTypeReferenceElement, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:270
		{
			yyVAL.node = yyDollar[2].node
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:275
		{
			yyVAL.node = appendNode(NodeOpSyntaxCombinationElement, yyDollar[1].node, yyDollar[3].node)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:279
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:284
		{
			yyVAL.node = newNode(NodeOpSyntaxKeywordElement, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:289
		{
			yyVAL.node = newNode(NodeOpSyntaxSymbolElement, ">", yyDollar[1].token, yyDollar[1].location)
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:292
		{
			yyVAL.node = newNode(NodeOpSyntaxSymbolElement, "<", yyDollar[1].token, yyDollar[1].location)
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:295
		{
			yyVAL.node = newNode(NodeOpSyntaxSymbolElement, "=>", yyDollar[1].token, yyDollar[1].location)
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:298
		{
			yyVAL.node = newNode(NodeOpSyntaxSymbolElement, "->", yyDollar[1].token, yyDollar[1].location)
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:301
		{
			yyVAL.node = newNode(NodeOpSyntaxSymbolElement, ":", yyDollar[1].token, yyDollar[1].location)
		}
	case 74:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:307
		{
			yyVAL.node = appendNode(NodeOpSyntaxVariableKeywordElement, newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location), yyDollar[3].node)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:312
		{
			yyVAL.node = yyDollar[2].node
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
//line macro.y:316
		{
			yyVAL.node = newNode(NodeOpSyntaxParameterListElement, true, emptyToken, emptyLocation)
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:321
		{
			yyVAL.node = appendNode(NodeOpSyntaxParameterListElement, yyDollar[1].node)
		}
	case 78:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:325
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 79:
		yyDollar = yyS[yypt-0 : yypt+1]
//line macro.y:329
		{
			yyVAL.node = newNode(NodeOpSyntaxParameterListElement, nil, emptyToken, emptyLocation)
		}
	case 80:
		yyDollar = yyS[yypt-5 : yypt+1]
//line macro.y:334
		{
			yyVAL.node = yyDollar[3].node
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:339
		{
			yyVAL.node = appendNode(NodeOpSyntaxAttributeListElement, yyDollar[1].node)
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:343
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:348
		{
			yyVAL.node = newNode(NodeOpName, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
//line macro.y:352
		{
			yyVAL.node = newNode(NodeOpName, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 85:
		yyDollar = yyS[yypt-8 : yypt+1]
//line macro.y:357
		{
			yyVAL.node = yyDollar[5].node
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:362
		{
			yyVAL.node = appendNode(NodeOpSyntaxArgumentListElement, yyDollar[1].node)
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:366
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
//line macro.y:371
		{
			yyVAL.node = newNode(NodeOpTypeDef, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 89:
		yyDollar = yyS[yypt-4 : yypt+1]
//line macro.y:375
		{
			yyVAL.node = newNode(NodeOpTypeDef, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[3].node)
		}
	}
	goto yystack /* stack new state and value */
}
