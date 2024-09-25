package core

import "github.com/google/uuid"

type Generator struct {
	requestId string
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) generate() string {
	return uuid.New().String()

}
