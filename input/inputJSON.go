package input

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	//Metric represents a single time series point
	Metric struct {
		Name      string
		Host      string
		Timestamp string
		Type      string
		Value     float64
		Sampling  float64
		Tags      map[string]string
	}

	//Event represents a single event instance
	Event struct {
		Name           string
		Text           string
		Host           string
		AggregationKey string
		Priority       string
		Timestamp      string
		AlertType      string
		Tags           map[string]string
		SourceType     string
	}

	metricsHTTPHandler struct {
		metricsIn chan Metric
	}

	eventsHTTPHandler struct {
		eventsIn chan Event
	}
)

//http handler function, unmarshalls json encoded metric into metric struct
func (handler *metricsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var receivedMetric Metric
	err := decoder.Decode(&receivedMetric)

	if err == nil {
		//add an aditional tag specifing the host which forwarded aggregateD the metric
		//this might often be the same as the client specified host field but in situations
		//where the client is behind NAT, i.e many EVE clients this information is useful.
		if receivedMetric.Tags == nil {
			receivedMetric.Tags = make(map[string]string)
		}
		receivedMetric.Tags["Source"] = r.RemoteAddr

		handler.metricsIn <- receivedMetric
	} else {
		fmt.Println("error parsing metric")
		fmt.Println(err)
	}

	r.Body.Close()
}

//unmarshall json encoded events into event struct
func (handler *eventsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var receivedEvent Event
	err := decoder.Decode(&receivedEvent)
	fmt.Println(receivedEvent.Text)
	if err == nil {
		if receivedEvent.Tags == nil {
			receivedEvent.Tags = make(map[string]string)
		}
		receivedEvent.Tags["Source"] = r.RemoteAddr
		handler.eventsIn <- receivedEvent

	} else {
		fmt.Println("error parsing event")
		fmt.Println(err)
	}

	r.Body.Close()
}

//ServeHTTP exposes /events and /metrics and proceses JSON encoded events
func ServeHTTP(port string, metricsIn chan Metric, eventsIn chan Event) {
	server := http.NewServeMux()

	metricsHandler := new(metricsHTTPHandler)
	metricsHandler.metricsIn = metricsIn

	eventsHandler := new(eventsHTTPHandler)
	eventsHandler.eventsIn = eventsIn

	server.Handle("/metrics", metricsHandler)
	server.Handle("/events", eventsHandler)

	log.Fatal(http.ListenAndServe(":"+port, server))
}
