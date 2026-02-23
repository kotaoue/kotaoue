package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/kotaoue/tools/fetch-bookmeter/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("expected subcommand: fetch-wish or update-readme")
	}

	switch os.Args[1] {
	case "fetch-wish":
		return service.RunFetchWish(os.Args[2:])
	case "update-readme":
		return service.RunUpdateReadme(os.Args[2:])
	default:
		return fmt.Errorf("unknown subcommand: %s", os.Args[1])
	}
}
