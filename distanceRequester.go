package main

import (
	"context"
	"errors"
	"os"
	"time"

	"googlemaps.github.io/maps"
)

// DistanceRequester ..
type DistanceRequester struct {
	client *maps.Client
}

// Prepare ...
func (dr *DistanceRequester) Prepare() error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return errors.New("Google API Key not found")
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	dr.client = c
	return nil
}

// GetDistance ...
func (dr DistanceRequester) GetDistance(origin, destination string, travelMode maps.Mode) (info *DistanceInfo, found bool, err error) {
	r := &maps.DistanceMatrixRequest{
		Origins:       []string{origin},
		Destinations:  []string{destination},
		Mode:          travelMode,
		DepartureTime: "now",
		Units:         maps.UnitsMetric,
	}

	matrix, err := dr.client.DistanceMatrix(context.Background(), r)
	if err != nil {
		return nil, false, err
	}

	if matrix.Rows[0].Elements[0].Status == "NOT_FOUND" {
		return nil, false, nil
	}

	var duration time.Duration
	var delay time.Duration = 0
	if travelMode == maps.TravelModeDriving {
		duration = matrix.Rows[0].Elements[0].DurationInTraffic
		delay = duration - matrix.Rows[0].Elements[0].Duration
	} else {
		duration = matrix.Rows[0].Elements[0].Duration
	}

	info = &DistanceInfo{
		Methode:      travelMode,
		Duration:     int64(duration.Seconds()),
		DurationText: duration.String(),
		Delay:        int64(delay.Seconds()),
		DelayText:    delay.String(),
		Destination:  matrix.DestinationAddresses[0],
		Location:     matrix.OriginAddresses[0],
		Distance:     matrix.Rows[0].Elements[0].Distance.HumanReadable,
	}

	return info, true, nil
}

// DistanceInfo ...
type DistanceInfo struct {
	Methode      maps.Mode `json:"methode"`
	Duration     int64     `json:"duration"`
	DurationText string    `json:"durationText"`
	Delay        int64     `json:"delay"`
	DelayText    string    `json:"delayText"`
	Distance     string    `json:"distance"`
	Destination  string    `json:"destination"`
	Location     string    `json:"location"`
	Link         string    `json:"link"`
}
