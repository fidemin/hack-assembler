package assembler

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


