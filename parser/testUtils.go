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

func TestIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier.got %v", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident value not %v.got %v", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident TokenLiteral not %v.got %v.", value, ident.TokenLiteral())
		return false
	}
	return true
}

func TestLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return TestIntegerLiteral(t, exp, int64(v))
	case int64:
		return TestIntegerLiteral(t, exp, v)
	case string:
		return TestIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled.got [%v]", exp)
	return false
}

func TestInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression.got [%v]", exp)
	}

	if !TestLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if !TestLiteralExpression(t, opExp.Right, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %v.but got %v", operator, opExp.Operator)
		return false
	}
	return true
}
