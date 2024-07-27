package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) { // Tはユニットテスト
	input := `=+(){},;`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{ token.ASSIGN, "=" },
		{ token.PLUS, "+" },
		{ token.LPAREN, "(" },
		{ token.RPAREN, ")" },
		{ token.LBRACE, "{" },
		{ token.RBRACE, "}" },
		{ token.COMMA, "," },
		{ token.SEMICOLON, ";" },
		{ token.EOF, "" },
	}

	 l := New(input)
	 for i, tt := range tests {
		token := l.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokenliteral wrong. expected=%q, got=%q", i, tt.expectedLiteral, token.Literal)
		}
	 }
	
}
