package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kotaoue/kotaoue/tools/fit/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cmd := flag.String("cmd", "", "subcommand to run (update-pedometer)")
	readmeFile := flag.String("readme", "README.md", "Path to README.md")
	flag.Parse()

	switch *cmd {
	case "update-pedometer":
		return service.RunUpdatePedometer(*readmeFile)
	default:
		return fmt.Errorf("expected -cmd flag with value: update-pedometer")
	}
}
