package lexer

import "local.packages/token"

// 字句
type Lexer struct {
	// 入力値
	input string
	// 現在の文字の位置
	position int
	// これから読み込む位置（現在の文字の次）
	readPosition int
	// 現在検査中の文字
	ch byte
}

// Lexerを生成して返す。
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 次の文字を読んで、入力値の現在位置を進める。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 末端に到達した場合。
		// ASCIIコードの"NUL"文字に対応している。
		l.ch = 0
	} else {
		// それ以外の場合は、その位置にある文字を読み取る。
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 次の文字からtoken.Tokenを生成して返す。
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// "=="の場合
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.EQ, Literal: literal}
		} else {
			// "=" の場合
			t = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			// "!="の場合
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			// "!" の場合
			t = newToken(token.BANG, l.ch)
		}
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.MINUS, l.ch)
	case '/':
		t = newToken(token.SLASH, l.ch)
	case '*':
		t = newToken(token.ASTERISK, l.ch)
	case '<':
		t = newToken(token.LT, l.ch)
	case '>':
		t = newToken(token.GT, l.ch)
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			// 識別子の場合
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t
		} else if isDigit(l.ch) {
			// 整数リテラルの場合
			t.Type = token.INT
			t.Literal = l.readNumber()
			return t
		} else {
			// 不明なトークンの場合
			t = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return t
}

// token.Tokenを生成して返す。
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 現在位置の次の位置の文字を返す。
// positionは進めない。
// また、現在位置が末尾の時は0を返す。
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// 連続する文字を識別子として取り出して文字列として返す。
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// a-zA-z_にマッチする場合にtrueを返す。
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 連続する数字を取り出して文字列として返す。
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 0-9にマッチする場合にtrueを返す。
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 空白文字を読み飛ばす。
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
