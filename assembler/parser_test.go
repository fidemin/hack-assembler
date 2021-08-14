package assembler

import (
	"strings"
	"testing"
)

func TestParserAdvance(t *testing.T) {
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
