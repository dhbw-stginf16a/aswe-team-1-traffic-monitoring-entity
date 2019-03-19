package trafficmonitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RegisterBody ...
type RegisterBody struct {
	Concern  string `json:"concern"`
	Endpoint string `json:"endpoint"`
}

// CentralConnector abstracts connection to central node
type CentralConnector struct {
	url string
}

// NewCentralConnector create a new CentralConnector
func NewCentralConnector(url string) *CentralConnector {
	return &CentralConnector{url}
}

// Register to central node
func (con CentralConnector) Register(hostname string, maxTries int) error {
	body := &RegisterBody{
		Concern:  "traffic",
		Endpoint: hostname }

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	for ; maxTries != 0; maxTries-- {
		resp, err := http.Post(fmt.Sprint(con.url, "/monitoring"), "application/json", bytes.NewReader(data))
		if err != nil || resp.StatusCode != 204 {
			log.Print("Registration attempt failed. Pause and retry")
			<-time.After(time.Second)
		} else {
			return nil
		}
	}
	return errors.New("Register Timed out")
}

// GetGlobalPref from central node
func (con CentralConnector) GetGlobalPref(key string) string {

	resp, err := http.Get(fmt.Sprint(con.url, "/preferences/global"))
	if err != nil || resp.StatusCode != 200 {
		return ""
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return ""
	}

	m := f.(map[string]interface{})

	val, ok := m[key]
	if !ok {
		return ""
	}
	result, ok := val.(string)
	if !ok {
		return ""
	}

	return result
}

// GetUserPref from central node
func (con CentralConnector) GetUserPref(userid string) map[string]string {
	return nil
}
