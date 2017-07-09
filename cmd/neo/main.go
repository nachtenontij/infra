package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nachtenontij/infra/member"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func main() {
	(&Program{}).Run(os.Args[1:])
}

type Program struct {
	Url     string
	Session string
}

func (p *Program) Run(args []string) {
	fs := flag.NewFlagSet("neo", flag.ExitOnError)
	fs.StringVar(&p.Url, "url", "https://nachtenontij.nl",
		"Url of ontijd")
	fs.StringVar(&p.Session, "session", "",
		"Sets session key manually - for bootstrapping")
	fs.Parse(args)

	args = fs.Args()

	// TODO: autogenerate
	if len(args) == 0 {
		fmt.Println("subcommands: enlist")
		os.Exit(2)
	}

	switch args[0] {
	case "enlist":
		(&Enlist{Program: *p}).Run(args[1:])
	default:
		fmt.Printf("%s is not a valid command", os.Args[1])
	}
}

func (p *Program) Authorization() string {
	return p.Session
}

func (p *Program) Request(method, name string, useQueryString bool,
	request, response interface{}) (err error) {

	var body io.Reader
	u := p.Url + "/api/" + name

	reqdata, err := json.Marshal(request)
	if err != nil {
		return
	}
	query := url.Values{"request": {string(reqdata)}}.Encode()

	if useQueryString {
		u = u + "?" + query
	} else {
		body = strings.NewReader(query)
	}

	httpreq, err := http.NewRequest(method, u, body)
	if err != nil {
		return
	}

	auth := p.Authorization()
	if auth != "" {
		httpreq.Header.Set("Authorization", "basic "+auth)
	}

	if !useQueryString {
		httpreq.Header.Set("Content-Type",
			"application/x-www-form-urlencoded")
	}

	httpresp, err := http.DefaultClient.Do(httpreq)
	if err != nil {
		return
	}

	defer httpresp.Body.Close()
	data, err := ioutil.ReadAll(httpresp.Body)

	err = json.Unmarshal(data, response)
	if err != nil {
		return fmt.Errorf("could not unmarshal %s: %s",
			string(data), err)
	}

	return nil
}

type Enlist struct {
	Program

	File string
}

func (c *Enlist) Run(args []string) {
	fs := flag.NewFlagSet("enlist", flag.ExitOnError)
	// no flags yet
	fs.Parse(args)

	var req member.EnlistRequest
	var resp member.EnlistResponse

	// let the user fill the struct
	if !FillStruct(&req) {
		os.Exit(2)
	}

	err := c.Program.Request("POST", "enlist", false, req, &resp)
	if err != nil {
		log.Fatalf("request failed: %s", err)
	}

	fmt.Printf("response: %s\n", resp)
}

func FillStruct(obj interface{}) bool {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalf("could not create temporary file: %s", err)
	}
	defer os.Remove(tmpfile.Name())

	ground, err := yaml.Marshal(obj)
	if err != nil {
		log.Fatalf("could not put %s into yaml: %s\n", obj, err)
	}

	_, err = tmpfile.Write(ground)
	if err != nil {
		log.Fatalf("could not write to temporary file: %s", err)
	}

	for {
		err = EditFile(tmpfile.Name())
		if err != nil {
			log.Fatalf("failed to start vim: %s", err)
		}

		excited, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			log.Fatalf("failed to read temporary file: %s", err)
		}

		if reflect.DeepEqual(excited, ground) {
			fmt.Printf("no changes - aborting\n")
			os.Exit(2)
		}

		fmt.Println()
		fmt.Print(string(excited))
		fmt.Println()

		err = yaml.Unmarshal(excited, obj)
		if err != nil {
			fmt.Printf("ERROR: Could not parse YAML:\n%s\n", err)
			Confirm("Retry?")
			continue
		}

		if Choose("Commit?") {
			return true
		}
	}
}

func EditFile(name string) error {
	cmd := exec.Command("vim", "-c", "set syntax=yaml", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func Confirm(message string) {
	fmt.Printf("\n%s (Press control-c to abort.)", message)
	fmt.Scanln()
}

func Choose(message string) bool {
	fmt.Printf("\n%s ", message)
	reply := ""
	for {
		fmt.Print("([Yes], No, Abort) ")
		fmt.Scanln(&reply)
		switch strings.ToLower(reply) {
		case "yes", "y", "":
			return true
		case "no", "n":
			return false
		case "abort", "a":
			os.Exit(2)
		default:
			continue
		}
	}

}
