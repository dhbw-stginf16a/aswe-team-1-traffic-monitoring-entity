package trafficmonitor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCentralConnector(t *testing.T) {

	t.Run("Global Preferences", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Add("content-type", "application/json")
			resp.Write([]byte("{\"testprop\": \"testval\"}"))
		}))
		defer ts.Close()

		connector := NewCentralConnector(ts.URL)
		result := connector.GetGlobalPref("testprop")
		if result != "testval" {
			t.Error("Wrong result string")
		}

		result = connector.GetGlobalPref("nottestprop")
		if result != "" {
			t.Error("Non existing key returns nonempty response")
		}
	})

	t.Run("Register", func(t *testing.T) {

		counter := 2
		ts := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if counter > 0 {
				resp.WriteHeader(404)
				counter = counter - 1
			} else {
				data, err := ioutil.ReadAll(req.Body)
				if err != nil {
					resp.WriteHeader(404)
					return
				}
				log.Print("Body:   ", string(data))
				req.Body.Close()

				var body RegisterBody
				err = json.Unmarshal(data, &body)
				if err != nil {
					t.Fatal("Unmarshal Error")
					resp.WriteHeader(404)
					return
				}

				if body.Concern != "traffic" {
					resp.WriteHeader(404)
					return
				}

				if body.Endpoint != "http://hostname:8080/api/v1" {
					resp.WriteHeader(404)
					return
				}

				resp.WriteHeader(204)
			}
		}))
		defer ts.Close()

		connector := NewCentralConnector(ts.URL)

		err := connector.Register("http://hostname:8080/api/v1", 3)
		if err != nil {
			t.Fatal("Error while register")
		}

		if counter != 0 {
			t.Fatal("Not enough retries")
		}
	})
}
