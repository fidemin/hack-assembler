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
	CCommand CommandType = "C"
	LCommand CommandType = "L"
	NCommand CommandType = "N"
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
	// TODO: comment should be considered
	// TODO: only whitespace line should be considered

	// reset
	p.currentCommand = ""
	p.currentCommandType = NCommand

	if p.scanner.Scan() {
		p.currentCommand = strings.ReplaceAll(strings.TrimSpace(p.scanner.Text()), " ", "")
		p.currentCommandType = p.commandType()
		return true
	}

	if err := p.scanner.Err(); err != nil {
		panic(fmt.Sprintf("error when reading command file: %s", err))
	}

	return false
}

func (p *Parser) commandType() CommandType {
	if p.currentCommand[0] == '@' {
		return ACommand
	}

	if p.currentCommand[0] == '(' {
		return LCommand
	}

	return CCommand
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

func (p *Parser) destCompJump() (dest string, comp string, jump string) {
	if p.currentCommandType == CCommand {
		// C type command: dest=comp;jump
		// jump or dest can be omitted
		jumpSplit := strings.Split(p.currentCommand, ";")
		if len(jumpSplit) == 2 {
			// e.g. D;JGT
			jump = jumpSplit[1]
		}
		destComp := jumpSplit[0]
		destCompSplit := strings.Split(destComp, "=")
		if len(destCompSplit) == 2 {
			// e.g. D=D+A
			dest = destCompSplit[0]
			comp = destCompSplit[1]
		} else {
			// e.g. D;JGT
			comp = destCompSplit[0]
		}
	}
	return
}
