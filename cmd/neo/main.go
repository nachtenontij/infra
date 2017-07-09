package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nachtenontij/infra/member"
	"gopkg.in/yaml.v2"
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
	Url string
}

func (p *Program) Run(args []string) {
	fs := flag.NewFlagSet("neo", flag.ExitOnError)
	fs.StringVar(&p.Url, "url", "https://nachtenontij.nl", "Url of ontij")
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

type Enlist struct {
	Program

	File string
}

func (c *Enlist) Run(args []string) {
	fs := flag.NewFlagSet("enlist", flag.ExitOnError)
	fs.StringVar(&c.File, "file", "", "Read request from given file.")
	fs.Parse(args)

	var req member.EnlistRequest

	// let the user fill the struct
	if !FillStruct(&req) {
		os.Exit(2)
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.PostForm(c.Program.Url+"/api/enlist",
		url.Values{"request": {string(data)}})

	if err != nil {
		log.Fatalf("failed to POST: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response: %s", err)
	}

	fmt.Printf("response: %s\n", body)
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
	fmt.Printf("\n%s\n", message)
	reply := ""
	for {
		fmt.Print("(yes, no, abort)")
		fmt.Scanln(&reply)
		switch strings.ToLower(reply) {
		case "yes", "y":
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
