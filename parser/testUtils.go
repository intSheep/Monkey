package parser

import (
	"Monkey/ast"
	"fmt"
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

func TestIntegerLiteral(t *testing.T, right ast.Expression, want int64) bool {
	integ, ok := right.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("want *ast.IntegerLiteral,but got [%v]", right)
		return false
	}
	if integ.Value != want {
		t.Errorf("integ.Value want [%v],but got [%v]", want, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", want) {
		t.Errorf("integ.TokenLiteral want [%v],but got [%v]", want, integ.TokenLiteral())
		return false
	}
	return true
}
