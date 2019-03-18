package trafficmonitor

import (
	"context"
	"time"

	"googlemaps.github.io/maps"
)

// DistanceRequester ..
type DistanceRequester struct {
	client *maps.Client
}

// Init ...
func (dr *DistanceRequester) Init(apiKey string) error {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	dr.client = c
	return nil
}

func createRequest(origin, destination, arriveby string, travelMode maps.Mode) *maps.DistanceMatrixRequest {
	r := &maps.DistanceMatrixRequest{
		Origins:       []string{origin},
		Destinations:  []string{destination},
		Mode:          travelMode,
		DepartureTime: "now",
		Units:         maps.UnitsMetric}

	if arriveby != "" {
		r.DepartureTime = ""
		r.ArrivalTime = arriveby
	}

	return r
}

func parseResult(matrix *maps.DistanceMatrixResponse, travelMode maps.Mode) *DistanceInfo {
	var duration time.Duration
	var delay time.Duration
	if travelMode == maps.TravelModeDriving {
		duration = matrix.Rows[0].Elements[0].DurationInTraffic
		delay = duration - matrix.Rows[0].Elements[0].Duration
	} else {
		duration = matrix.Rows[0].Elements[0].Duration
	}

	info := &DistanceInfo{
		Mode:      travelMode,
		Duration:     int64(duration.Seconds()),
		DurationText: duration.String(),
		Delay:        int64(delay.Seconds()),
		DelayText:    delay.String(),
		Destination:  matrix.DestinationAddresses[0],
		Location:     matrix.OriginAddresses[0],
		Distance:     matrix.Rows[0].Elements[0].Distance.HumanReadable,
	}

	return info
}

// GetDistance ...
func (dr DistanceRequester) GetDistance(origin, destination, arriveby string, travelMode maps.Mode) (info *DistanceInfo, found bool, err error) {
	r := createRequest(origin, destination, arriveby, travelMode)

	matrix, err := dr.client.DistanceMatrix(context.Background(), r)
	if err != nil {
		return nil, false, err
	}

	if matrix.Rows[0].Elements[0].Status == "NOT_FOUND" {
		return nil, false, nil
	}

	info = parseResult(matrix, travelMode)

	return info, true, nil
}

// DistanceInfo ...
type DistanceInfo struct {
	Mode      maps.Mode `json:"travelmode"`
	Duration     int64     `json:"duration"`
	DurationText string    `json:"durationText"`
	Delay        int64     `json:"delay"`
	DelayText    string    `json:"delayText"`
	Distance     string    `json:"distance"`
	Destination  string    `json:"destination"`
	Location     string    `json:"location"`
	Link         string    `json:"link"`
}
