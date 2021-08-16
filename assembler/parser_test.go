package assembler

import (
	"strings"
	"testing"
)

func TestParser_Advance(t *testing.T) {
	commands := `@i
D=D-A`
	reader := strings.NewReader(commands)
	parser := NewParser(reader)

	tests := []struct {
		command     string
		commandType CommandType
		advance bool
	}{
		{command: "@i", commandType: ACommand, advance: true},
		{command: "D=D-A", commandType: CCommand, advance: true},
		{command: "", commandType: NCommand, advance: false},
	}

	for _, test := range tests {
		if got := parser.Advance(); got != test.advance {
			t.Errorf("Advance() = %t, want %t", got, test.advance)
		}

		if parser.currentCommand != test.command  {
			t.Errorf("currentCommand = %s, want %s", parser.currentCommand, test.command)
		}

		if parser.currentCommandType != test.commandType {
			t.Errorf(
				"currentCommandType = %s, want %s",
				parser.currentCommandType, test.commandType)
		}
	}
}

func TestParser_commandType(t *testing.T) {
	tests := []struct {
		command string
		wanted CommandType
	}{
		{command: "@100", wanted: ACommand},
		{command: "(LOOP)", wanted: LCommand},
		{command: "D=D-A", wanted: CCommand},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		if got := parser.commandType(); got != test.wanted {
			t.Errorf("commandType() = %s, want %s", got, test.wanted)
		}
	}
}

func TestParser_symbol(t *testing.T) {
	tests := []struct {
		command string
		commandType CommandType
		wanted string
	}{
		{command: "(LOOP)", commandType: LCommand, wanted: "LOOP"},
		{command: "@100", commandType: ACommand, wanted: "100"},
		{command: "0;JMP", commandType: CCommand, wanted: ""},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		parser.currentCommandType = test.commandType
		if got := parser.symbol(); got != test.wanted {
			t.Errorf("symbol() = %s, want %s", got, test.wanted)
		}
	}
}

func TestParser_destCompJump(t *testing.T) {
	tests := []struct {
		command string
		commandType CommandType
		dest string
		comp string
		jump string
	} {
		{command: "D;JGT", commandType: CCommand, dest: "", comp: "D", jump: "JGT"},
		{command: "D=D+A", commandType: CCommand, dest: "D", comp: "D+A", jump: ""},
		{command: "D=D+A;JGT", commandType: CCommand, dest: "D", comp: "D+A", jump: "JGT"},
		{command: "@i", commandType: ACommand, dest: "", comp: "", jump: ""},
		{command: "(LOOP)", commandType: LCommand, dest: "", comp: "", jump: ""},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		parser.currentCommandType = test.commandType
		dest, comp, jump := parser.destCompJump()
		if (dest != test.dest) || (comp != test.comp) || (jump != test.jump) {
			t.Errorf(
				"destCompJump() = %s, %s, %s, want %s, %s, %s",
				dest, comp, jump, test.dest, test.comp, test.jump)
		}
	}
}
