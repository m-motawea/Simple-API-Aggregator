package aggregator

import (
	"github.com/mie00/tokenizer"
	"net/http"
	"strings"
)

type APICall struct {
	request  http.Request
	response http.ResponseWriter
	body     []byte
	// headers         map[string]string
	queryParameters map[string]string
	baseURL         string
	line            string
	fullPath        string
}

func (call *APICall) prepare() {
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
