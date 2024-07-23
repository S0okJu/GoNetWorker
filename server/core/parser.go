package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parser struct {
	path string
	cfg  Config
}

func newParser(path string) Parser {
	return Parser{path: path}
}

func (p *Parser) parseJson(Parser Parser) {
	// Read file
	file, err := os.ReadFile(Parser.path)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}

	// Parse the json file

	err = json.Unmarshal([]byte(file), &p.cfg)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
	}
}
