package test

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/parser"
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Println("program:", program.String())
	parser.CheckErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program statement num want [1] , but got [%v]", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] want type [*ast.ExpressionStatement] , but got [%v]", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier,got [%v]", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("value want [foobar],but got [%v]", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.Tokenliteral want [foobar],but got [%v]", ident.TokenLiteral())
	}
}

func TestIntgerLiteralExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Println("program:", program.String())
	parser.CheckErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program statement num want [1] , but got [%v]", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] want type [*ast.ExpressionStatement] , but got [%v]", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral,got [%v]", stmt.Expression)
	}

	if ident.Value != 5 {
		t.Fatalf("value want [5],but got [%v]", ident.Value)
	}

	if ident.TokenLiteral() != "5" {
		t.Fatalf("ident.Tokenliteral want [5],but got [%v]", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		fmt.Println("program:", program.String())
		parser.CheckErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program statement num want [1] , but got [%v]", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] want type [*ast.ExpressionStatement] , but got [%v]", program.Statements[0])
		}

		prefix, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression,got [%v]", stmt.Expression)
		}

		if prefix.Operator != tt.operator {
			t.Fatalf("want operator[%v],but got [%v]", tt.operator, prefix.Operator)
		}
		if !parser.TestIntegerLiteral(t, prefix.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		parser.CheckErrors(t, p)
		fmt.Println(program.String())
		if len(program.Statements) != 1 {
			t.Fatalf("program statement num want [1] , but got [%v]", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] want type [*ast.ExpressionStatement] , but got [%v]", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression,got [%v]", stmt.Expression)
		}

		if !parser.TestIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("want operator[%v],but got [%v]", tt.operator, exp.Operator)
		}
		if !parser.TestIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}
