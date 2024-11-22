// Code generated by goyacc -o logi.go logi.y. DO NOT EDIT.

//line logi.y:2
package logi

import __yyfmt__ "fmt"

//line logi.y:2

import (
	"github.com/tislib/logi/pkg/parser/lexer"
)

//line logi.y:10
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
const DefinitionKeyword = 57350
const SyntaxKeyword = 57351
const FuncKeyword = 57352
const BracketOpen = 57353
const BracketClose = 57354
const BraceOpen = 57355
const BraceClose = 57356
const Comma = 57357
const Colon = 57358
const Semicolon = 57359
const Equal = 57360
const GreaterThan = 57361
const LessThan = 57362
const Dot = 57363
const Arrow = 57364
const ParenOpen = 57365
const ParenClose = 57366
const Eol = 57367
const IfKeyword = 57368
const ElseKeyword = 57369
const ReturnKeyword = 57370
const SwitchKeyword = 57371
const CaseKeyword = 57372
const VarKeyword = 57373
const Plus = 57374
const Minus = 57375
const Star = 57376
const Slash = 57377
const Percent = 57378
const Exclamation = 57379
const And = 57380
const Or = 57381
const Xor = 57382

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
	"FuncKeyword",
	"BracketOpen",
	"BracketClose",
	"BraceOpen",
	"BraceClose",
	"Comma",
	"Colon",
	"Semicolon",
	"Equal",
	"GreaterThan",
	"LessThan",
	"Dot",
	"Arrow",
	"ParenOpen",
	"ParenClose",
	"Eol",
	"IfKeyword",
	"ElseKeyword",
	"ReturnKeyword",
	"SwitchKeyword",
	"CaseKeyword",
	"VarKeyword",
	"Plus",
	"Minus",
	"Star",
	"Slash",
	"Percent",
	"Exclamation",
	"And",
	"Or",
	"Xor",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line logi.y:472

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 79,
	14, 81,
	25, 81,
	-2, 92,
	-1, 105,
	4, 30,
	5, 30,
	6, 30,
	7, 30,
	11, 30,
	13, 30,
	20, 30,
	23, 30,
	25, 30,
	-2, 95,
	-1, 106,
	4, 29,
	5, 29,
	6, 29,
	7, 29,
	11, 29,
	13, 29,
	20, 29,
	25, 29,
	-2, 98,
	-1, 107,
	4, 31,
	5, 31,
	6, 31,
	7, 31,
	11, 31,
	13, 31,
	20, 31,
	23, 31,
	25, 31,
	-2, 96,
	-1, 108,
	4, 32,
	5, 32,
	6, 32,
	7, 32,
	11, 32,
	13, 32,
	20, 32,
	23, 32,
	25, 32,
	-2, 97,
	-1, 111,
	14, 1,
	25, 1,
	-2, 4,
}

const yyPrivate = 57344

const yyLast = 344

