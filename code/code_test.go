package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		name     string
		op       Opcode
		operands []int
		expected []byte
	}{
		{"OpConstant", OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{"OpAdd", OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instruction := Make(tt.op, tt.operands...)

			if len(instruction) != len(tt.expected) {
				t.Errorf("instruction has wrong length. want=%d, got=%d",
					len(tt.expected), len(instruction))
			}

			for i, b := range tt.expected {
				if instruction[i] != tt.expected[i] {
					t.Errorf("wrong byte at pos %d. want=%d, got=%d",
						i, b, instruction[i])
				}
			}

		})
	}
}

func TestReadOperand(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		Instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found:%s", err)
		}
		operandsRead, n := ReadOperands(def, Instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong.want =%d,got= %d", want, operandsRead[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
	}
	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Fatalf("instruction wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
	}
}
