package lexer

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
