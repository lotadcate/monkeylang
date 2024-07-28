package token

/*
トークンには種類がある（タイプ）
トークンはタイプとリテラル（トークンの個別の情報）が必要

let five = 5;
let ten = 10;
let add = fn(x,y) {
  x + y;
}
let result = add(five, ten);
*/
type TokenType string
type Token struct {
	Type TokenType
	Literal string
}

// 優先順位
const (
	_int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

const (
	// 識別子
	IDENT = "IDENT" // five, ten, add, x, y

	// リテラル
	INT = "INT" // 5, 10

	// キーワード
	FUNCTION = "FUNCTION" // fn
	LET = "LET" // let
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"

	// 演算子
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	LT = "<"
	GT = ">"
	EQ = "=="
	NOT_EQ = "!="

	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
)

var keywords = map[string]TokenType { // 将来の追加を考えてvarにしている
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

// 引数の識別子がキーワードかどうか
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok 
	} else {
		return IDENT
	}
}

