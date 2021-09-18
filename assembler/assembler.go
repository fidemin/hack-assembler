package assembler

import "io"

type Assembler struct {
	reader io.Reader
	writer io.Writer
}

func New(reader io.Reader, writer io.Writer) *Assembler {
	return &Assembler{
		reader: reader,
		writer: writer,
	}
}

func (a *Assembler) WriteBinaryCode() {
	parser := NewParser(a.reader)

	for parser.Advance() {
		command := parser.ParseOne()
		if command.IsNil() {
			continue
		}
		codebits := NewCodeBits(command)
		bits, err := codebits.Generate()
		if err != nil {
			panic(err.Error())
		}
		if _, err = a.writer.Write([]byte(bits)); err != nil {
			panic(err.Error())
		}
		a.writer.Write([]byte("\n"))
	}
}
