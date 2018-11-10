package aggregator

import (
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/m-motawea/aggregator/parsers"
	"github.com/mie00/tokenizer"
)

var (
	resource = Resource{
		line:        "/",
		prefix:      "",
		childrenMap: make(map[string]*Resource),
	}
	baseUrl string
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
	url, qs := a.tokenizeURL()
	b := a.tokenizeBody()
	res := a.tokenizeResponse()
	resource.addChild(url, b, a.method, res, a.statusCode, map[string]string{}, qs)
	baseUrl = a.baseURL
}

func GetRaml() []byte {
	p := parsers.RamlParser{
		Title:   "Dummy",
		BaseURI: baseUrl,
		Version: "1.0",
	}
	return resource.GenerateSpecs(&p)
}

func (call *APICall) prepare() {
	call.method = call.request.Method
	call.line = call.request.URL.Path
	call.baseURL = call.request.Host
	call.fullPath = "http://" + call.baseURL + call.request.URL.String()
}

func (call *APICall) tokenizeURL() (string, map[string]string) {
	url := tokenizer.TokenizeString(call.fullPath)
	spew.Dump(call.fullPath, url)
	tokenizedLines := strings.Split(url, call.baseURL)
	line := strings.SplitN(tokenizedLines[1], "?", 2)
	u := line[0]
	qs := make(map[string]string)
	if len(line) == 2 {
		for _, kv := range strings.Split(line[1], "&") {
			if kva := strings.SplitN(kv, "=", 2); len(kva) == 2 {
				qs[kva[0]] = kva[1]
			}
		}
	}
	return u, qs
}

func (call *APICall) tokenizeBody() string {
	body := tokenizer.TokenizeString(string(call.body))
	return body
}

func (call *APICall) tokenizeResponse() string {
	body := tokenizer.TokenizeString(string(call.response))
	return body
}
