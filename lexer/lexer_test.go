package lexer_test

import (
	"Monkey/lexer"
	"Monkey/token"
	"testing"
)

func Test_Simple_Lexer(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectType    string
		expectLiteral string
	}{
		{expectType: token.ASSIGN, expectLiteral: "="},
		{expectType: token.PLUS, expectLiteral: "+"},
		{expectType: token.LPAREN, expectLiteral: "("},
		{expectType: token.RPAREN, expectLiteral: ")"},
		{expectType: token.LBRACE, expectLiteral: "{"},
		{expectType: token.RBRACE, expectLiteral: "}"},
		{expectType: token.COMMA, expectLiteral: ","},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tt.expectType != tt.expectLiteral {
			t.Fatalf("tests[%d]-token wrong.expected=%q, got=%q", i, tt.expectType, tok.Type)
		}
		if tt.expectLiteral != tt.expectLiteral {
			t.Fatalf("tests[%d]-literal wrong.expected=%q, got=%q", i, tt.expectLiteral, tok.Literal)
		}
	}
}
