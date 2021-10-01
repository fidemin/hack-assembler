package assembler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser parses hack assembly program codes.
type Parser struct {
	Commands []Command
	scanner *bufio.Scanner
	currentCommand string
	currentRAMAddr uint16
	currentROMAddr uint16
	symbolTable map[string]uint16
}

// NewParser returns *Parser object which has commands
func NewParser(reader io.Reader) *Parser {
	parser := &Parser{}
	parser.scanner = bufio.NewScanner(reader)

	// initializing symbolTable for predefined symbols
	parser.symbolTable = map[string]uint16 {
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

	// initialize ROM Addr counter
	parser.currentROMAddr = 0

	// initialize RAM Addr counter
	parser.currentRAMAddr = 16

	parser.Commands = []Command{}

	return parser
}

// Parse parses all hack assembly program codes to Command structs
func (p *Parser) Parse() error {
	p.parseToCommands()
	if err := p.fillSymbolTable(); err != nil {
		return err
	}
	p.parseACommandSymbolToInt()
	return nil
}

func (p *Parser) parseToCommands() {
	for p.Advance() {
		p.Commands = append(p.Commands, p.ParseOne())
	}
}

func (p *Parser) fillSymbolTable() error {
	// check LCommand's symbol exists in predefined symbol.
	// If exists, returns error
	for _, command := range p.Commands {
		if command.CommandType == LCommand {
			if _, ok := p.symbolTable[command.Symbol]; ok {
				return errors.New(fmt.Sprintf("%s label already exists in predefined symbol", command.Symbol))
			}
		}
	}

	// fill SymbolTable with label and ROM address
	// TODO: check when LCommand has same label -> parse error
	for _, command := range p.Commands {
		if command.CommandType == LCommand {
			if _, ok := p.symbolTable[command.Symbol]; !ok {
				p.symbolTable[command.Symbol] = p.currentROMAddr
			}
		}
		if command.CommandType == CCommand || command.CommandType == ACommand {
			p.currentROMAddr += 1
		}
	}
	return nil
}

func (p *Parser) parseACommandSymbolToInt() {
	for i, command := range p.Commands {
		if command.CommandType == ACommand {
			symbol := command.Symbol
			symbolInt, err := strconv.ParseUint(symbol, 10, 64)
			// symbol is int string
			if err == nil {
				p.Commands[i].SymbolInt = uint16(symbolInt)
				continue
			}

			// symbol is not int string but variable
			if symbolInt16, ok := p.symbolTable[symbol]; ok {
				p.Commands[i].SymbolInt = symbolInt16
			} else {
				p.symbolTable[symbol] = p.currentRAMAddr
				p.Commands[i].SymbolInt = p.currentRAMAddr
				p.currentRAMAddr += 1
			}
		}
	}
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

// ParseOne parses currentCommand and convert it to Command object
func (p *Parser) ParseOne() Command {
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
