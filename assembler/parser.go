package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var predefinedSymbol = map[string]uint16 {
	"SP": 0,
	"LCL": 1,
	"ARG": 2,
	"THIS": 3,
	"THAT": 4,
	"R0": 0,
	"R1": 1,
	"R2": 2,
	"R3": 3,
	"R4": 4,
	"R5": 5,
	"R6": 6,
	"R7": 7,
	"R8": 8,
	"R9": 9,
	"R10": 10,
	"R11": 11,
	"R12": 12,
	"R13": 13,
	"R14": 14,
	"R15": 15,
	"SCREEN": 16384,
	"KBD": 24576,
}

// Parser parses hack assembly program line by line.
type Parser struct {
	scanner *bufio.Scanner
	currentCommand string
	currentRAMAddr uint16
	reservedSymbol map[string]uint16
}

// NewParser returns *Parser object which has commands
func NewParser(reader io.Reader) *Parser {
	parser := &Parser{}
	parser.scanner = bufio.NewScanner(reader)
	parser.currentRAMAddr = 16
	parser.reservedSymbol = make(map[string]uint16)
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
	symbol := strings.TrimLeft(p.currentCommand, "@")

	symbolInt, err := strconv.ParseUint(symbol, 10, 16)
	if err != nil {
		// symbol is string
		// check predefinedSymbol first
		if symbolInt, ok := predefinedSymbol[symbol]; ok {
			return string(symbolInt)
		}

		// if predefiendSymbol has no symbol, check reservedSymbol
		if symbolInt, ok := p.reservedSymbol[symbol]; ok {
			return string(symbolInt)
		}

		// Otherwise, allocate symbol to new RAM address
		currentRamAddr := p.currentRAMAddr
		p.reservedSymbol[symbol] = currentRamAddr
		p.currentRAMAddr += 1
		return string(currentRamAddr)
	}
	return string(symbolInt)
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
