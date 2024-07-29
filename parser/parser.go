package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParsefn func() ast.Expression // 前置構文解析関数
	infixParseFn func(ast.Expression) ast.Expression // 中値構文解析関数
)

// トークンを見ていってASTを作成する
type Parser struct {
	l *lexer.Lexer
	errors []string 
	curToken token.Token // 今見ているトークン
	peekToken token.Token // 次のトークン
	prefixParsefns map[token.TokenType]prefixParsefn
	infixParseFns map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{ l: l, errors: []string{} }
	p.prefixParsefns = make(map[token.TokenType]prefixParsefn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{ Token: p.curToken, Value: p.curToken.Literal }
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tt token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tt, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // ASTのルートノード
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET: // let ときているとき
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{ Token: p.curToken }
	if !p.expectPeek(token.IDENT) { // let <identifier> ときているかチェック
		return nil
	}

	stmt.Name = &ast.Identifier{ Token: p.curToken, Value: p.curToken.Literal }
	if !p.expectPeek(token.ASSIGN) { // let <identifier>  = ときているかチェック
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) { // とりあえず、=まで来たらおけ
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{ Token: p.curToken }
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression {
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(token.PREFIX)
	return expression
}

// 式の構文解析
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{ Token: p.curToken }
	stmt.Expression = p.parseExpression(token.LOWEST) // 優先順位
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParsefns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()


    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }
        p.nextToken()

		// prefixParseFnからの式を渡している
        leftExp = infix(leftExp)
    }
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{ Token: p.curToken }
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) curTokenIs(tt token.TokenType) bool {
	return p.curToken.Type == tt
}

func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) expectPeek(tt token.TokenType) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	} else {
		p.peekError(tt)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParsefn) {
	p.prefixParsefns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

var precedences = map[token.TokenType]int{
    token.EQ:       token.EQUALS,
    token.NOT_EQ:   token.EQUALS,
    token.LT:       token.LESSGREATER,
    token.GT:       token.LESSGREATER,
    token.PLUS:     token.SUM,
    token.MINUS:    token.SUM,
    token.SLASH:    token.PRODUCT,
    token.ASTERISK: token.PRODUCT,
}

func (p *Parser) peekPrecedence() int {
    if precedence, ok := precedences[p.peekToken.Type]; ok { return precedence }
    return token.LOWEST
}

func (p *Parser) curPrecedence() int {
    if p, ok := precedences[p.curToken.Type]; ok { return p }
    return token.LOWEST
}

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
