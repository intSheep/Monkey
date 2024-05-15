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
