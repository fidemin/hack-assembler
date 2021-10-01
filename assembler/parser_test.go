package assembler

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	reader := strings.NewReader(
		`@i
@100
(LOOP)
@LOOP1
D;JGT
D=D+A
D=D+A;JGT
(LOOP1)`)
	parser := NewParser(reader)
	if err := parser.Parse(); err != nil {
		t.Errorf("error happens: %s", err.Error())
		return
	}

	wantedCommands := []Command{
		{CommandType: ACommand, Symbol: "i", SymbolInt: 16},
		{CommandType: ACommand, Symbol: "100", SymbolInt: 100},
		{CommandType: LCommand, Symbol: "LOOP"},
		{CommandType: ACommand, Symbol: "LOOP1", SymbolInt: 6},
		{CommandType: CCommand, Dest: "", Comp: "D", Jump: "JGT"},
		{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: ""},
		{CommandType: CCommand, Dest: "D", Comp: "D+A", Jump: "JGT"},
		{CommandType: LCommand, Symbol: "LOOP1"},
	}

	if len(parser.Commands) != len(wantedCommands) {
		t.Errorf("len(parser.Commands) = %d, but want %d", len(parser.Commands), len(wantedCommands))
	}

	for i, got := range parser.Commands {
		wanted := wantedCommands[i]
		if got != wanted {
			t.Errorf("parser.Commands[%d] = +%v, but want +%v", i, got, wanted)
		}
	}
}

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
			t.Errorf("parser.Commands[%d] = +%v, but want +%v", i, got, wanted)
		}
	}
}

func TestParser_fillSymbolTable_success(t *testing.T) {
	reader := strings.NewReader(
		`@i
@100
(LOOP)
D;JGT
D=D+A
D=D+A;JGT
(LOOP1)`)

	parser := NewParser(reader)
	parser.parseToCommands()
	err := parser.fillSymbolTable()
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}

	if parser.symbolTable["LOOP"] != 2 {
		t.Errorf("LOOP label ROM addr %d, but want 2", parser.symbolTable["LOOP"])
		return
	}

	if parser.symbolTable["LOOP1"] != 5 {
		t.Errorf("LOOP1 label ROM addr %d, but want 5", parser.symbolTable["LOOP1"])
		return
	}
}

func TestParser_fillSymbolTable_error(t *testing.T) {
	reader := strings.NewReader(
		`@i
@100
(R0)
D;JGT
D=D+A
D=D+A;JGT
(LOOP1)`)

	parser := NewParser(reader)
	parser.parseToCommands()
	err := parser.fillSymbolTable()
	if err == nil {
		t.Errorf("parser.fillSymbolTable() should not be nil")
		return
	}
	fmt.Println(fmt.Sprintf("Error Message Check: %s", err))
}

func TestParser_parseACommandSymbolToInt(t *testing.T) {
	reader := strings.NewReader(
		`@i
@100
D;JGT
@R0
D=D+A
D=D+A;JGT
@i`)
	parser := NewParser(reader)
	parser.parseToCommands()
	parser.parseACommandSymbolToInt()

	if parser.Commands[0].SymbolInt != uint16(16) {
		t.Errorf("RAM Addr %d, but want 16", parser.Commands[0].SymbolInt)
		return
	}

	if parser.Commands[1].SymbolInt != uint16(100) {
		t.Errorf("RAM Addr %d, but want 100", parser.Commands[1].SymbolInt)
		return
	}

	if parser.Commands[3].SymbolInt != uint16(0) {
		t.Errorf("RAM Addr %d, but want 0", parser.Commands[3].SymbolInt)
		return
	}

	if parser.Commands[6].SymbolInt != uint16(16) {
		t.Errorf("RAM Addr %d, but want 16", parser.Commands[6].SymbolInt)
		return
	}

	if parser.currentRAMAddr != uint16(17) {
		t.Errorf("parser.currentRAMAddr %d, but want 17", parser.currentRAMAddr)
		return
	}

	if parser.symbolTable["i"] != uint16(16) {
		t.Errorf("parser.symbolTable[i] %d, but want 16", parser.symbolTable["i"])
		return
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
