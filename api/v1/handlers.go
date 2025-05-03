package v1

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/ericogr/nav-pos/internal/app"
)

func (a *HandleRequests) HandleSession(w http.ResponseWriter, r *http.Request) {
	if !a.validateMethod(w, r, http.MethodGet) {
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
	if !a.validateAndHandleSession(w, r) {
		return
	}

	if !a.validateMethod(w, r, http.MethodPost) {
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
	if !a.validateAndHandleSession(w, r) {
		return
	}

	if !a.validateMethod(w, r, http.MethodGet) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	telemetryData := a.Telemetry.GetTelemetryData()
	json.NewEncoder(w).Encode(telemetryData)
}

func (a *HandleRequests) HandleTile(w http.ResponseWriter, r *http.Request) {
	if !a.validateMethod(w, r, http.MethodGet) {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 6 || parts[0] != "api" || parts[1] != "v1" || parts[2] != "tile" {
		http.Error(w, "Invalid path. Use /api/v1/tile/{z}/{x}/{y}.pbf", http.StatusBadRequest)
		return
	}

	z, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}

	x, err := strconv.Atoi(parts[4])
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}

	yStr := strings.TrimSuffix(parts[5], ".pbf")
	y, err := strconv.Atoi(yStr)
	if err != nil {
		http.Error(w, "Invalid z parameter", http.StatusBadRequest)
		return
	}

	tile, err := a.MapService.GetTile(x, y, z)
	if err != nil {
		http.Error(w, "Error getting tile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Length", strconv.Itoa(len(tile)))
	_, err = w.Write(tile)
	if err != nil {
		http.Error(w, "Error writing tile", http.StatusInternalServerError)
		return
	}
}

func (a *HandleRequests) validateMethod(w http.ResponseWriter, r *http.Request, httpMethod string) bool {
	valid := httpMethod == r.Method

	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{
			"error": "Method not allowed",
		}
		json.NewEncoder(w).Encode(response)
	}

	return valid
}

func (a *HandleRequests) validateAndHandleSession(w http.ResponseWriter, r *http.Request) bool {
	valid := a.validateSessionId(r)

	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"error": "Bad request",
		}
		json.NewEncoder(w).Encode(response)
	}

	return valid
}

func (a *HandleRequests) validateSessionId(r *http.Request) bool {
	sessionId := r.Header.Get(HEADER_SESSION_ID)
	if sessionId == "" {
		return false
	}

	return sessionId == strconv.Itoa(a.SessionId)
}
