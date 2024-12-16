package lexer

func isWhitespace(ch rune) bool    { return ch == ' ' || ch == '\t' || ch == '\n' }
func isEol(ch rune) bool           { return ch == '\n' || ch == '\r' }
func isDigit(r rune) bool          { return r >= '0' && r <= '9' }
func isAlpha(r rune) bool          { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') }
func isAlphaNum(r rune) bool       { return isAlpha(r) || isDigit(r) }
func isIdentifierChar(r rune) bool { return isAlpha(r) || isDigit(r) || r == '_' || r == '$' }
