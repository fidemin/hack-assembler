package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CommandType string

const (
	ACommand CommandType = "A"
	CCommand             = "C"
	LCommand             = "L"
)

// Parser parses hack assembly program line by line.
type Parser struct {
	scanner *bufio.Scanner
	currentCommand string
	currentCommandType CommandType
}

// NewParser returns *Parser object which has commands
func NewParser(reader io.Reader) *Parser {
	parser := &Parser{}
	parser.scanner = bufio.NewScanner(reader)
	return parser
}

// Advance reads next line and make it to current command
func (p *Parser) Advance() bool {
	if p.scanner.Scan() {
		p.currentCommand = strings.TrimSpace(p.scanner.Text())
		return true
	}

	if err := p.scanner.Err(); err != nil {
		panic(fmt.Sprintf("error when reading command file: %s", err))
	}

	return false
}

func (p *Parser) commandType() CommandType {
	if p.currentCommand[0] == '@' {
		p.currentCommandType = ACommand
	} else if p.currentCommand[0] == '(' {
		p.currentCommandType = LCommand
	} else {
		p.currentCommandType = CCommand
	}
	return p.currentCommandType
}

func (p *Parser) symbol() string {
	if p.currentCommandType == ACommand {
		return strings.TrimLeft(p.currentCommand, "@")
	}

	if p.currentCommandType == LCommand {
		return strings.TrimRight(strings.TrimLeft(p.currentCommand, "("), ")")
	}

	return ""
}

