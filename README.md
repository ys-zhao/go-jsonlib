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
