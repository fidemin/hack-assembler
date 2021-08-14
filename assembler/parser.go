package assembler

import (
	"bufio"
	"fmt"
	"io"
)

type Parser struct {
	scanner *bufio.Scanner
	currentCommand string
}


func NewParser(reader io.Reader) *Parser {
	parser := &Parser{}
	parser.scanner = bufio.NewScanner(reader)
	return parser
}

// Advance reads next line and make it to current command
func (p *Parser) Advance() bool {
	if p.scanner.Scan() {
		p.currentCommand = p.scanner.Text()
		return true
	}

	if err := p.scanner.Err(); err != nil {
		panic(fmt.Sprintf("error when reading command file: %s", err))
	}

	return false
}

