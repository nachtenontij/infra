package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nachtenontij/infra/member"
	"golang.org/x/crypto/ssh/terminal"
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
		fmt.Println("subcommands: enlist, su, passwd, login")
		os.Exit(2)
	}

	switch args[0] {
	case "enlist":
		p.Enlist(args[1:])
	case "su":
		p.SelectUser(args[1:])
	case "passwd":
		p.Passwd(args[1:])
	case "login":
		p.Login(args[1:])
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

	if httpresp.StatusCode != 200 {
		return fmt.Errorf("response %s: %s",
			httpresp.Status, string(data))
	}

	err = json.Unmarshal(data, response)
	if err != nil {
		return fmt.Errorf("could not unmarshal %s: %s",
			string(data), err)
	}

	return nil
}

func (p *Program) Enlist(args []string) {
	fs := flag.NewFlagSet("enlist", flag.ExitOnError)
	// no flags yet
	fs.Parse(args)

	var req member.EnlistRequest
	var resp member.EnlistResponse

	for {
		// let the user fill the struct
		if !FillStruct(&req) {
			os.Exit(2)
		}

		err := p.Request("POST", "enlist", false, req, &resp)
		if err == nil {
			break
		}

		fmt.Printf("request failed: %s\n", err)
		Confirm("Retry?")
	}

	fmt.Printf("response: %s\n", resp)
}

func (p *Program) SelectUser(args []string) {
	fs := flag.NewFlagSet("su", flag.ExitOnError)
	// no flags yet
	fs.Parse(args)

	args = fs.Args()
	if len(args) != 1 {
		fmt.Println("usage: neo su <handle>")
		os.Exit(2)
	}

	req := member.SelectUserRequest{Handle: args[0]}
	var resp member.SelectUserResponse

	err := p.Request("POST", "su", false, req, &resp)
	if err != nil {
		log.Fatalf("request failed: %s\n", err)
	}

	fmt.Printf("response: %s\n", resp)
}

func (p *Program) Passwd(args []string) {
	password := Password("new password: ")

	req := member.PasswdRequest{Password: password}
	var resp member.PasswdResponse

	err := p.Request("POST", "passwd", false, req, &resp)
	if err != nil {
		log.Fatalf("request failed: %s\n", err)
	}

	fmt.Printf("response: %s\n", resp)
}

func (p *Program) Login(args []string) {
	if len(args) != 1 {
		fmt.Println("usage: neo login <username>")
		os.Exit(2)
	}
	req := member.LoginRequest{
		Handle:   args[0],
		Password: Password("password: "),
	}
	var resp member.LoginResponse

	err := p.Request("POST", "login", false, req, &resp)
	if err != nil {
		log.Fatalf("request failed: %s\n", err)
	}

	fmt.Printf("response: %s\n", resp)
}

func FillStruct(obj interface{}) bool {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalf("could not create temporary file: %s\n", err)
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

func Password(message string) string {
	fmt.Print(message)
	data, err := terminal.ReadPassword(0)
	fmt.Println()
	if err != nil {
		log.Fatalf("could not read password: %s\n", err)
	}
	return string(data)
}
