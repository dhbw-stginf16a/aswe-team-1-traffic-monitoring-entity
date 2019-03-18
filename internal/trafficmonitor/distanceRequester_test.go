package trafficmonitor

import (
	"testing"
	"time"

	"googlemaps.github.io/maps"
)

func TestRequestCreation(t *testing.T) {
	t.Run("Without arriveby set", func(t *testing.T) {
		r := createRequest("Stuttgart", "Stuttgart", "", maps.TravelModeTransit)
		if r.Origins[0] != "Stuttgart" {
			t.Fail()
		}
		if r.Destinations[0] != "Stuttgart" {
			t.Fail()
		}
		if r.Mode != maps.TravelModeTransit {
			t.Fail()
		}
		if r.DepartureTime != "now" {
			t.Fail()
		}
	})

	t.Run("With arriveby set", func(t *testing.T) {
		r := createRequest("Stuttgart", "Stuttgart", "1234567", maps.TravelModeTransit)
		if r.Origins[0] != "Stuttgart" {
			t.Fail()
		}
		if r.Destinations[0] != "Stuttgart" {
			t.Fail()
		}
		if r.Mode != maps.TravelModeTransit {
			t.Fail()
		}
		if r.DepartureTime != "" {
			t.Fail()
		}
		if r.ArrivalTime != "1234567" {
			t.Fail()
		}
	})
}

func TestResponseParser(t *testing.T) {

	matrix := &maps.DistanceMatrixResponse{
		DestinationAddresses: []string{"Stuttgart"},
		OriginAddresses:      []string{"Stuttgart"},
		Rows: []maps.DistanceMatrixElementsRow{
			{
				Elements: []*maps.DistanceMatrixElement{
					{
						Distance: maps.Distance{
							HumanReadable: "24km"},
						Duration:          time.Second,
						DurationInTraffic: time.Minute},
				},
			},
		},
	}

	t.Run("Travel mode driving", func(t *testing.T) {
		r := parseResult(matrix, maps.TravelModeDriving)
		if r.Duration != 60 {
			t.Fail()
		}
		if r.Delay != 59 {
			t.Fail()
		}
		if r.Destination != "Stuttgart" {
			t.Fail()
		}
		if r.Location != "Stuttgart" {
			t.Fail()
		}
		if r.Mode != maps.TravelModeDriving {
			t.Fail()
		}
		if r.Distance != "24km" {
			t.Fail()
		}
	})

	t.Run("Travel mode transit", func(t *testing.T) {
		r := parseResult(matrix, maps.TravelModeTransit)
		if r.Duration != 1 {
			t.Fail()
		}
		if r.Delay != 0 {
			t.Fail()
		}
		if r.Destination != "Stuttgart" {
			t.Fail()
		}
		if r.Location != "Stuttgart" {
			t.Fail()
		}
		if r.Mode != maps.TravelModeTransit {
			t.Fail()
		}
		if r.Distance != "24km" {
			t.Fail()
		}
	})
}
