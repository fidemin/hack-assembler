package assembler

import (
	"strings"
	"testing"
)


func TestParser_parseToCommands(t *testing.T) {
	reader := strings.NewReader(
		`@i
@100
(LOOP)
D;JGT
D=D+A
D=D+A;JGT`)

	wantedCommands := []Command{
		{CommandType: ACommand, Symbol: "i"},
		{CommandType: ACommand, Symbol: "100"},
		{CommandType: LCommand, Symbol: "LOOP"},
		{CommandType: CCommand, Dest: "", Comp: "D", Jump: "JGT"},
		{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: ""},
		{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: "JGT"},
	}

	parser := NewParser(reader)
	parser.parseToCommands()

	if len(parser.Commands) != len(wantedCommands) {
		t.Errorf("len(parser.Commands) = %d, but want %d", len(parser.Commands), len(wantedCommands))
	}

	for i, got := range parser.Commands {
		wanted := wantedCommands[i]
		if got != wanted {
			t.Errorf("parser.Commands[%d] = %s, but want %s", i, got, wanted)
		}
	}
}

func TestParser_Advance(t *testing.T) {
	commands := `@i
D=D-A`
	reader := strings.NewReader(commands)
	parser := NewParser(reader)

	tests := []struct {
		command string
		advance bool
	}{
		{command: "@i", advance: true},
		{command: "D=D-A", advance: true},
		{command: "", advance: false},
	}

	for _, test := range tests {
		if got := parser.Advance(); got != test.advance {
			t.Errorf("Advance() = %t, want %t", got, test.advance)
		}

		if parser.currentCommand != test.command  {
			t.Errorf("currentCommand = %s, want %s", parser.currentCommand, test.command)
		}
	}
}

func TestParser_ParseOne(t *testing.T) {
	tests := []struct {
		commandString string
		command Command
	}{
		{commandString: "@i", command: Command{CommandType: ACommand, Symbol: "i"}},
		{commandString: "@100", command: Command{CommandType: ACommand, Symbol: "100"}},
		{commandString: "(LOOP)", command: Command{CommandType: LCommand, Symbol: "LOOP"}},
		{commandString: "D;JGT", command: Command{CommandType: CCommand, Dest: "", Comp: "D", Jump: "JGT"}},
		{commandString: "D=D+A", command: Command{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: ""}},
		{commandString: "D=D+A;JGT", command: Command{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: "JGT"}},
	}

	for _, test := range tests {
		parser := Parser{}
		parser.currentCommand = test.commandString
		if got := parser.ParseOne(); got != test.command {
			t.Errorf("ParseOne() = %+v, want %+v", got, test.command)
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

func TestParser_symbolFromACommand(t *testing.T) {
	tests := []struct {
		command string
		wanted string
	}{
		{command: "@i", wanted: "i"},
		{command: "@100", wanted: "100"},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		if got := parser.symbolFromACommand(); got != test.wanted {
			t.Errorf("symbolFromACommand() = %s, want %s", got, test.wanted)
		}
	}
}

func TestParser_symbolFromLCommand(t *testing.T) {
	tests := []struct {
		command string
		wanted string
	}{
		{command: "(LOOP)", wanted: "LOOP"},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		if got := parser.symbolFromLCommand(); got != test.wanted {
			t.Errorf("symbolFromLCommand() = %s, want %s", got, test.wanted)
		}
	}
}

func TestParser_destCompJump(t *testing.T) {
	tests := []struct {
		command string
		dest string
		comp string
		jump string
	} {
		{command: "D;JGT", dest: "", comp: "D", jump: "JGT"},
		{command: "D=D+A", dest: "D", comp: "D+A", jump: ""},
		{command: "D=D+A;JGT", dest: "D", comp: "D+A", jump: "JGT"},
	}

	for _, test := range tests {
		parser := &Parser{}
		parser.currentCommand = test.command
		dest, comp, jump := parser.destCompJump()
		if (dest != test.dest) || (comp != test.comp) || (jump != test.jump) {
			t.Errorf(
				"destCompJump() = %s, %s, %s, want %s, %s, %s",
				dest, comp, jump, test.dest, test.comp, test.jump)
		}
	}
}
