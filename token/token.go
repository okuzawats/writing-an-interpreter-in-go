package token

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

	// EOF End of File
	EOF = "EOF"

	// IDENT 識別子
	IDENT = "IDENT"

	// INT リテラル
	INT = "INT"

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

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

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
