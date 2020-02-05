# jsonlib
This is a small json library for my own projects, if you are interested in the same functionalities, grab it and help yourself. :)

# Examples

## Get JSON data
Returns a JSON object from a URL
```go
package main

import (
	"log"

	"github.com/d3sw/jsonlib"
)

func main() {
	url := "https://api.github.com/users/octocat"
	var res struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := jsonlib.GetJSON(url, nil, &res); err != nil {
		log.Fatal("failed to get json from url", err)
	}
	log.Println("Got user info.", "ID:", res.ID, "name:", res.Name)
}
```

## Post JSON data
Posts a JSON object to an URL and returns a JSON response
```go
package main

import (
	"log"

	"github.com/d3sw/jsonlib"
)

func main() {
	url := "https://api.github.com/users/octocat"
	req := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{1234, "name1234"}
	var res interface{}
	if err := jsonlib.PostJSON(url, nil, &req, &res); err != nil {
		log.Fatal("failed to post to server.", err)
	}
	log.Println("Post never succeed.")
}
```

## Converts JSON to XML
A single API to convert json to xml
```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ys-zhao/jsonlib"
)

func main() {
	url := "https://api.github.com/users/octocat"
	res, err := http.Get(url)
	if err != nil {
		fmt.Fatalf("main: failed to get json from '%s'", url)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	xmlStr, _ := jsonlib.JSON2XML(string(data), jsonlib.J2XWithRootTag("root"), jsonlib.J2XWithIndent(true, "", "  "))
	fmt.Printf("main: json2xml from '%s'\n", url)
	fmt.Println(xmlStr)
}
```

## Converts xml to json
A single API to convert json to xml
```go
package main

import (
	"fmt"

	"github.com/ys-zhao/jsonlib"
)

var xmlStr = `<root><name>foo</name><age>21</age></root>`
func main() {
	jsonStr, _ := jsonlib.XML2JSON(xmlStr, jsonlib.X2JWithOmitRoot(false), jsonlib.X2JWithIndent(true, "", "  "))
	fmt.Println("main: json2xml with root node...")
	fmt.Println(jsonStr)

	jsonStr, _ = jsonlib.XML2JSON(xmlStr, jsonlib.X2JWithOmitRoot(true), jsonlib.X2JWithIndent(true, "", "  "))
	fmt.Println("main: json2xml without root node...")
	fmt.Println(jsonStr)
}
```