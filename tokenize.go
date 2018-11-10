package aggregator

import (
	"github.com/m-motawea/aggregator/parsers"
	"github.com/mie00/tokenizer"
	"net/http"
	"strings"
)

var (
	resource = Resource{}
	baseUrl  string
)

type APICall struct {
	request    http.Request
	response   []byte
	body       []byte
	statusCode int
	// headers         map[string]string
	queryParameters map[string]string
	baseURL         string
	line            string
	fullPath        string
	method          string
}

func RegsiterCall(r http.Request, body []byte, response []byte, statusCode int) {
	a := APICall{
		request:    r,
		response:   response,
		body:       body,
		statusCode: statusCode,
	}
	a.prepare()
	url := a.tokenizeURL()
	b := a.tokenizeBody()
	res := a.tokenizeResponse()
	resource.addChild(url, b, a.method, res, a.statusCode, map[string]string{}, a.queryParameters)
	baseUrl = a.baseURL
}

func GetRaml() []byte {
	p := parsers.RamlParser{
		Title:   "Dummy",
		BaseURI: baseUrl,
		Version: "1.0",
	}
	return resource.GenerateSpecs(p)
}

func (call *APICall) prepare() {
	call.method = call.request.Method
	call.line = call.request.URL.Path
	call.baseURL = call.request.Host
	q := call.request.URL.Query()
	for k, v := range q {
		if len(v) >= 1 {
			// TODO: handle multiple query parameters
			call.queryParameters[k] = v[0]
		}
	}
	call.fullPath = call.baseURL + call.line
}

func (call *APICall) tokenizeURL() string {
	url := tokenizer.TokenizeString(call.fullPath)
	tokenizedLines := strings.Split(url, call.baseURL)
	if len(tokenizedLines) > 1 {
		line := tokenizedLines[1]
		return strings.Split(line, "?")[0]
	} else {
		return "/"
	}
}

func (call *APICall) tokenizeBody() string {
	body := tokenizer.TokenizeString(string(call.body))
	return body
}

func (call *APICall) tokenizeResponse() string {
	body := tokenizer.TokenizeString(string(call.response))
	return body
}
