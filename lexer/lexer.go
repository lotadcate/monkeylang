package lexer

import "monkey/token"

/*
lexerは入力を先頭から読み込む
入力、今見ている文字、今見ている文字の位置（次に何が来るかを見る必要がある）、次の文字の位置
*/
type Lexer struct {
	input string
	character byte
	position int
	nextPosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 次の一文字を読んでinput文字列の現在位置を進める
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.character = 0 // NULに対応
	} else {
		l.character = l.input[l.nextPosition] // 次の文字
	}

	// positionの更新
	l.position = l.nextPosition
	l.nextPosition += 1
}

// l.characterを見てその文字に対応したトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.character {
	case '=':
		tok = newToken(token.ASSIGN, l.character)
	case '+':
		tok = newToken(token.PLUS, l.character)
	case '(':
		tok = newToken(token.LPAREN, l.character)
	case ')':
		tok = newToken(token.RPAREN, l.character)
	case '{':
		tok = newToken(token.LBRACE, l.character)
	case '}':
		tok = newToken(token.RBRACE, l.character)
	case ',':
		tok = newToken(token.COMMA, l.character)
	case ';':
		tok = newToken(token.SEMICOLON, l.character)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

// 指定のトークンタイプでその文字をトークン化する
func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{ Type: tokenType, Literal: string(character)}
}