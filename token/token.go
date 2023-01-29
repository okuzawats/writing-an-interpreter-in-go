package token

// トークンの種別を表すstringの別名
type TokenType string

// トークンを表す構造体
type Token struct {
	Type    TokenType
	Literal string
}

// 定数定義のブロック
const (
	// 未知のトークン
	ILLEGAL = "ILLEGAL"

	// End of File
	EOF = "EOF"

	// 識別子＋リテラル
	IDENT = "IDENT"
	INT   = "INT"

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
