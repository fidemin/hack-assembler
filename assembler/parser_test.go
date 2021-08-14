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
		{command: "@100", wanted: A_COMMAND},
		{command: "(LOOP)", wanted: L_COMMAND},
		{command: "D=D-A", wanted: C_COMMAND},
	}

	for _, data := range testDataList {
		parser := &Parser{}
		parser.currentCommand = data.command
		if got := parser.commandType(); got != data.wanted {
			t.Errorf("commandType() = %s, want %s", got, data.wanted)
		}
	}
}
