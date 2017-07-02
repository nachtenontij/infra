package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	var url string

	flag.StringVar(&url, "url", "https://nachtenontij.nl", "Url of ontij")
	flag.Parse()

	resp, err := http.Post(url+"/api/login", "application/json",
		strings.NewReader("{hi:there}"))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("response: %s\n", body)
}
