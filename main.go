package main

import (
	"log"

	"github.com/theguarantors/tiger/cmd"
)

func main() {
	server, err := cmd.InitializeService()
	if err != nil {
		log.Fatal("cannot start server: ", err)
		return
	}

	if err := server.Start(); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
