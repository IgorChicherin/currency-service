package main

import (
	"flag"
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/db"
	"github.com/IgorChicherin/currency-service/server"
	"os"
)

func main() {
	environment := flag.String("e", "develop", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.Init()
	server.Init()
}
