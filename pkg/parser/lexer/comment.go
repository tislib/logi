package lexer

func (sc *lexer) handleComments(r rune) {
	nc, _ := sc.peekChar()

	// single line comments
	if nc == '/' {
		l := sc.peekUntil(func(ch rune) bool {
			return ch == '\n'
		})

		sc.discard(len(l))
		return
	}

	// multi line comments
	if nc == '*' {
		for {
			r = sc.read()
			if r == '*' {
				nc, _ := sc.peekChar()
				if nc == '/' {
					sc.discardChar()
					return
				}
			}
		}
	}
}
