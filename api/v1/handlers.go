package v1

import (
	"log"
	"math/rand"
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/ericogr/nav-pos/internal/app"
)

func (a *HandleRequests) HandleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{
			"error": "Method not allowed",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/text")

	if a.SessionId == 0 {
		a.SessionId = rand.Int()
		log.Printf("Session ID: %d\n", a.SessionId)
	}

	w.Write([]byte(strconv.Itoa(a.SessionId)))
}

func (a *HandleRequests) HandleAircraft(w http.ResponseWriter, r *http.Request) {
	if !a.validateSessionId(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"error": "Bad request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{
			"error": "Method not allowed",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bbox := app.BoundingBox{}
	err := json.NewDecoder(r.Body).Decode(&bbox)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	airCrafts, err := a.Aircraft.GetAircrafts(bbox)
	if err != nil {
		http.Error(w, "Error calling aircrafts api", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(airCrafts)
}

func (a *HandleRequests) HandleTelemetry(w http.ResponseWriter, r *http.Request) {
	if !a.validateSessionId(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"error": "Bad request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{
			"error": "Method not allowed",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	telemetryData := a.Telemetry.GetTelemetryData()
	json.NewEncoder(w).Encode(telemetryData)
}

func (a *HandleRequests) validateSessionId(r *http.Request) bool {
	sessionId := r.Header.Get(HEADER_SESSION_ID)
	if sessionId == "" {
		return false
	}

	return sessionId == strconv.Itoa(a.SessionId)
}
