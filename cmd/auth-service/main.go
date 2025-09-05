package main

import (
	"log"

	"github.com/Nezent/auth-service/cmd/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
