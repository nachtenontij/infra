package main

import (
	"flag"
	"github.com/nachtenontij/infra/common/server"
	"log"
)

func main() {
	flag.StringVar(&server.Settings.BindAddress,
		"bind", "127.0.0.1:8080", "Address to bind to")
	flag.StringVar(&server.Settings.DatabaseAddress,
		"dbaddr", "localhost", "Address of mongo server")
	flag.StringVar(&server.Settings.DatabaseName,
		"dbname", "neo", "Name of mongo database to use")
	flag.Parse()
	log.Fatal(server.ListenAndServe())
}
