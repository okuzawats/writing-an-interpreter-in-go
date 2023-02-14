package parser

import (
	"fmt"

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

	prefixParseFn map[token.TokenType]prefixParseFn
	infixParseFn  map[token.TokenType]infixParseFn
}

// New Parserを生成する。
// Lexerを受け取り、トークンを読み込むことでParserが初期化される。
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

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
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

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
		return nil
	}
}

// LetStatementを構築して返す。
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

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

type (
	// 前置構文解析関数
	prefixParseFn func() ast.Expression
	// 中置構文解析関数
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
