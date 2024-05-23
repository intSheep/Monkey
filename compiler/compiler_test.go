package compiler

import (
	"Monkey/ast"
	"Monkey/code"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"fmt"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedInstructions []code.Instructions
	expectedConstants    []any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{input: `1+2`, expectedConstants: []any{1, 2}, expectedInstructions: []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpAdd), code.Make(code.OpPop)}},
		{input: `1;2`, expectedConstants: []any{1, 2}, expectedInstructions: []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpPop), code.Make(code.OpConstant, 1), code.Make(code.OpPop)}},
		{input: `2-1`, expectedConstants: []any{2, 1}, expectedInstructions: []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpSub), code.Make(code.OpPop)}},
		{input: `4/2`, expectedConstants: []any{4, 2}, expectedInstructions: []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpDiv), code.Make(code.OpPop)}},
		{input: `3*7`, expectedConstants: []any{3, 7}, expectedInstructions: []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpMul), code.Make(code.OpPop)}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			runCompilerTest(t, tt)
		})
	}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{input: `true`, expectedConstants: []any{}, expectedInstructions: []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpPop)}},
		{input: `false`, expectedConstants: []any{}, expectedInstructions: []code.Instructions{code.Make(code.OpFalse), code.Make(code.OpPop)}},
		{input: `true == false`, expectedConstants: []any{}, expectedInstructions: []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpFalse), code.Make(code.OpPop)}},
		{input: `true !+ false`, expectedConstants: []any{}, expectedInstructions: []code.Instructions{code.Make(code.OpFalse), code.Make(code.OpPop)}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			runCompilerTest(t, tt)
		})
	}

}

func runCompilerTest(t *testing.T, tt compilerTestCase) {
	t.Helper()

	program := parse(tt.input)

	compiler := New()
	err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compiler error:%s", err)
	}

	bytecode := compiler.Bytecode()

	err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
	if err != nil {
		t.Fatalf("testInstructions fail: %v", err)
	}

	err = testConstants(tt.expectedConstants, bytecode.Constants)
	if err != nil {
		t.Fatalf("testConstants fail: %v", err)
	}

}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	contacted := concatInstructions(expected)
	if len(contacted) != len(actual) {
		return fmt.Errorf("wrong instructions length. \nwant=%q\ngot =%q", contacted, actual)
	}

	for i, ins := range contacted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q", i, contacted, actual)
		}
	}
	return nil
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, instruction := range instructions {
		out = append(out, instruction...)
	}
	return out
}

func testConstants(expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("")
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got =%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
