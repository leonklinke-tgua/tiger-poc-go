package main

import (
	"log"

	"github.com/theguarantors/tiger/cmd"
)

func main() {
	// config, configErr := config.LoadConfig()
	// if configErr != nil {
	// 	log.Fatal("cannot load config: ", configErr)
	// }

	server, err := cmd.InitializeService() // You need to import the InitializeService function from wire_gen.go
	if err != nil {
		log.Fatal("cannot start server: ", err)
		return
	}

	if err := server.Start(); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
