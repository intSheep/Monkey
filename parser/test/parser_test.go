package test

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/parser"
	"fmt"
	"github.com/stretchr/testify/require"
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

	infixTest2 := []struct {
		input      string
		leftValue  bool
		operator   string
		rightValue bool
	}{
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false==false", false, "==", false},
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

		if !parser.TestInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}

	for _, tt := range infixTest2 {
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

		if exp.Operator != tt.operator {
			t.Fatalf("want operator[%v],but got [%v]", tt.operator, exp.Operator)
		}

		if !parser.TestInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input  string
		target bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Println(program.String())
		parser.CheckErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program statement num want [1] , but got [%v]", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("stmt not *ast.ExpressionStatement. got [%v]", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Errorf("exp want *ast.Boolean,but got [%v]", stmt.Expression)
		}

		if exp.Value != tt.target {
			t.Errorf("stmt.Value want [%v] but got [%v] ", tt.target, exp.Value)
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		{"a+b-c", "((a + b) - c)"},
		{"a*b*c", "((a * b) * c)"},
		{"a*b/c", "((a * b) / c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3+4;-5*5", "(3 + 4)((-5) * 5)"},
		{"5>4==3<4", "((5 > 4) == (3 < 4))"},
		{"5<4!=3>4", "((5 < 4) != (3 > 4))"},
		{"3+4*5==3*1+4*5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3>5==false", "((3 > 5) == false)"},
		{"3<5==true", "((3 < 5) == true)"},
		{"1+(2+3)+4", "((1 + (2 + 3)) + 4)"},
		{"(5+5)*2", "((5 + 5) * 2)"},
		{"2/(5+5)", "(2 / (5 + 5))"},
		{"-(5+5)", "(-(5 + 5))"},
		{"!(true==true)", "(!(true == true))"},
		{"a+add(b*c)+d", "((a + add((b * c))) + d)"},
		{"add(a,b,1,2*3,4+5,add(6,7*8))", "add(a,b,1,(2 * 3),(4 + 5),add(6,(7 * 8)))"},
		{"(3-2)", "(3 - 2)"},
		//{"a*[1,2,3,4][b*c]*d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		//{"add(a*b[2], b[1], 2 * [1,2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		parser.CheckErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"if(x<y){x}", "if ( (x < y) ) { x }"},
		{"if(x<y){x}else{y}", "if ( (x < y) ) { x } else{ y }"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		parser.CheckErrors(t, p)
		fmt.Println(program.String())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`fn(x,y){x+y;}`, `fn(x,y){ (x + y) }`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		parser.CheckErrors(t, p)

		fmt.Println(program.String())
		if program.String() != tt.expect {
			t.Fatalf("want [%v],but got  [%v]", tt.expect, program.String())
		}
	}
}

func TestCallFunction(t *testing.T) {
	input := `add(1,2*3,4+5);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	parser.CheckErrors(t, p)

	require.Equal(t, len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.Expression,got [%v]", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt not *ast.CallExpression,got [%v]", stmt.Expression)
	}

	if !parser.TestIdentifier(t, exp.Function, "add") {
		return
	}

	require.Equal(t, len(exp.Arguments), 3)
	parser.TestLiteralExpression(t, exp.Arguments[0], 1)
	parser.TestInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	parser.TestInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x=5;", "x", 5},
		{"let y =true;", "y", true},
		{"let foobar =y;", "foobar", "y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			parser.CheckErrors(t, p)

			require.Equal(t, len(program.Statements), 1)
			stmt := program.Statements[0]
			if !parser.TestLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.LetStatement).Value
			if !parser.TestLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})

	}
}

// parser/parser_test.go

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	parser.CheckErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	parser.CheckErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	parser.TestIntegerLiteral(t, array.Elements[0], 1)
	parser.TestInfixExpression(t, array.Elements[1], 2, "*", 2)
	parser.TestInfixExpression(t, array.Elements[2], 3, "+", 3)
}

// parser/parser_test.go

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	parser.CheckErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !parser.TestIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !parser.TestInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}
