package main

import (
	"log"

	"github.com/meisbokai/GolangApiTest/cmd/api/server"
)

func main() {
	// Start API server
	app, err := server.NewServerApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
