package main

import (
	"log"
	"os"

	"github.com/dhbw-stginf16a/aswe-team-1-traffic-monitoring-entity/internal/trafficmonitor"
)

func main() {

	manager := trafficmonitor.NewManager("centralnode:8080")

	err := manager.Init()
	if err != nil {
		log.Fatal(err)
	}

	manager.Serve()

	manager.Shutdown()

	os.Exit(0)

}
