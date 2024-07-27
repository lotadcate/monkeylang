package lexer

import (
	"monkey/token"
)

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

//先読み
func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

// l.characterを見てその文字に対応したトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	switch l.character {
	case '=':
		if l.peekChar() == '=' { // 次のトークンを覗き見
			ch := l.character
			l.readChar()
			literal := string(ch) + string(l.character)
			tok = token.Token{ Type: token.EQ, Literal: literal }
		} else {
			tok = newToken(token.ASSIGN, l.character)
		}
	case '+':
		tok = newToken(token.PLUS, l.character)
	case '-':
		tok = newToken(token.MINUS, l.character)
	case '*':
		tok = newToken(token.ASTERISK, l.character)
	case '/':
		tok = newToken(token.SLASH, l.character)
	case '!':
		if l.peekChar() == '=' {
			ch := l.character
			l.readChar()
			literal := string(ch) + string(l.character)
			tok = token.Token{ Type: token.NOT_EQ, Literal: literal }
		} else {
			tok = newToken(token.BANG, l.character)
		}
	case '<':
		tok = newToken(token.LT, l.character)
	case '>':
		tok = newToken(token.GT, l.character)
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
	default: // 読んでいる文字が識別子、リテラル、キーワードだった場合
		if isLetter(l.character) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.character){
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}else {
			tok = newToken(token.ILLEGAL, l.character)
		}
	}
	l.readChar()
	return tok
}

// 指定のトークンタイプでその文字をトークン化する
func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{ Type: tokenType, Literal: string(character)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.character) {
		l.readChar()
	}
	return l.input[position:l.position] // 識別子の初めの文字から終わりの文字まで（識別子自体）
}
// 小文字/大文字のアルファベット、アンダースコアを英字としている
func isLetter(character byte) bool {
	return ('a' <= character && character <= 'z') || ('A' <= character && character <= 'Z') || character == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.character) {
		l.readChar()
	}
	return l.input[position:l.position]
}
func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

// 空白は無視
func (l *Lexer) skipWhiteSpace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		l.readChar()
	}
}