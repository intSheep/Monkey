package parser

import (
	"Monkey/ast"
	"testing"
)

func TestLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral want let,but got [%q]", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement.but got [%v]", letStmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s,but got [%s]", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not %s , but got [%s]", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}
