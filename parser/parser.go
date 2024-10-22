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
	_           int = iota // 次にくる定数にインクリメントしながら数を与えるための定義
	LOWEST                 // 最も低い優先順位
	EQUALS                 // ==
	LESSGREATER            // >, <
	SUM                    // +
	PRODUCT                // *
	PREFIX                 // -X, !X
	CALL                   // myFunction(X
)

// 優先順位テーブル
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// 現在位置の次の位置のトークンの優先順位を返す。
func (p *Parser) peekPrecedence() int {
	// 優先順位テーブルに対象のトークンが見つかれば、見つかった優先順位を返す。
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	// 見つからなかったら、最も低い優先順位を返す。
	return LOWEST
}

// 現在位置のトークンの優先順位を返す。
func (p *Parser) curPrecedence() int {
	// 優先順位テーブルに対象のトークンが見つかれば、見つかった優先順位を返す。
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	// 見つからなかったら、最も低い優先順位を返す。
	return LOWEST
}

// New Parserを生成する。
// Lexerを受け取り、トークンを読み込むことでParserが初期化される。
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// `prefixParseFns` を初期化し、前置演算子の構文解析関数を登録する。
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// `infixParseFns` を初期化し、中置演算子の構文解析関数を登録する。
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ReturnStatementを構築して返す。
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 式を解析して返す。
// 前置に関連付けられた構文解析関数を呼び出し、その結果を返す。
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	// 前置に関連付けられたトークンがなければ `nil` を返す。
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// 前置に関連付けられたトークンがあれば解析して返す。
	leftExp := prefix()

	// 次のトークンに紐つけられている `infixParseFn` を探し、その返り値を渡す。
	// これをより低い優先順位のトークンに遭遇するまで繰り返し実行する。
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// 前置式を解析して返す。
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// 前置式に対応する式を読み込み、Rightに詰める。
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 真偽値を解析して返す。
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// 丸括弧 `()` 内の式を指揮を解析して返す。
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatements()

	// else句がある場合はそのブロックを解析する
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatements()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatements()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// 引数がない場合
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// 最初の引数
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// カンマで区切られたそれ以降の引数
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// 括弧が閉じられていない場合は `nil` を返す
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatements() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
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

// 中置式を解析し、Expressionノードを返す。
// Leftは引数として受け取り、構文を解析してRightを取り出してExpressionに紐つけている。
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// 関数呼び出しに引数がない場合
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	// カンマ区切りの引数の解析
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	// 括弧が閉じられていない場合
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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
