package ast

import "monkey/token"

/*
let <identifier> = <expression>; <- これで一つのstatement

ASTはNodeだけで構成される
Programノードは全てのASTのルートノード
*/

type Node interface {
	TokenLiteral() string // ノードが関連づけられているトークンのリテラル値を返す、デバッグとテストのために使う
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

type LetStatement struct { 
	Token token.Token // let
	Name *Identifier // 識別子 ex) x, y, ...
	Value Expression // 値 ex) 5, add(2, 3), ...
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }


type Identifier struct { // 実装の簡素化のため、識別子は値を生成する（文ではなく式）、letは本当は値を生成しない（文だから）
	Token token.Token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// return <expression>
type ReturnStatement struct {
	Token token.Token // return
	Value Expression
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }