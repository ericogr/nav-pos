package v1

import (
	"log"

	"encoding/json"
	"net/http"

	"github.com/ericogr/nav-pos/internal/app"
)

func (a *HandleRequests) HandleAircraft(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleAircraft called")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bbox := app.BoundingBox{}
	err := json.NewDecoder(r.Body).Decode(&bbox)
	if err != nil {
		http.Error(w, "Erro decoding JSON", http.StatusBadRequest)
		return
	}

	airCrafts, err := a.Aircraft.GetAircrafts(bbox)
	if err != nil {
		http.Error(w, "Erro calling airCrafts api", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(airCrafts)
}

func (a *HandleRequests) HandleTelemetry(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTelemetry called")

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	telemetryData := a.Telemetry.GetTelemetryData()
	if telemetryData.GPSValid {
		json.NewEncoder(w).Encode(telemetryData)
	} else {
		w.Write([]byte("{}"))
	}
}