var yyAct = [...]uint8{
	4, 49, 9, 12, 13, 42, 15, 163, 33, 101,
	104, 18, 19, 80, 66, 56, 35, 127, 32, 72,
	29, 60, 24, 188, 173, 192, 179, 53, 25, 190,
	7, 95, 58, 165, 164, 167, 166, 111, 14, 7,
	171, 61, 170, 151, 122, 121, 14, 8, 63, 172,
	59, 6, 14, 70, 14, 30, 91, 114, 115, 116,
	117, 118, 119, 120, 123, 124, 7, 161, 61, 68,
	140, 79, 109, 28, 156, 130, 92, 61, 14, 30,
	7, 94, 97, 89, 88, 84, 90, 14, 79, 89,
	88, 84, 90, 110, 142, 126, 128, 133, 14, 134,
	135, 23, 125, 93, 139, 14, 61, 159, 132, 7,
	95, 81, 131, 82, 14, 176, 83, 21, 155, 138,
	174, 145, 99, 148, 175, 55, 144, 146, 79, 14,
	143, 98, 150, 147, 54, 157, 41, 141, 26, 149,
	69, 162, 102, 168, 154, 152, 136, 160, 57, 137,
	158, 48, 129, 22, 68, 89, 88, 84, 90, 17,
	16, 45, 44, 43, 46, 1, 113, 20, 47, 145,
	52, 180, 181, 178, 153, 186, 177, 51, 182, 183,
	50, 169, 14, 168, 103, 38, 61, 189, 184, 187,
	185, 193, 26, 85, 191, 168, 87, 151, 122, 121,
	86, 194, 89, 88, 84, 90, 151, 122, 121, 73,
	78, 114, 115, 116, 117, 118, 119, 120, 123, 124,
	114, 115, 116, 117, 118, 119, 120, 123, 124, 112,
	122, 121, 45, 44, 77, 46, 89, 88, 84, 90,
	3, 76, 11, 114, 115, 116, 117, 118, 119, 120,
	123, 124, 107, 105, 106, 108, 75, 14, 81, 47,
	82, 52, 2, 83, 10, 74, 71, 65, 51, 39,
	27, 50, 100, 14, 81, 40, 82, 37, 96, 83,
	45, 44, 43, 46, 36, 34, 31, 47, 5, 52,
	0, 45, 44, 43, 46, 0, 51, 0, 47, 50,
	52, 62, 45, 44, 67, 46, 0, 51, 0, 47,
	50, 52, 7, 0, 0, 0, 0, 0, 51, 0,
	0, 50, 64, 45, 44, 43, 46, 0, 0, 0,
	47, 0, 52, 0, 165, 164, 167, 166, 0, 51,
	0, 171, 50, 170,
}

var yyPact = [...]int16{
	41, 41, 5, 5, 27, 5, 154, -1000, 153, 27,
	5, 5, 27, 27, -1000, 104, 78, -1000, 27, 27,
	-1000, 5, 125, 49, 157, -1000, 5, 110, -1000, -1000,
	142, 287, 276, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 5, -1000, -1000,
	298, 129, 5, 232, -1000, 5, -1000, 56, 89, 276,
	6, -1000, -1000, 157, -1000, 107, -1000, 142, 319, 136,
	248, 12, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	211, 79, 198, 146, 52, -1000, -1000, -1000, -1000, -1000,
	-1000, 73, 142, -1000, 6, -1000, 84, 319, -1000, 5,
	134, -1000, 228, 55, -1000, 121, 52, -1000, -1000, 80,
	85, -1000, 151, 198, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 115, 105, -1000, -1000, 198, 179, -1000, 188, 142,
	198, -1000, 99, 62, 5, 157, 88, 136, -1000, 53,
	5, 330, -1000, -1000, 188, -1000, 188, -1000, -1000, 25,
	-3, 103, 102, 100, 188, -1000, -1000, 157, -1000, -1000,
	-1000, -1000, 21, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	5, 5, 125, 125, 198, -1000, 198, 319, -1000, 121,
	21, 29, -4, -1000, 188, 188, 14, -1000, 125, 13,
	5, -1000, -1000, 29, -1000,
}

var yyPgo = [...]int16{
	0, 15, 262, 288, 151, 286, 18, 8, 285, 284,
	278, 277, 16, 275, 272, 9, 136, 270, 20, 269,
	267, 14, 5, 266, 19, 265, 256, 241, 234, 210,
	209, 13, 200, 196, 193, 17, 185, 1, 184, 10,
	181, 7, 175, 174, 240, 166, 165, 0, 21,
}

var yyR1 = [...]int8{
	0, 47, 47, 47, 48, 48, 46, 46, 46, 46,
	46, 46, 46, 2, 3, 4, 5, 5, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 8,
	12, 12, 12, 9, 10, 10, 10, 11, 36, 37,
	38, 38, 38, 39, 41, 41, 41, 41, 41, 41,
	40, 42, 42, 42, 13, 14, 14, 15, 15, 16,
	16, 17, 17, 18, 19, 19, 20, 20, 21, 1,
	1, 22, 23, 23, 23, 24, 24, 24, 24, 24,
	24, 30, 25, 26, 27, 27, 27, 27, 28, 29,
	29, 31, 31, 31, 31, 32, 32, 32, 33, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 34, 35, 43, 43, 43, 44,
}

