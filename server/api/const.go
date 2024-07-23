package api

type Code struct {
	number int
}

func (c *Code) Number() int {
	return c.number
}

var (
	OK                = Code{number: 200}
	CouldNotParseJson = Code{number: 4002}
)
