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

const (
	// 識別子
	IDENT = "IDENT" // five, ten, add, x, y

	// リテラル
	INT = "INT" // 5, 10

	// キーワード
	FUNCTION = "FUNCTION" // fn
	LET = "LET" // let

	// 演算子
	ASSIGN = "="
	PLUS = "+"

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
}

// 引数の識別子がキーワードかどうか
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok 
	} else {
		return IDENT
	}
}