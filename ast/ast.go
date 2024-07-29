package ast

import (
	"bytes"
	"monkey/token"
)

/*
let <identifier> = <expression>; <- これで一つのstatement

ASTはNodeだけで構成される
Programノードは全てのASTのルートノード
*/

type Node interface {
	TokenLiteral() string // ノードが関連づけられているトークンのリテラル値を返す、デバッグとテストのために使う
	String() string
}

// 文
type Statement interface { 
	Node
	statementNode() //コンパイラに情報を与えるために存在
}

// 式
type Expression interface { 
	Node
	expressionNode() //コンパイラに情報を与えるために存在
}

type Program struct {
	Statements []Statement
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct { 
	Token token.Token // let
	Name *Identifier // 識別子 ex) x, y, ...
	Value Expression // 値 ex) 5, add(2, 3), ...
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil { out.WriteString(ls.Value.String()) }
	out.WriteString(";")
	return out.String()
}

type Identifier struct { // 実装の簡素化のため、識別子は値を生成する（文ではなく式）、letは本当は値を生成しない（文だから）
	Token token.Token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

// return <expression>
type ReturnStatement struct {
	Token token.Token // return
	Value Expression
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil { out.WriteString(rs.Value.String()) }
	out.WriteString(";")
	return out.String()
}

// 式文、x + 10;
type ExpressionStatement struct { // Statementを実装することでProgramのStatementsスライスに追加できる
	Token token.Token
	Expression Expression
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil { return es.Expression.String() }
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// !5 -15 ..etc
type PrefixExpression struct {
	Token token.Token
	Operator string 
	Right Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}
func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}