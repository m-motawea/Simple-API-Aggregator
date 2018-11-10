package parsers

type Parser interface {
	Parse(r ResourceInterface) []byte
}

type ResourceInterface interface {
	GenerateMap() map[string]interface{}
}
