//+build !test

package trafficmonitor

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
	centralURL      string
	centralCon      *CentralConnector
	requestEndpoint *RequestEndpoint
	distReq         DistanceRequester
}

// NewManager creates new Manager
func NewManager() *Manager {
	return &Manager{requestEndpoint: &RequestEndpoint{}, distReq: &GoogleDistanceRequester{}}
}

// FetchApiKey fetches google api key from env or central node
func (m *Manager) fetchAPIKey() error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey != "" {
		m.apiKey = apiKey
		return nil
	}

	apiKey = m.centralCon.GetGlobalPref("traffic/apiKey")
	if apiKey == "" {
		return errors.New("Google API key not found")
	}

	m.apiKey = apiKey
	return nil
}

// Fetch Own Url
func (m *Manager) fetchHostName() error {
	hostname := os.Getenv("OWN_URL")
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

// Fetch CentralNode URL
func (m *Manager) fetchCentralNodeURL() error {
	centralURL := os.Getenv("CENTRAL_URL")
	if centralURL == "" {
		return errors.New("Central URL not found")
	}
	m.centralURL = centralURL
	return nil
}

// Init Manager
func (m *Manager) Init() error {
	log.Println("Fetching Hostname")
	if err := m.fetchHostName(); err != nil {
		log.Println("Error fetching hostname")
		return err
	}
	log.Println("Own URL: ", m.hostname)

	log.Println("Fetching Central Node url")
	if err := m.fetchCentralNodeURL(); err != nil {
		log.Println("Error fetching hostname")
		return err
	}
	m.centralCon = NewCentralConnector(m.centralURL)
	log.Println("Central Node URL: ", m.centralURL)

	log.Println("Registering at central node")
	if err := m.centralCon.Register(m.hostname, -1); err != nil {
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
