package token

// Monkey言語の字句を表すトークンを格納するモジュール。
// LET、RETURN等、Monkey言語で定義されるトークンが存在し、それぞれのトークンは固有の文字列表現を持つ。

// TokenType トークンの種別を表すstringの別名
type TokenType string

// Token トークンを表す構造体
type Token struct {
	Type    TokenType
	Literal string
}

// 予約語とそのTokenTypeへのマッピング
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// 定数定義のブロック
const (
	// ILLEGAL 未知のトークン
	ILLEGAL = "ILLEGAL"

	// EOF ファイル末端（End of File）を表すトークン
	EOF = "EOF"

	// IDENT 識別子を表すトークン：add、foobar、x、y、...
	IDENT = "IDENT"

	// INT 整数リテラルを表すトークン：1234567...
	INT = "INT"

	// STRING 文字列リテラルを表すトークン
	STRING = "STRING"

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// LookupIdentifier 識別子が予約語にマッチしたら予約語に対応するTokenTypeを、
// マッチしなかったらIDENTを返す。
func LookupIdentifier(identifier string) TokenType {
	if t, ok := keywords[identifier]; ok {
		return t
	}
	return IDENT
}
