package trafficmonitor

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"googlemaps.github.io/maps"
)

type MockDistanceRequester struct {
	info  *DistanceInfo
	found bool
	err   error
}

func (r MockDistanceRequester) GetDistance(origin, destination, arriveby string, travelMode maps.Mode) (info *DistanceInfo, found bool, err error) {
	return r.info, r.found, r.err
}

func (r MockDistanceRequester) Init(apiKey string) error {
	return nil
}

func TestRequest(t *testing.T) {
	requester := &MockDistanceRequester{}
	endpoint := &RequestEndpoint{}

	err := endpoint.Init(requester)
	if err != nil {
		t.Fatal("Error endpoint init: ", err)
	}

	endpoint.StartServe()

	<-time.After(time.Second)

	t.Run("Traffic route request", func(t *testing.T) {

		requester.found = true
		requester.err = nil
		requester.info = &DistanceInfo{
			Delay:        100,
			DelayText:    "100",
			Destination:  "Hamburg",
			Location:     "Stuttgart",
			Duration:     100,
			DurationText: "100",
			Distance:     "4km",
			Link:         "Link",
			Mode:         "driving"}

		body := `
		{
			"type" : "traffic_route",
			"payload": {
				"location": "Stuttgart",
				"destination": "Hamburg",
				"travelmode": [
					"driving"
				]
			}
		}`

		resp, err := http.Post("http://localhost:8080/api/v1/request", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal("Error during request: ", err)
		}

		if resp.StatusCode != 200 {
			t.Fatal("Traffic route request failed with code: ", resp.StatusCode)
		}

	})

	t.Run("Invalid type", func(t *testing.T) {

		requester.found = true
		requester.err = nil
		requester.info = &DistanceInfo{
			Delay:        100,
			DelayText:    "100",
			Destination:  "Hamburg",
			Location:     "Stuttgart",
			Duration:     100,
			DurationText: "100",
			Distance:     "4km",
			Link:         "Link",
			Mode:         "driving"}

		body := `
		{
			"type" : "traffic_invalid_type",
			"payload": {
				"location": "Stuttgart",
				"destination": "Hamburg",
				"travelmode": [
					"driving"
				]
			}
		}`

		resp, err := http.Post("http://localhost:8080/api/v1/request", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal("Error during request: ", err)
		}

		if resp.StatusCode != 404 {
			t.Fatal("Request did not failed but returned code: ", resp.StatusCode)
		}
	})

	t.Run("Invalid Json request", func(t *testing.T) {

		requester.found = true
		requester.err = nil
		requester.info = &DistanceInfo{
			Delay:        100,
			DelayText:    "100",
			Destination:  "Hamburg",
			Location:     "Stuttgart",
			Duration:     100,
			DurationText: "100",
			Distance:     "4km",
			Link:         "Link",
			Mode:         "driving"}

		body := `
		{
			"type" : "traffic_route",
			"payload": {
				"location": "Stuttgart",
				"destination": "Hamburg",
				"travelmode": [
					"driving"
				]
			}[][]
		}`

		resp, err := http.Post("http://localhost:8080/api/v1/request", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal("Error during request: ", err)
		}

		if resp.StatusCode != 400 {
			t.Fatal("Request did not fail with 400 but code: ", resp.StatusCode)
		}

	})

	t.Run("DistanceRequester error", func(t *testing.T) {

		requester.found = true
		requester.err = errors.New("Some Error")
		requester.info = &DistanceInfo{
			Delay:        100,
			DelayText:    "100",
			Destination:  "Hamburg",
			Location:     "Stuttgart",
			Duration:     100,
			DurationText: "100",
			Distance:     "4km",
			Link:         "Link",
			Mode:         "driving"}

		body := `
		{
			"type" : "traffic_route",
			"payload": {
				"location": "Stuttgart",
				"destination": "Hamburg",
				"travelmode": [
					"driving"
				]
			}
		}`

		resp, err := http.Post("http://localhost:8080/api/v1/request", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal("Error during request: ", err)
		}

		if resp.StatusCode != 500 {
			t.Fatal("Response code is not 500 but: ", resp.StatusCode)
		}

	})

	t.Run("Distance Requester not found", func(t *testing.T) {

		requester.found = false
		requester.err = nil
		requester.info = &DistanceInfo{
			Delay:        100,
			DelayText:    "100",
			Destination:  "Hamburg",
			Location:     "Stuttgart",
			Duration:     100,
			DurationText: "100",
			Distance:     "4km",
			Link:         "Link",
			Mode:         "driving"}

		body := `
		{
			"type" : "traffic_route",
			"payload": {
				"location": "Stuttgart",
				"destination": "Hamburg",
				"travelmode": [
					"driving"
				]
			}
		}`

		resp, err := http.Post("http://localhost:8080/api/v1/request", "application/json", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatal("Error during request: ", err)
		}

		if resp.StatusCode != 404 {
			t.Fatal("Response code is not 404 but: ", resp.StatusCode)
		}

	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	endpoint.Shutdown(ctx)
}
