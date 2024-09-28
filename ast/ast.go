package ast

import (
	"bytes"

	"local.packages/token"
)

// Node ノード
// ASTのノードは、すべてNodeインターフェースを実装する必要がある、すなわち
//  `TokenLiteral` メソッドを提供し、そのNodeが関連付けられているトークンのリテラル値を返す必要がある。
type Node interface {
	TokenLiteral() string
	// デバッグのためにNodeごとに固有の文字列を出力する。
	String() string
}

// Statement 文
// 文は値を返さない。
type Statement interface {
	Node
	statementNode() // 式と文を間違えていたらコンパイラが教えてくれる
}

// Expression 式
// 式は値を返す。
type Expression interface {
	Node
	expressionNode() // 式と文を間違えていたらコンパイラが教えてくれる
}

// Program すべてのASTのルートノード
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

// Program#stringの実装
// 各Statementの文字列出力を一つの文字列にまとめて返す。
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement let文
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // 束縛の識別子を保持する
	Value Expression  // 値を保持する式を保持する
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// LetStatementの文字列出力を返す。
// `LET name = value;` 形式の文字列を返す。
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier 識別子
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Identifierの文字列出力を返す。
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement return文
type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// ReturnStatementの文字列出力を返す。
// `RETURN value;` 形式の文字列を返す。
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement 式文
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// ExpressionStatementの文字列出力を返す。
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral 整数リテラル
type IntegerLiteral struct {
	Token token.Token
	// ソースコード中の整数リテラルが表現する値を格納するためのフィールド
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// 前置
type PrefixExpression struct {
	Token token.Token // 前置トークン、例えば「!」
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
