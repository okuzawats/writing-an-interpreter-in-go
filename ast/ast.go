package ast

import "local.packages/token"

// ノード
type Node interface {
	TokenLiteral() string
}

// 文
type Statement interface {
	Node
	statementNode() // 式と文を間違えていたらコンパイラが教えてくれる
}

// 式
type Expression interface {
	Node
	expressionNode() // 式と文を間違えていたらコンパイラが教えてくれる
}

// すべてのASTのルートノード
type Program struct {
	// Monkeyプログラムの文の集まりが格納される。
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// let文
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // 束縛の識別子を保持する
	Value Expression  // 値を保持する式を保持する
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// 識別子
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
