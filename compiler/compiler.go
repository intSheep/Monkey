package compiler

import (
	"Monkey/ast"
	"Monkey/code"
	"Monkey/object"
	"fmt"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object // 常量池
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer)) //
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil
}

// Bytecode 包含编译器生成的instructions和求值的constants
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// addConstant 辅助函数，往常量池添加常量，并返回索引
func (c *Compiler) addConstant(constant object.Object) int {
	c.constants = append(c.constants, constant)
	return len(c.constants) - 1
}

// emit 生成指令，并将其添加至内存
// 返回指令的位置
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	inst := code.Make(op, operands...)
	pos := c.addInstruction(inst)
	return pos
}

// addInstruction 辅助函数
// 用于将指令添加到内存中
func (c *Compiler) addInstruction(inst code.Instructions) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, inst...)
	return posNewInstruction
}
