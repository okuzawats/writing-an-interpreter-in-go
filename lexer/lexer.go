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
		//
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 次の文字からtoken.Tokenを生成して返す。
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	switch l.ch {
	case '=':
		t = newToken(token.ASSIGN, l.ch)
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	}

	l.readChar()
	return t
}

// token.Tokenを生成して返す。
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
