package main

import (
	"fmt"
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
		log.Fatalf("failed to get json from url '%s'. err: '%v'", url, err)
	}
	fmt.Println("Got user info.", "ID:", res.ID, "name:", res.Name)
}
