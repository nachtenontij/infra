package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	var u string

	flag.StringVar(&u, "url", "https://nachtenontij.nl", "Url of ontij")
	flag.Parse()

	resp, err := http.PostForm(u+"/api/login",
		url.Values{"request": {"{}"}})

	if err != nil {
		log.Fatal(fmt.Errorf("failed to POST: %s", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read response: %s", err))
	}

	fmt.Printf("response: %s\n", body)
}
