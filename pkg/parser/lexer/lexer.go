package lexer

import (
	"bufio"
	"strconv"
	"strings"
)

type lexer struct {
	config  LexerConfig
	buf     *bufio.Reader
	debug   bool
	readStr string
}

func (s *lexer) GetReadString() any {
	return s.readStr
}

func (s *lexer) Next() (token Token, err error) {
	for {
		r := s.read()
		if r == 0 { // EOF
			return Token{}, ErrEOF
		}

		// handle comments
		if s.config.HandleComments && r == '/' {
			s.handleComments(r)
			continue
		}

		for _, tokenConfig := range s.config.Tokens {
			token, matched := s.matchToken(tokenConfig, r)

			if matched {
				return *token, nil
			}
		}

		if isWhitespace(r) {
			continue
		}

		return Token{}, ErrInvalidToken
	}
}

func (s *lexer) matchToken(config TokenConfig, startingChar rune) (*Token, bool) {
	// StartsWith   string
	//	EndsWith     string
	//	Equals       string
	//	IsWhiteSpace bool
	//	IsAlphaNum   bool
	//	IsAlpha      bool
	//	IsDigit      bool
	//	IsEol        bool

	var isSingleChar = config.IsEol
	var needsToLookAhead = !isSingleChar && (config.StartsWith != "" || config.EndsWith != "" || config.Equals != "")
	var isLengthKnown = config.Equals != "" || config.EqualsCaseInsensitive != "" || isSingleChar || len(config.EqualOneOf) > 0

	if isSingleChar {
		if config.IsEol {
			if isEol(startingChar) {
				return &Token{Id: config.Id, Value: string(startingChar)}, true
			}
		}

		return nil, false
	}

	if isLengthKnown {
		if config.Equals != "" {
			data, _ := s.buf.Peek(len(config.Equals) - 1)
			var value = string(startingChar) + string(data)

			if value == config.Equals {
				s.buf.Discard(len(data))
				return &Token{Id: config.Id, Value: value}, true
			}
		}

		if config.EqualsCaseInsensitive != "" {
			data, _ := s.buf.Peek(len(config.EqualsCaseInsensitive) - 1)
			var value = string(startingChar) + string(data)

			if strings.ToLower(value) == strings.ToLower(config.EqualsCaseInsensitive) {
				s.buf.Discard(len(data))
				return &Token{Id: config.Id, Value: value}, true
			}
		}

		if len(config.EqualOneOf) > 0 {
			for _, equal := range config.EqualOneOf {
				data, _ := s.buf.Peek(len(equal) - 1)
				var value = string(startingChar) + string(data)

				if value == equal {
					s.buf.Discard(len(data))
					return &Token{Id: config.Id, Value: value}, true
				}
			}
		}

		return nil, false
	}

	if needsToLookAhead {
		if config.StartsWith != "" {
			var length = len(config.StartsWith)
			var startConditionMatched bool
			if len(config.StartsWith) == 1 && startingChar == rune(config.StartsWith[0]) {
				startConditionMatched = true
			} else {
				data, _ := s.buf.Peek(len(config.StartsWith) - 1)
				var startValue = string(startingChar) + string(data)

				if startValue == config.StartsWith {
					startConditionMatched = true
				}
			}

			if startConditionMatched {
				var endConditionMatched bool
				if config.EndsWith != "" {
					var i = 0
					for {
						// peek and check without read

						data, _ := s.buf.Peek(i)

						if len(data) == 0 {
							break
						}

						if len(data) == 1 && data[0] == config.EndsWith[0] {
							endConditionMatched = true
							length = i + 1
							break
						}

						if len(data) > 1 {
							if data[len(data)-1] == config.EndsWith[0] {
								endConditionMatched = true
								length = i + 1
								break
							}
						}

						i++
					}
				}

				if endConditionMatched {
					data, _ := s.buf.Peek(length - 1)
					var value = string(startingChar) + string(data)
					s.buf.Discard(len(data))
					return &Token{Id: config.Id, Value: value}, true
				}
			}
		}

		return nil, false
	}

	if config.IsAlpha {
		if isAlpha(startingChar) {
			var value = string(startingChar)
			for {
				r := s.read()
				if !isAlphaNum(r) {
					s.unread()
					break
				}

				value += string(r)
			}

			return &Token{Id: config.Id, Value: value}, true
		}
	}

	if config.IsAlphaNum {
		if isAlphaNum(startingChar) {
			var value = string(startingChar)
			for {
				r := s.read()
				if !isAlphaNum(r) {
					s.unread()
					break
				}

				value += string(r)
			}

			return &Token{Id: config.Id, Value: value}, true
		}
	}

	if config.IsIdentifier {
		if isAlpha(startingChar) {
			var value = string(startingChar)
			for {
				r := s.read()
				if !isAlphaNum(r) {
					s.unread()
					break
				}

				value += string(r)
			}

			return &Token{Id: config.Id, Value: value}, true
		}
	}

	if config.IsDigit {
		if isDigit(startingChar) || startingChar == '-' {
			var value = string(startingChar) + s.peekUntil(func(ch rune) bool {
				return !isDigit(ch) && ch != '.' && !isAlpha(ch)
			})

			_, err := strconv.ParseFloat(value, 64)

			if err == nil {
				s.buf.Discard(len(value) - 1)

				return &Token{Id: config.Id, Value: value}, true
			}

			return nil, false
		}

		return nil, false
	}

	if config.IsString {
		if startingChar == '"' {
			var data = ""

			for {
				r := s.read()
				if r == '"' || r == 0 {
					break
				}
				data += string(r)
			}

			return &Token{Id: config.Id, Value: data}, true
		}
	}

	return nil, false
}

func (s *lexer) read() rune {
	ch, _, _ := s.buf.ReadRune()
	s.readStr += string(ch)
	return ch
}

func (s *lexer) unread() {
	_ = s.buf.UnreadRune()
	s.readStr = s.readStr[:len(s.readStr)-1]
}

func (s *lexer) peekUntil(endFunc func(ch rune) bool) string {
	var o = 0
	var l = 16

	for {
		data, err := s.buf.Peek(l)

		if len(data) < l {
			l = len(data)
		}

		for i := o; i < l; i++ {
			r := rune(data[i])
			if endFunc(r) {
				l = i

				return string(data[:l])
			}
		}

		o += l
		l += l

		if l > 1024 {
			return ""
		}

		if err != nil {
			return string(data)
		}
	}

	return ""
}
