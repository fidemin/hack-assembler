package assembler

import (
	"bytes"
	"strings"
	"testing"
)

func TestAssembler_WriteBinaryCode(t *testing.T) {
	reader := strings.NewReader(
		`@14
D;JGT
D=D+A
M=M-1;JMP
`)
	writer := new(bytes.Buffer)
	assembler := New(reader, writer)
	assembler.WriteBinaryCode()

	wanted := `1000000000001110
1110001100000001
1110000010010000
1111110010001111
`
	got := writer.String()
	if got != wanted {
		t.Errorf(`assembly code is not properly converted to binary code
wanted:
%s
but got:
%s
`, wanted, got)
	}
}
