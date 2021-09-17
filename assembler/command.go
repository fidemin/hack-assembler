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

// IsNil determins Command struct has no data.
// This can be used when to check Command struct is like nil.
func (c Command) IsNil() bool {
	return c.CommandType == "" && c.Symbol == "" &&
		c.Dest == "" && c.Comp == "" && c.Jump == ""
}
