Go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Tracker struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

var trackers = []Tracker{}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/trackers", getTrackers).Methods("GET")
	r.HandleFunc("/trackers", createTracker).Methods("POST")
	r.HandleFunc("/trackers/{id}", getTracker).Methods("GET")
	r.HandleFunc("/trackers/{id}", updateTracker).Methods("PUT")
	r.HandleFunc("/trackers/{id}", deleteTracker).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getTrackers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(trackers)
}

func createTracker(w http.ResponseWriter, r *http.Request) {
	var tracker Tracker
	_ = json.NewDecoder(r.Body).Decode(&tracker)
	tracker.Timestamp = time.Now().Format(time.RFC3339)
	tracker.Status = "pending"
	trackers = append(trackers, tracker)
	json.NewEncoder(w).Encode(tracker)
}

func getTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, tracker := range trackers {
		if tracker.ID == params["id"] {
			json.NewEncoder(w).Encode(tracker)
			return
		}
	}
	json.NewEncoder(w).Encode(&Tracker{})
}

func updateTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, tracker := range trackers {
		if tracker.ID == params["id"] {
			_ = json.NewDecoder(r.Body).Decode(&trackers[index])
			trackers[index].Timestamp = time.Now().Format(time.RFC3339)
			json.NewEncoder(w).Encode(trackers[index])
			return
		}
	}
	json.NewEncoder(w).Encode(&Tracker{})
}

func deleteTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, tracker := range trackers {
		if tracker.ID == params["id"] {
			trackers = append(trackers[:index], trackers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(&Tracker{})
}

func trackURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, tracker := range trackers {
		if tracker.ID == params["id"] {
			resp, err := http.Get(tracker.URL)
			if err != nil {
				tracker.Status = "error"
			} else {
				tracker.Status = "ok"
			}
			trackers = append(trackers[:index], tracker)
			break
		}
	}
}

func automateTracking() {
	for range time.Tick(time.Second * 10) {
		for _, tracker := range trackers {
			resp, err := http.Get(tracker.URL)
			if err != nil {
				tracker.Status = "error"
			} else {
				tracker.Status = "ok"
			}
			tracker.Timestamp = time.Now().Format(time.RFC3339)
		}
	}
}

func init() {
	go automateTracking()
}