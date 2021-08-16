package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Parser parses hack assembly program line by line.
type Parser struct {
	scanner *bufio.Scanner
	currentCommand string
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

	if p.scanner.Scan() {
		p.currentCommand = strings.ReplaceAll(strings.TrimSpace(p.scanner.Text()), " ", "")
		return true
	}

	if err := p.scanner.Err(); err != nil {
		panic(fmt.Sprintf("error when reading command file: %s", err))
	}

	return false
}

// Parse parses currentCommand and convert it to Command object
func (p *Parser) Parse() Command {
	var symbol, dest, comp, jump string
	commandType := p.commandType()

	if commandType == ACommand {
		symbol = p.symbolFromACommand()
	} else if commandType == LCommand {
		symbol = p.symbolFromLCommand()
	} else if commandType == CCommand {
		dest, comp, jump = p.destCompJump()
	}

	return Command{
		CommandType: commandType,
		Symbol: symbol,
		Dest: dest,
		Comp: comp,
		Jump: jump,
	}
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

func (p *Parser) symbolFromACommand() string {
	return strings.TrimLeft(p.currentCommand, "@")
}

func (p *Parser) symbolFromLCommand() string {
	return strings.TrimRight(strings.TrimLeft(p.currentCommand, "("), ")")
}

func (p *Parser) destCompJump() (dest string, comp string, jump string) {
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
	return
}