var yyR2 = [...]int8{
	0, 1, 2, 0, 1, 2, 2, 2, 2, 1,
	3, 3, 0, 3, 2, 5, 2, 3, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 5, 1, 4, 0, 1, 1, 5,
	1, 4, 0, 3, 1, 1, 1, 1, 1, 1,
	5, 1, 4, 0, 5, 1, 3, 1, 2, 3,
	2, 1, 4, 2, 3, 2, 1, 4, 1, 1,
	4, 5, 1, 3, 0, 1, 1, 1, 1, 1,
	1, 1, 3, 1, 5, 7, 3, 5, 2, 3,
	5, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 3, 4, 1, 3, 0, 4,
}

var yyChk = [...]int16{
	-1000, -46, -2, -44, -47, -3, 10, 25, 6, -47,
	-2, -44, -47, -47, 25, -47, 6, 6, -47, -47,
	-4, 13, -16, 23, -47, -22, 13, -17, 24, -18,
	6, -5, -6, -7, -8, -12, -9, -11, -36, -19,
	-13, -16, -22, 6, 5, 4, 7, 11, -4, -37,
	23, 20, 13, -47, 24, 15, -1, 6, -47, -6,
	-48, -7, 25, -47, 24, -20, -21, 6, -6, 11,
	-47, -23, -24, -30, -25, -26, -27, -28, -29, -35,
	-31, 26, 28, 31, 6, -34, -32, -33, 5, 4,
	7, -47, 20, 14, -48, 25, -10, -6, 24, 15,
	-14, -15, 6, -38, -39, 5, 6, 4, 7, -47,
	-48, 25, 18, -45, 32, 33, 34, 35, 36, 37,
	38, 20, 19, 39, 40, 23, -31, -35, -31, 6,
	23, -18, -1, -47, 15, -47, 12, 15, -12, -47,
	15, 16, 14, -24, -31, 18, -31, 18, 18, -31,
	-22, 18, -1, -43, -31, 19, 12, -47, -21, 19,
	-15, 14, -47, -41, 5, 4, 7, 6, -37, -40,
	13, 11, 24, 27, 18, 24, 15, -6, -39, 5,
	-47, -47, -22, -22, -31, -31, -42, -41, 27, -47,
	15, -22, 12, -47, -41,
}

