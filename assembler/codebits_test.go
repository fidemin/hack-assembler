package assembler

import "testing"

func TestCodeBits_generateFromACommand(t *testing.T) {
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
		got, err := codebits.generateFromACommand()

		if test.isErr {
			if err == nil {
				t.Errorf("CodeBits.generateFromACommand(): error should not be nil with command %+v", test.command)
			}
		} else {
			if err != nil {
				t.Errorf("CodeBits.generateFromACommand(): unexpected error '%s'  with %+v", err.Error(), test.command)
			} else {
				if got != test.wanted {
					t.Errorf("CodeBits.generateFromACommand() = %s, but want %s", got, test.wanted)
				}
			}

		}
	}
}
