package assembler

import (
	"strings"
	"testing"
)

func TestParser_Advance(t *testing.T) {
	commands := `@i
M=1`
	reader := strings.NewReader(commands)
	parser := NewParser(reader)

	// test first line of commands
	wantedBool1 := true
	wantedCommand1 := "@i"

	if got := parser.Advance(); got != wantedBool1 {
		t.Errorf("Advance() = %t, want %t", got, wantedBool1)
	}

	if parser.currentCommand != wantedCommand1 {
		t.Errorf("currentCommand = %s, want %s", parser.currentCommand, wantedCommand1)
	}

	// test second line of commands
	wantedBool2 := true
	wantedCommand2 := "M=1"

	if got := parser.Advance(); got != wantedBool2 {
		t.Errorf("Advance() = %t, want %t", got, wantedBool2)
	}

	if parser.currentCommand != wantedCommand2 {
		t.Errorf("currentCommand = %s, want %s", parser.currentCommand, wantedCommand1)
	}

	// test no next line
	wantedBool3 := false

	if got := parser.Advance(); got != wantedBool3 {
		t.Errorf("Advance() = %t, want %t", got, wantedBool2)
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
