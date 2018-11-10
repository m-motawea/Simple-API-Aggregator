package parsers

import (
	"fmt"

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
	headYML, err := yaml.Marshal(resourceMap)
	if err != nil {
		return nil
	}
	return []byte(fmt.Sprintf(`
#%%RAML 1.0
title: %s
version: %s
baseUri: %s
%s
`, parser.Title, parser.BaseURI, parser.Version, headYML))
}
