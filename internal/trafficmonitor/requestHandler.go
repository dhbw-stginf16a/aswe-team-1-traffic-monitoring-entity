package trafficmonitor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"googlemaps.github.io/maps"
)

// Request ...
type Request struct {
	Type    string         `json:"type"`
	Payload RequestPayload `json:"payload"`
}

// RequestPayload ...
type RequestPayload struct {
	Location    string      `json:"location"`
	Destination string      `json:"destination"`
	ArriveBy    string      `json:"arriveby"`
	TravelMode  []maps.Mode `json:"travelmode"`
}

// Response ...
type Response struct {
	Type    string          `json:"type"`
	Payload ResponsePayload `json:"payload"`
}

// ResponsePayload ...
type ResponsePayload struct {
	Routes []DistanceInfo `json:"routes"`
}

// RequestHandler ...
type RequestHandler struct {
	distanceRequester *DistanceRequester
}

// Init ...
func (handler *RequestHandler) Init(distReq *DistanceRequester) error {
	handler.distanceRequester = distReq
	return nil
}

// ServeHTTP ...
func (handler RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body ", err)
		http.Error(w, "Internal Error", 500)
		return
	}
	req := &Request{}
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Println("Unable to parse body", err)
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	switch req.Type {
	case "traffic_route":
		handler.ServeTrafficRoute(w, r, req)
		break
	}

}

// ServeTrafficRoute ...
func (handler RequestHandler) ServeTrafficRoute(w http.ResponseWriter, r *http.Request, request *Request) {

	resultArray := make([]DistanceInfo, 0, len(request.Payload.TravelMode))

	for _, travelMode := range request.Payload.TravelMode {
		info, found, err := handler.distanceRequester.GetDistance(request.Payload.Location, request.Payload.Destination, request.Payload.ArriveBy, travelMode)
		if err != nil {
			log.Println("Error getting distance data", err)
			http.Error(w, "Internal Error", 500)
			return
		}
		if !found {
			log.Println("Destination or Location invalid")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		resultArray = append(resultArray, *info)
	}

	response := Response{Type: "traffic_route", Payload: ResponsePayload{Routes: resultArray}}

	wrapper := make([]Response, 0, 1)
	wrapper[0] = response

	data, err := json.Marshal(wrapper)
	if err != nil {
		log.Println("Error parsing response into json", err)
		http.Error(w, "Internal Error", 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(data)
}
