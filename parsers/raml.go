package parsers

import (
	"gopkg.in/yaml.v2"
)

type RamlParser struct {
	Title   string `yaml:Title`
	BaseURI string `yaml:baseUri`
	Version string `yaml:version`
}

type ResourceRaml struct {
	prefix string
}

func (parser *RamlParser) Parse(root ResourceInterface) []byte {
	resourceMap := root.GenerateMap()
	resourceMap["Title"] = parser.Title
	resourceMap["baseUri"] = parser.BaseURI
	resourceMap["version"] = parser.Version
	headYML, err := yaml.Marshal(resourceMap)
	if err != nil {
		return nil
	}
	return headYML
}