var yyDef = [...]int8{
	3, -2, 3, 3, 9, 3, 0, 1, 0, 8,
	3, 3, 6, 7, 2, 0, 0, 14, 10, 11,
	13, 3, 0, 0, 0, 118, 3, 0, 60, 61,
	0, 3, 0, 18, 20, 21, 22, 23, 24, 25,
	26, 27, 28, 29, 30, 31, 32, 3, 37, 38,
	0, 0, 3, 74, 59, 3, 63, 69, 0, 0,
	16, 19, 4, 36, 60, 0, 66, 29, 68, 0,
	42, 3, 72, 75, 76, 77, 78, 79, 80, -2,
	83, 0, 0, 0, 98, 91, 93, 94, 95, 96,
	97, 0, 0, 15, 17, 5, 3, 34, 64, 3,
	0, 55, 57, 3, 40, -2, -2, -2, -2, 0,
	0, -2, 0, 0, 99, 100, 101, 102, 103, 104,
	105, 106, 107, 108, 109, 0, 0, 92, 88, 0,
	117, 62, 0, 0, 3, 0, 0, 0, 58, 0,
	3, 0, 71, 73, 82, 110, 113, 112, 111, 0,
	86, 0, 89, 0, 115, 70, 33, 0, 67, 54,
	56, 39, 0, 43, 44, 45, 46, 47, 48, 49,
	3, 3, 0, 0, 0, 114, 0, 35, 41, 0,
	42, 53, 84, 87, 90, 116, 3, 51, 0, 0,
	3, 85, 50, 0, 52,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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
//line logi.y:57
		{
			registerRootNode(yylex, yyDollar[1].node)
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:60
		{
			registerRootNode(yylex, yyDollar[1].node)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:65
		{
			registerRootNode(yylex, yyDollar[2].node)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:68
		{
			registerRootNode(yylex, yyDollar[2].node)
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:74
		{
			yyVAL.node = appendNode(NodeOpDefinition, yyDollar[1].node, yyDollar[3].node)
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:79
		{
			yyVAL.node = appendNode(NodeOpSignature, newNode(NodeOpMacro, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location), newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location))
		}
	case 15:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:86
		{
			yyVAL.node = appendNode(NodeOpBody, yyDollar[3].node)
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:91
		{
			yyVAL.node = appendNode(NodeOpStatements, yyDollar[1].node)
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:95
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:100
		{
			yyVAL.node = appendNode(NodeOpStatement, yyDollar[1].node)
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:104
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[2].node)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:111
		{
			yyVAL.node = newNode(NodeOpIdentifier, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:116
		{
			yyVAL.node = newNode(NodeOpValue, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:120
		{
			yyVAL.node = newNode(NodeOpValue, yyDollar[1].number, yyDollar[1].token, yyDollar[1].location)
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:124
		{
			yyVAL.node = newNode(NodeOpValue, yyDollar[1].bool, yyDollar[1].token, yyDollar[1].location)
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:129
		{
			yyVAL.node = yyDollar[3].node
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:134
		{
			yyVAL.node = appendNode(NodeOpArray, yyDollar[1].node)
		}
	case 35:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:138
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 36:
		yyDollar = yyS[yypt-0 : yypt+1]
//line logi.y:142
		{
			yyVAL.node = appendNode(NodeOpArray)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:147
		{
			yyVAL.node = appendNode(NodeOpStruct, yyDollar[1].node)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:152
		{
			yyVAL.node = yyDollar[1].node
		}
	case 39:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:157
		{
			yyVAL.node = yyDollar[3].node
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:162
		{
			yyVAL.node = appendNode(NodeOpJsonObject, yyDollar[1].node)
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:165
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 42:
		yyDollar = yyS[yypt-0 : yypt+1]
//line logi.y:168
		{
			yyVAL.node = appendNode(NodeOpJsonObject)
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:173
		{
			yyVAL.node = newNode(NodeOpJsonObjectItem, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[3].node)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:178
		{
			yyVAL.node = newNode(NodeOpJsonObjectItemValue, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:181
		{
			yyVAL.node = newNode(NodeOpJsonObjectItemValue, yyDollar[1].number, yyDollar[1].token, yyDollar[1].location)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:184
		{
			yyVAL.node = newNode(NodeOpJsonObjectItemValue, yyDollar[1].bool, yyDollar[1].token, yyDollar[1].location)
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:187
		{
			yyVAL.node = newNode(NodeOpJsonIdentifier, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:190
		{
			yyVAL.node = yyDollar[1].node
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:193
		{
			yyVAL.node = yyDollar[1].node
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:198
		{
			yyVAL.node = yyDollar[3].node
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:203
		{
			yyVAL.node = appendNode(NodeOpJsonArray, yyDollar[1].node)
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:206
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 53:
		yyDollar = yyS[yypt-0 : yypt+1]
//line logi.y:209
		{
			yyVAL.node = appendNode(NodeOpJsonArray)
		}
	case 54:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:215
		{
			yyVAL.node = yyDollar[3].node
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:220
		{
			yyVAL.node = appendNode(NodeOpAttributeList, yyDollar[1].node)
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:224
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:229
		{
			yyVAL.node = newNode(NodeOpAttribute, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:233
		{
			yyVAL.node = newNode(NodeOpAttribute, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:238
		{
			yyVAL.node = yyDollar[2].node
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:242
		{
			yyVAL.node = appendNode(NodeOpArgumentList)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:247
		{
			yyVAL.node = appendNode(NodeOpArgumentList, yyDollar[1].node)
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:251
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:256
		{
			yyVAL.node = newNode(NodeOpArgument, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[2].node)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:261
		{
			yyVAL.node = yyDollar[2].node
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:265
		{
			yyVAL.node = appendNode(NodeOpParameterList)
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:270
		{
			yyVAL.node = appendNode(NodeOpParameterList, yyDollar[1].node)
		}
	case 67:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:274
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[4].node)
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:279
		{
			yyVAL.node = newNode(NodeOpParameter, yyDollar[1].node, yyDollar[1].token, yyDollar[1].location)
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:285
		{
			yyVAL.node = newNode(NodeOpTypeDef, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 70:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:289
		{
			yyVAL.node = newNode(NodeOpTypeDef, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[3].node)
		}
	case 71:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:297
		{
			yyVAL.node = appendNode(NodeOpCodeBlock, yyDollar[3].node)
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:302
		{
			yyVAL.node = appendNode(NodeOpStatements, yyDollar[1].node)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:306
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 74:
		yyDollar = yyS[yypt-0 : yypt+1]
//line logi.y:310
		{
			yyVAL.node = appendNode(NodeOpStatements)
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:319
		{
			yyVAL.node = appendNode(NodeOpFunctionCall, yyDollar[1].node)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:324
		{
			yyVAL.node = appendNode(NodeOpAssignment, yyDollar[1].node, yyDollar[3].node)
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:329
		{
			yyVAL.node = appendNode(NodeOpExpression, yyDollar[1].node)
		}
	case 84:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:334
		{
			yyVAL.node = appendNode(NodeOpIf, yyDollar[3].node, yyDollar[5].node)
		}
	case 85:
		yyDollar = yyS[yypt-7 : yypt+1]
//line logi.y:338
		{
			yyVAL.node = appendNode(NodeOpIfElse, yyDollar[3].node, yyDollar[5].node, yyDollar[7].node)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:342
		{
			yyVAL.node = appendNode(NodeOpIf, yyDollar[2].node, yyDollar[3].node)
		}
	case 87:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:346
		{
			yyVAL.node = appendNode(NodeOpIfElse, yyDollar[2].node, yyDollar[3].node, yyDollar[5].node)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:351
		{
			yyVAL.node = appendNode(NodeOpReturn, yyDollar[2].node)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:356
		{
			yyVAL.node = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location), yyDollar[3].node)
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
//line logi.y:360
		{
			yyVAL.node = appendNode(NodeOpVariableDeclaration, newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location), yyDollar[3].node, yyDollar[5].node)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:369
		{
			yyVAL.node = newNode(NodeOpLiteral, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:373
		{
			yyVAL.node = newNode(NodeOpLiteral, yyDollar[1].number, yyDollar[1].token, yyDollar[1].location)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:377
		{
			yyVAL.node = newNode(NodeOpLiteral, yyDollar[1].bool, yyDollar[1].token, yyDollar[1].location)
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:382
		{
			yyVAL.node = newNode(NodeOpVariable, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:386
		{
			yyVAL.string = "+"
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:389
		{
			yyVAL.string = "-"
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:393
		{
			yyVAL.string = "*"
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:397
		{
			yyVAL.string = "/"
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:401
		{
			yyVAL.string = "%"
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:405
		{
			yyVAL.string = "!"
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:409
		{
			yyVAL.string = "&&"
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:413
		{
			yyVAL.string = "<"
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:417
		{
			yyVAL.string = ">"
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:421
		{
			yyVAL.string = "||"
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:425
		{
			yyVAL.string = "^"
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:429
		{
			yyVAL.string = "=="
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:433
		{
			yyVAL.string = ">="
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
//line logi.y:437
		{
			yyVAL.string = "<="
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:442
		{
			yyVAL.node = appendNode(NodeOpBinaryExpression, yyDollar[1].node, yyDollar[3].node, newNode(NodeOpOperator, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location))
		}
	case 114:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:447
		{
			yyVAL.node = newNode(NodeOpFunctionCall, yyDollar[1].string, yyDollar[1].token, yyDollar[1].location, yyDollar[3].node)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
//line logi.y:452
		{
			yyVAL.node = appendNode(NodeOpFunctionParams, yyDollar[1].node)
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
//line logi.y:456
		{
			yyVAL.node = appendNodeTo(&yyDollar[1].node, yyDollar[3].node)
		}
	case 117:
		yyDollar = yyS[yypt-0 : yypt+1]
//line logi.y:460
		{
			yyVAL.node = appendNode(NodeOpFunctionParams)
		}
	case 118:
		yyDollar = yyS[yypt-4 : yypt+1]
//line logi.y:468
		{
			yyVAL.node = appendNode(NodeOpFunctionDefinition, newNode(NodeOpName, yyDollar[2].string, yyDollar[2].token, yyDollar[2].location), yyDollar[3].node, yyDollar[4].node)
		}
	}
	goto yystack /* stack new state and value */
}
