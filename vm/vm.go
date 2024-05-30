package vm

import (
	"Monkey/code"
	"Monkey/compiler"
	"Monkey/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 65536

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	sp      int // 指向栈顶下一个位置的指针
	globals []object.Object
}

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm

}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		// 直接取op并转化为操作码，而不是使用lookup，因为这会很慢
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUnit16(vm.instructions[ip+1:]) //ReadUnit16期望读取两个字节，因此不用特地使用[ip+1:ip+3]
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpNotEqual, code.OpEqual, code.OpGreaterThan:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return nil
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUnit16(vm.instructions[ip+1:]))
			ip = pos - 1 // pos减一是因为循环的时候还会加一
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUnit16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUnit16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUnit16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= len(vm.stack) {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) StackTop() object.Object {
	return vm.stack[vm.sp-1]
}

// LastPoppedStackElem 指向栈顶下一个位置
// 用于测试栈清理的辅助函数
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.BOOLEAN_OBJ && rightType == object.BOOLEAN_OBJ:
		return vm.executeBinaryBooleanOperation(op, left, right)
	default:
		return fmt.Errorf("unsupport types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpEqual:
		if leftValue == rightValue {
			return vm.push(True)
		} else {
			return vm.push(False)
		}
	case code.OpNotEqual:
		if leftValue != rightValue {
			return vm.push(True)
		} else {
			return vm.push(False)
		}
	case code.OpGreaterThan:
		if leftValue > rightValue {
			return vm.push(True)
		} else {
			return vm.push(False)
		}

	default:
		return fmt.Errorf("unkonwn integer operator:%d", op)
	}
	return vm.push(&object.Integer{result})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation:%s", operand.Type())
	}
	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{-value})
}
func (vm *VM) executeBinaryBooleanOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	var result bool

	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	default:
		return fmt.Errorf("unkonwn integer operator:%d", op)
	}
	if result {
		return vm.push(True)
	} else {
		return vm.push(False)
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
