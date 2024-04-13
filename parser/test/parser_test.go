package test

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/parser"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `let x=5;
let y=10;
let foobar = 838383;`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	parser.CheckErrors(t, p)
	if program == nil {
		t.Fatalf("parse program return nil")
	}

	if nums := len(program.Statements); nums != 3 {
		t.Fatalf("statement nums want 3, but got %d", nums)
	}

	tests := []struct{ expectIdentifier string }{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !parser.TestLetStatement(t, stmt, tt.expectIdentifier) {
			return
		}

	}
}

func TestReturnStatement(t *testing.T) {
	input := `return 5;
	return 10;
	return 1234;
`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	parser.CheckErrors(t, p)

	if program == nil {
		t.Fatalf("parse program return nil")
	}

	if nums := len(program.Statements); nums != 3 {
		t.Fatalf("parse program statements want [3] ,but got [%v]", nums)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStmt,got [%v]", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.Tokenliteral not 'return',but got %v", returnStmt.TokenLiteral())
		}

	}

}
