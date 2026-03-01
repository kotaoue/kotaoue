package main

import (
	"log"
	"os"

	"github.com/kotaoue/kotaoue/tools/rss-to-readme/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	return service.RunUpdateReadme(os.Args[1:])
}
