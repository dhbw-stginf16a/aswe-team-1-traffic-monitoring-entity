package main

import (
	"log"
	"os"
)

func main() {

	manager := NewManager("centralnode:8080")

	err := manager.Init()
	if err != nil {
		log.Fatal(err)
	}
	
	manager.Serve()

	manager.Shutdown()

	os.Exit(0)

}
