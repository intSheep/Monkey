package test

import (
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
		t.Fatalf("parseprogram return nil")
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
