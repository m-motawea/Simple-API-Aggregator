package parsers

import (
	"os"
)

type Parser interface {
	Parse(r ResourceInterface, f *os.File) []byte
}

type ResourceInterface interface {
	GenerateMap() map[string]interface{}
}
