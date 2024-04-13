package test

import (
	"Monkey/ast"
	"Monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{Statements: []ast.Statement{
		&ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "myVar"},
				Value: "myVar",
			},
			Value: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
				Value: "anothervar",
			},
		},
	},
	}

	if program.String() != "let myVar=anothervar;" {
		t.Fatalf("program.String() want [let myVar=anothervar;], but got %v", program.String())
	}
}
