package assembler

type CommandType string

const (
	ACommand CommandType = "A"
	CCommand CommandType = "C"
	LCommand CommandType = "L"
)

type Command struct {
	CommandType CommandType
	Symbol string
	Dest string
	Comp string
	Jump string
}
