package main

import (
	"os"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"

	"github.com/wesleimp/rain/cmd/cli"
)

var (
	version = "0.1.0"
)

func main() {
	log.SetHandler(cliHandler.Default)

	err := cli.Execute(version, os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}