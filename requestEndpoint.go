package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// RequestEndpoint ...
type RequestEndpoint struct {
	server *http.Server
	router *mux.Router
}

// Prepare ...
func (endpoint *RequestEndpoint) Init(distReq *DistanceRequester) error {
	endpoint.router = mux.NewRouter()

	requestHandler := &RequestHandler{}
	if err := requestHandler.Init(distReq); err != nil {
		return err
	}

	endpoint.router.PathPrefix("/api/v1/request").Methods("POST").Handler(requestHandler)

	endpoint.server = &http.Server{
		Handler:      endpoint.router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return nil
}

// StartServe ...
func (endpoint *RequestEndpoint) StartServe() {
	go func() {
		if err := endpoint.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

// Shutdown ...
func (endpoint *RequestEndpoint) Shutdown(ctx context.Context) {
	endpoint.server.Shutdown(ctx)
}
