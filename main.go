package main

import "log"

func main() {
	// Create db connection
	// TODO: Add connection config as params
	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()

	// fmt.Println("First test...")
}
