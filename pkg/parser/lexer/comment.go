package lexer

func (sc *lexer) handleComments(r rune) {
	// single line comments
	if sc.read() == '/' {
		for {
			r = sc.read()
			if isEol(r) || r == 0 {
				sc.unread()
				return
			}
		}
	} else {
		sc.unread()
	}

	// multi line comments
	if sc.read() == '*' {
		for {
			r = sc.read()
			if r == '*' {
				if sc.read() == '/' {
					sc.unread()
					return
				}
			}
		}
	} else {
		sc.unread()
	}
}
