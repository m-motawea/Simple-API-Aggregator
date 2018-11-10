package aggregator

import (
	"fmt"
	"github.com/aggregator/parsers"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func getRandomLine() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	file, _ := os.Open("urls.txt")
	data, _ := ioutil.ReadAll(file)
	strData := string(data)
	lines := strings.Split(strData, "\n")
	return lines[r1.Intn(len(lines))]
}

func TestAggregator(t *testing.T) {
	root := new(Resource)
	root.line = "/"
	root.prefix = ""
	root.childrenMap = make(map[string]*Resource)
	i := 0
	var lines [6]string
	lines[0] = "/"
	lines[1] = "/api"
	lines[2] = "/api/books"
	lines[3] = "/api/books/{book_id}"
	lines[4] = "/api/books/{book_id}/authors"
	lines[5] = "/api/authors"
	bodyStr := "{\"key_1\": \"value_1\"}"
	headers := make(map[string]string)
	headers["Content-Type"] = "applicatio/json"
	queryParams := make(map[string]string)
	queryParams["is_item"] = "True"
	for i < len(lines) {
		fmt.Println("parsing line: ", lines[i])
		if lines[i] != "/" {
			root.addChild(lines[i], "", "GET", "{\"success\": true}", 200, headers, queryParams)
			root.addChild(lines[i], bodyStr, "POST", "{\"success\": true}", 201, headers, queryParams)
		}
		i++
	}
	child1, ok := root.childrenMap["api"]
	if !ok {
		t.Errorf("api Resource not found")
	}
	child2, ok := child1.childrenMap["books"]
	if !ok {
		t.Errorf("api/books Resource not found")
	}
	child3, ok := child2.childrenMap["{book_id}"]
	if !ok {
		t.Errorf("api/books Resource not found")
	}
	if !child3.isItem {
		t.Errorf("api/books/{book_id} Resource isItem is False")
	}
	child4, ok := child3.childrenMap["authors"]
	if !ok {
		t.Errorf("api/books/{book_id}/authors Resource not found")
	}
	if child4.isItem {
		t.Errorf("authors is not item")
	}
	pBody, ok := child4.methodToBody["POST"]
	if !ok {
		t.Errorf("no body found for POST method")
	}
	if pBody != bodyStr {
		t.Errorf("POST body does not match")
	}
	gBody, ok := child4.methodToBody["GET"]
	if !ok {
		t.Errorf("no body found for GET method")
	}
	if gBody != "" {
		t.Errorf("GET body does not match")
	}
	fmt.Println("**************************************************")
	// root.buildTreeDescription()
	// root.GenerateMap()

	parser := new(parsers.RamlParser)
	parser.Title = "First API"
	parser.BaseURI = "http://myfirstraml.com"
	parser.Version = "v1"
	f, _ := os.Create("/home/maged/test_raml.yml")
	root.GenerateSpecs(parser, f)
}
