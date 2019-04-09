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
