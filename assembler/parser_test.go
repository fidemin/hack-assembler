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
	testDataList := []struct {
		command string
		wanted CommandType
	}{
		{command: "@100", wanted: ACommand},
		{command: "(LOOP)", wanted: LCommand},
		{command: "D=D-A", wanted: CCommand},
	}

	for _, data := range testDataList {
		parser := &Parser{}
		parser.currentCommand = data.command
		if got := parser.commandType(); got != data.wanted {
			t.Errorf("commandType() = %s, want %s", got, data.wanted)
		}
	}
}

func TestParser_symbol(t *testing.T) {
	testDataList := []struct {
		command string
		commandType CommandType
		wanted string
	}{
		{command: "(LOOP)", commandType: LCommand, wanted: "LOOP"},
		{command: "@100", commandType: ACommand, wanted: "100"},
		{command: "0;JMP", commandType: CCommand, wanted: ""},
	}

	for _, data := range testDataList {
		parser := &Parser{}
		parser.currentCommand = data.command
		parser.currentCommandType = data.commandType
		if got := parser.symbol(); got != data.wanted {
			t.Errorf("symbol() = %s, want %s", got, data.wanted)
		}
	}
}
