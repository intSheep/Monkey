package lexer_test

import (
	"Monkey/lexer"
	"Monkey/token"
	"testing"
)

func Test_Simple_Lexer(t *testing.T) {
	input := "=+(){},;-/*<>!"

	tests := []struct {
		expectType    token.TokenType
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
		{expectType: token.MINUS, expectLiteral: "-"},
		{expectType: token.SLASH, expectLiteral: "/"},
		{expectType: token.ASTERISK, expectLiteral: "*"},
		{expectType: token.LT, expectLiteral: "<"},
		{expectType: token.GT, expectLiteral: ">"},
		{expectType: token.BANG, expectLiteral: "!"},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectType {
			t.Fatalf("tests[%d]-token wrong.expected=%q, got=%q", i, tt.expectType, tok.Type)
		}
		if tok.Literal != tt.expectLiteral {
			t.Fatalf("tests[%d]-literal wrong.expected=%q, got=%q", i, tt.expectLiteral, tok.Literal)
		}
	}
}

func Test_Complex_Lexer(t *testing.T) {
	input := `let five = 5;
	let ten =10;
	let add =fn(x,y){
	x+y;
};
let result = add(five,ten);
5 < 10 >5;
if(5<10){
	return true;
}else{
return false;
}
10 == 10;
10 != 9;
"foobar"
"foo bar"
`
	tests := []struct {
		expectType    token.TokenType
		expectLiteral string
	}{
		{expectType: token.LET, expectLiteral: "let"},
		{expectType: token.IDENT, expectLiteral: "five"},
		{expectType: token.ASSIGN, expectLiteral: "="},
		{expectType: token.INT, expectLiteral: "5"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.LET, expectLiteral: "let"},
		{expectType: token.IDENT, expectLiteral: "ten"},
		{expectType: token.ASSIGN, expectLiteral: "="},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.LET, expectLiteral: "let"},
		{expectType: token.IDENT, expectLiteral: "add"},
		{expectType: token.ASSIGN, expectLiteral: "="},
		{expectType: token.FUNCTION, expectLiteral: "fn"},
		{expectType: token.LPAREN, expectLiteral: "("},
		{expectType: token.IDENT, expectLiteral: "x"},
		{expectType: token.COMMA, expectLiteral: ","},
		{expectType: token.IDENT, expectLiteral: "y"},
		{expectType: token.RPAREN, expectLiteral: ")"},
		{expectType: token.LBRACE, expectLiteral: "{"},
		{expectType: token.IDENT, expectLiteral: "x"},
		{expectType: token.PLUS, expectLiteral: "+"},
		{expectType: token.IDENT, expectLiteral: "y"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.RBRACE, expectLiteral: "}"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.LET, expectLiteral: "let"},
		{expectType: token.IDENT, expectLiteral: "result"},
		{expectType: token.ASSIGN, expectLiteral: "="},
		{expectType: token.IDENT, expectLiteral: "add"},
		{expectType: token.LPAREN, expectLiteral: "("},
		{expectType: token.IDENT, expectLiteral: "five"},
		{expectType: token.COMMA, expectLiteral: ","},
		{expectType: token.IDENT, expectLiteral: "ten"},
		{expectType: token.RPAREN, expectLiteral: ")"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.INT, expectLiteral: "5"},
		{expectType: token.LT, expectLiteral: "<"},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.GT, expectLiteral: ">"},
		{expectType: token.INT, expectLiteral: "5"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.IF, expectLiteral: "if"},
		{expectType: token.LPAREN, expectLiteral: "("},
		{expectType: token.INT, expectLiteral: "5"},
		{expectType: token.LT, expectLiteral: "<"},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.RPAREN, expectLiteral: ")"},
		{expectType: token.LBRACE, expectLiteral: "{"},
		{expectType: token.RETURN, expectLiteral: "return"},
		{expectType: token.TRUE, expectLiteral: "true"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.RBRACE, expectLiteral: "}"},
		{expectType: token.ELSE, expectLiteral: "else"},
		{expectType: token.LBRACE, expectLiteral: "{"},
		{expectType: token.RETURN, expectLiteral: "return"},
		{expectType: token.FALSE, expectLiteral: "false"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.RBRACE, expectLiteral: "}"},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.EQ, expectLiteral: "=="},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.INT, expectLiteral: "10"},
		{expectType: token.NOT_EQ, expectLiteral: "!="},
		{expectType: token.INT, expectLiteral: "9"},
		{expectType: token.SEMICOLON, expectLiteral: ";"},
		{expectType: token.STRING, expectLiteral: "foobar"},
		{expectType: token.STRING, expectLiteral: "foo bar"},
		{expectType: token.EOF, expectLiteral: ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectType {
			t.Fatalf("tests[%d]-token wrong.expected=%q, got=%q", i, tt.expectType, tok.Type)
		}
		if tok.Literal != tt.expectLiteral {
			t.Fatalf("tests[%d]-literal wrong.expected=%q, got=%q", i, tt.expectLiteral, tok.Literal)
		}
	}

}
