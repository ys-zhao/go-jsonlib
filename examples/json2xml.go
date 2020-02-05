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
