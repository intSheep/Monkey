package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup 传入opcode的byte
// 得到opcode的定义
func Lookup(op byte) (*Definition, error) {
	res, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return res, nil
}

// Make 传入操作码、操作数，得到字节码,其中操作数使用大端编码
//
//	 for example :
//		Make(OpConstant, []int{65534})
//	 返回[]byte{byte(OpConstant), 255, 254}
func Make(op Opcode, operand ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := def.OperandWidths[0] + 1
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operand {
		with := def.OperandWidths[i]
		switch with {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += with
	}
	return instruction
}

// ReadOperands : Make 的逆过程
// 传入opcode的定义、字节码指令
// 返回字节码的操作数operands和指令长度
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUnit16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUnit16 辅助函数
// 将[]byte转化为uint16
func ReadUnit16(ins []byte) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (i Instructions) String() string {
	return ""
}
