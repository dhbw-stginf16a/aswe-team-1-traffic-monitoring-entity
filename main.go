package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	requestEndpoint := &RequestEndpoint{}
	if err := requestEndpoint.Prepare(); err != nil {
		log.Fatal(err)
	}

	requestEndpoint.StartServe()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	requestEndpoint.Shutdown(ctx)

	os.Exit(0)

}
