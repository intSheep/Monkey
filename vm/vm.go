package vm

import (
	"Monkey/code"
	"Monkey/compiler"
	"Monkey/object"
	"fmt"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 指向栈顶下一个位置的指针
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		// 直接取op并转化为操作码，而不是使用lookup，因为这会很慢
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUnit16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
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

func (vm *VM) StackTop() object.Object {
	return vm.stack[vm.sp-1]
}
