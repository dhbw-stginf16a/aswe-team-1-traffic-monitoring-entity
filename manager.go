package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

// Manager main application handler
type Manager struct {
	apiKey          string
	hostname        string
	centralCon      *CentralConnector
	requestEndpoint *RequestEndpoint
	distReq         *DistanceRequester
}

// NewManager creates new Manager
func NewManager(centralURL string) *Manager {
	return &Manager{centralCon: NewCentralConnector(centralURL), requestEndpoint: &RequestEndpoint{}, distReq: &DistanceRequester{}}
}

// FetchApiKey fetches google api key from env or central node
func (m *Manager) fetchAPIKey() error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return errors.New("API Key not found")
	}
	// TODO: Fetch from Pref Store

	m.apiKey = apiKey
	return nil
}

// FetchApiKey fetches google api key from env or central node
func (m *Manager) fetchHostName() error {
	hostname := os.Getenv("HOST_NAME")
	if hostname != "" {
		m.hostname = hostname
		return nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	m.hostname = hostname
	return nil
}

// Init Manager
func (m *Manager) Init() error {
	log.Println("Fetching Hostname")
	if err := m.fetchHostName(); err != nil {
		log.Println("Error fetching hostname")
		return err
	}
	log.Println("Hostname: ", m.hostname)

	log.Println("Registering at central node")
	if err := m.centralCon.Register(m.hostname); err != nil {
		log.Println("Error registering on central node")
		return err
	}
	log.Println("Registration successfull")

	log.Println("Fetching Google API key")
	if err := m.fetchAPIKey(); err != nil {
		log.Println("Google API key is missing")
		return err
	}
	log.Println("Key fetched successfully")

	log.Println("Initializing google api connection")
	if err := m.distReq.Init(m.apiKey); err != nil {
		log.Println("Error during google api initialization")
		return err
	}
	log.Println("Initialization successfull")

	log.Println("Initializing request endpoint")

	if err := m.requestEndpoint.Init(m.distReq); err != nil {
		log.Println("Error during request endpoint initialization")
		return err
	}
	log.Println("Initialization successfull")

	return nil
}

// Serve application
func (m *Manager) Serve() {
	m.requestEndpoint.StartServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// Shutdown server
func (m *Manager) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m.requestEndpoint.Shutdown(ctx)
}
