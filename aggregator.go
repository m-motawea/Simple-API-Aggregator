package aggregator

import (
	"fmt"
	"strings"

	"github.com/m-motawea/aggregator/parsers"
)

type codeToResponse map[int]string

type methodParameters map[string]string

type Resource struct {
	line             string
	prefix           string
	isItem           bool
	methodToBody     map[string]string
	childrenMap      map[string]*Resource
	headersMap       map[string]string
	methodToResponse map[string]codeToResponse
	queryParameters  map[string]methodParameters
}

func (r *Resource) buildDescription() {
	fmt.Println("\n\ngenerating strucutre for Resource: ", r.line)
	fmt.Println("\tprefix: ", r.prefix)
	if r.isItem {
		fmt.Println("\t\tType: collection-item")
	} else {
		fmt.Println("\t\tType: collection")
	}
	for k, v := range r.methodToBody {
		fmt.Println("\t\t\tMethod: ", k)
		fmt.Println("\t\t\t\tBody:", v)
	}
}

func (r *Resource) buildTreeDescription() {
	for _, v := range r.childrenMap {
		v.buildDescription()
		v.buildTreeDescription()
	}
}

func (r *Resource) addChild(line string, body string, method string, response string, code int, headers map[string]string, queryParams map[string]string) {
	prefixes := strings.Split(line, "/")
	loopLen := len(prefixes) - 1
	parentResource := r
	var tmpResource *Resource
	i := 1
	for i < loopLen {
		child, ok := parentResource.childrenMap["/"+prefixes[i]]
		if !ok {
			tmpResource = parentResource
			parentResource = new(Resource)
			parentResource.prefix = "/" + prefixes[i]
			parentResource.line = strings.Join(prefixes[:i], "/")
			parentResource.childrenMap = make(map[string]*Resource)
			if strings.Contains("/"+prefixes[i], "{") {
				parentResource.isItem = true
			} else {
				parentResource.isItem = false
			}
			tmpResource.childrenMap["/"+prefixes[i]] = parentResource
		} else {
			parentResource = child
		}
		i++
	}
	c, ok := parentResource.childrenMap["/"+prefixes[i]]
	if !ok {
		c := new(Resource)
		c.line = line
		c.prefix = "/" + prefixes[i]
		c.childrenMap = make(map[string]*Resource)
		c.methodToBody = make(map[string]string)
		c.methodToBody[method] = body
		codeRes := make(codeToResponse)
		codeRes[code] = response
		c.methodToResponse = make(map[string]codeToResponse)
		c.methodToResponse[method] = codeRes
		c.headersMap = make(map[string]string)
		c.headersMap = headers
		c.queryParameters = make(map[string]methodParameters)
		mParams := make(methodParameters)
		for k, v := range queryParams {
			mParams[k] = v
		}
		c.queryParameters[method] = mParams
		if strings.Contains(c.prefix, "{") {
			c.isItem = true
		} else {
			c.isItem = false
		}
		parentResource.childrenMap["/"+prefixes[i]] = c
	} else {

		currentParams, ok := c.queryParameters[method]
		if !ok {
			c.queryParameters[method] = make(methodParameters)
			currentParams, _ = c.queryParameters[method]
		}
		for k, v := range queryParams {
			currentParams[k] = v
		}
		c.queryParameters[method] = currentParams
		c.methodToBody[method] = body
		_, k := c.methodToResponse[method]
		if !k {
			c.methodToResponse[method] = make(codeToResponse)
		}
		codeResponse := c.methodToResponse[method]
		codeResponse[code] = response
		c.methodToResponse[method] = codeResponse
		c.headersMap = headers
	}
}

func (r *Resource) GenerateSpecs(parser parsers.Parser) []byte {
	return parser.Parse(r)
}

func (r *Resource) GenerateMap() map[string]interface{} {
	resourceMap := make(map[string]interface{})
	if r.isItem {
		resourceMap["Type"] = "Collection-Item"
	} else {
		resourceMap["Type"] = "Collection"
	}
	for k, v := range r.methodToBody {
		methodMap := make(map[string]interface{})
		methodMap["body"] = v
		methodMap["queryParameters"] = r.queryParameters[k]
		methodMap["headers"] = r.headersMap
		contentType, ok := r.headersMap["Content-Type"]
		responsesMap := make(map[int]interface{})
		for code, response := range r.methodToResponse[k] {
			codeMap := make(map[string]interface{})
			exampleList := []string{response}
			bodyMap := make(map[string]interface{})
			bodyMap["examples"] = exampleList
			if ok {
				bodyMap["Content-Type"] = contentType
			}
			codeMap["body"] = bodyMap
			responsesMap[code] = codeMap

		}
		methodMap["responses"] = responsesMap
		resourceMap[k] = methodMap
	}
	for _, v := range r.childrenMap {
		childMap := v.GenerateMap()
		resourceMap[v.prefix] = childMap
	}
	return resourceMap
}
