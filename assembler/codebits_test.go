package assembler

import "testing"

func TestCodeBits_fromACommand(t *testing.T) {
	tests := []struct {
		command Command
		wanted string
		isErr bool
	}{
		{command: Command{CommandType: ACommand, Symbol: "abc"}, wanted: "", isErr: true},
		{command: Command{CommandType: ACommand, Symbol: "-1"}, wanted: "", isErr: true},
		{command: Command{CommandType: ACommand, Symbol: "32767"}, wanted: "1111111111111111", isErr: false},
		{command: Command{CommandType: ACommand, Symbol: "0"}, wanted: "1000000000000000", isErr: false},
		{command: Command{CommandType: ACommand, Symbol: "14"}, wanted: "1000000000001110", isErr: false},
	}

	for _, test := range tests {
		codebits := &CodeBits{command: test.command}
		got, err := codebits.fromACommand()

		if test.isErr {
			if err == nil {
				t.Errorf(
					"CodeBits.fromACommand(): error should not be nil with command %+v", test.command)
			}
		} else {
			if err != nil {
				t.Errorf(
					"CodeBits.fromACommand(): unexpected error '%s'  with %+v",
					err.Error(), test.command)
			} else {
				if got != test.wanted {
					t.Errorf("CodeBits.fromACommand() = %s, but want %s", got, test.wanted)
				}
			}

		}
	}
}


func TestCodeBits_fromCCommand(t *testing.T) {
	tests := []struct {
		command Command
		wanted string
		isErr bool
	}{
		{
			command: Command{CommandType: CCommand, Comp: "A+1", Dest: "MD", Jump: ""},
			wanted: "1110110111011000", isErr: false,
		},
		{
			command: Command{CommandType: CCommand, Comp: "0", Dest: "", Jump: "JMP"},
			wanted: "1110101010000111", isErr: false,
		},
		{
			command: Command{CommandType: CCommand, Comp: "K+1", Dest: "MD", Jump: ""},
			wanted: "", isErr: true,
		},
		{
			command: Command{CommandType: CCommand, Comp: "A+1", Dest: "ABC", Jump: ""},
			wanted: "", isErr: true,
		},
		{
			command: Command{CommandType: CCommand, Comp: "A+1", Dest: "MD", Jump: "JMT"},
			wanted: "", isErr: true,
		},
	}

	for _, test := range tests {
		codebits := &CodeBits{command: test.command}
		got, err := codebits.fromCCommand()

		if test.isErr {
			if err == nil {
				t.Errorf(
					"CodeBits.fromCCommand(): error should not be nil with command %+v", test.command)
			}
		} else {
			if err != nil {
				t.Errorf(
					"CodeBits.fromCCommand(): unexpected error '%s'  with %+v",
					err.Error(), test.command)
			} else {
				if got != test.wanted {
					t.Errorf("CodeBits.fromCCommand() = %s, but want %s", got, test.wanted)
				}
			}

		}
	}
}
