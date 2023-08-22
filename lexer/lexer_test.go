package lexer

import (
	"testing"

	"github.com/lautarojayat/writing-an-interpreter-in-go/token"
)

func TestNewToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expecedType     token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expecedType {
			t.Fatalf("tests[%d] - wrong token type. Expected %s but got %s", i, tt.expecedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token literal. Expected %s but got %s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
