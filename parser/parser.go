package parser

import (
	"fmt"
	"strconv"

	"local.packages/ast"
	"local.packages/lexer"
	"local.packages/token"
)

// Parser 構文解析器
type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	// `curToken.Type` に関連付けられた構文解析関数が前置かどうかをチェックするためのマップ
	prefixParseFns map[token.TokenType]prefixParseFn
	// `curToken.Type` に関連付けられた構文解析関数が中置かどうかをチェックするためのマップ
	infixParseFns map[token.TokenType]infixParseFn
}

// 優先順位の定義
const (
	_ int = iota // 次にくる定数にインクリメントしながら数を与えるための定義
	LOWEST       // 最も低い優先順位
	EQUALS       // ==
	LESSGREATER  // >, <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X, !X
	CALL         // myFunction(X
)

// New Parserを生成する。
// Lexerを受け取り、トークンを読み込むことでParserが初期化される。
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// `prefixParseMap` を初期化し、構文解析関数を登録する。
	// `token.IDENT` が出現したら `parseIdentifier` を呼び出す、等の登録を行なっている。
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// トークンを2つ読み込む。curTokenとpeekTokenがセットされる。
	p.nextToken()
	p.nextToken()

	return p
}

// 次のトークンを読み込む。
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram ファイル末尾に達するまでStatementを読み込み、読み込んだStatementを持つProgramを返す。
func (p *Parser) ParseProgram() *ast.Program {
	// ASTのルートノードの生成
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// ファイル末尾に至るまで、繰り返しトークンを読み込む。
	for p.curToken.Type != token.EOF {
		// 文を構文解析する。
		stmt := p.parseStatement()
		if stmt != nil {
			// parseStatementが文を返した場合は、その戻り値をStatementsに追加する。
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// 構文解析する対象が尽きたら、ASTのルートノードを返す。
	return program
}

// Statementを構築して返す。
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// `let` 文を解析して、LetStatementノードを構築して返す。
// 現在解析しているLETトークンに基づいて、LetStatementノードを構築している。
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// 後続するトークンに対するアサーション
	// 初期状態では後続するトークンに識別子 `IDENT` を期待している。
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// 識別子 `Identifier` を構築する。
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 次に等号 `ASSIGN` を期待している。
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// セミコロン `SEMICOLON` に遭遇するまで読み飛ばしている。
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ReturnStatementを構築して返す。
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 式文を解析する。
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 式を解析して返す。
// 前置に関連付けられた構文解析関数を呼び出し、その結果を返す。
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	// 前置に関連付けられたトークンがなければ `nil` を返す。
	if prefix == nil {
		return nil
	}

	// 前置に関連付けられたトークンがあれば解析して返す。
	leftEx := prefix()
	return leftEx
}

// 識別子を解析して返す。
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 整数リテラルを解析して返す。
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// 整数リテラルをint64として解釈する。
	// 解釈に失敗した場合はエラーを返し、成功した場合は `lit.Value` にint64を詰めて返す。
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
}

// 現在のトークンがtと等しい時にtrueを返す。
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次のトークンがtと等しい時にtrueを返す。
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 次のトークンがtと等しい時に次のトークンを読み込み、trueを返す。
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// 前置構文を `prefixParseFns` に登録する。
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// 中置構文を `prefixParseFns` に登録する。
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

type (
	// 前置構文解析関数
	prefixParseFn func() ast.Expression
	// 中置構文解析関数
	// 引数 `expression` は、中置演算子の「左側」に置かれる式。
	infixParseFn func(expression ast.Expression) ast.Expression
)

// Errors エラーの文字列のスライスを返す。
func (p *Parser) Errors() []string {
	return p.errors
}

// peekTokenが期待されたものでない場合にエラーのスライスに追加する。
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
