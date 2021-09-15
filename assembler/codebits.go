package assembler

import (
	"errors"
	"fmt"
)

// compToBits maps comp operation to bits.
// value is 7 bits string: a c1 c2 c3 c4 c5 c6
var compToBits = map[string] string {
	"0": "0101010",
	"1": "011111",
	"-1": "0111010",
	"D": "0001100",
	"A": "0110000",
	"!D": "0001101",
	"!A": "0110001",
	"-D": "0001111",
	"-A": "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0000010",
	"A-D": "0000111",
	"D&A" : "0000000",
	"D|A": "0010101",
	"M": "1110000",
	"!M": "1110001",
	"-M": "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

// destToBits maps dest operation to bits.
// value is 3 bits string: d1 d2 d3
var destToBits = map[string]string {
	"": "000",
	"M": "001",
	"D": "010",
	"MD": "011",
	"A": "100",
	"AM": "101",
	"AD": "110",
	"AMD": "111",
}

// jumpToBits maps jump operation to bits.
// value is 3 bits string: j1 j2 j3
var jumpToBits = map[string]string {
	"": "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

type CodeBits struct {
	command Command
}

func NewCodeBits(command Command) *CodeBits {
	return &CodeBits{
		command: command,
	}
}

func (b *CodeBits) generateFromACommand() (string, error) {
	symbol := b.command.Symbol

	// TODO: need to implement for that symbol is not integer but variable
	converter, err := NewIntToBitsConverter(symbol, 15, true)

	if err != nil {
		return "", err
	}

	return "1" + converter.ToBits(), nil
}

func (b *CodeBits) generateFromCCommand() (string, error) {
	compBits, ok := compToBits[b.command.Comp]
	if !ok {
		return "", errors.New(fmt.Sprintf("%s is not proper comp", b.command.Comp))
	}

	destBits, ok := destToBits[b.command.Dest]
	if !ok {
		return "", errors.New(fmt.Sprintf("%s is not proper dest", b.command.Dest))
	}

	jumpBits, ok := jumpToBits[b.command.Jump]
	if !ok {
		return "", errors.New(fmt.Sprintf("%s is not proper jump", b.command.Jump))
	}

	return "111" + compBits + destBits + jumpBits, nil
}
