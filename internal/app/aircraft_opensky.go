package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type AircraftOpenSky struct {
	Username string
	Password string
}

func (a *AircraftOpenSky) GetAircrafts(bbox BoundingBox) ([]AircraftData, error) {
	// Se temos uma área definida, adicionamos os parâmetros de bounding box
	baseURL := OPENSKY_NETWORK_API_STATES_ALL

	params := url.Values{}
	params.Add("lamin", strconv.FormatFloat(bbox.MinLatitude, 'f', -1, 64))
	params.Add("lomin", strconv.FormatFloat(bbox.MinLongitude, 'f', -1, 64))
	params.Add("lamax", strconv.FormatFloat(bbox.MaxLatitude, 'f', -1, 64))
	params.Add("lomax", strconv.FormatFloat(bbox.MaxLongitude, 'f', -1, 64))
	baseURL = baseURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	if a.Username != "" && a.Password != "" {
		req.SetBasicAuth(a.Username, a.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	a.logRemainingRequests(resp.Header)

	// Ler o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decodificar a resposta JSON
	var openSkyResp OpenSkyResponse
	err = json.Unmarshal(body, &openSkyResp)
	if err != nil {
		return nil, err
	}

	// Converter os dados para o formato de aeronave
	aircrafts := a.convertToAircraft(openSkyResp.States)
	return aircrafts, nil
}

func (a *AircraftOpenSky) logRemainingRequests(resp http.Header) {
	remaining := resp.Get(HEADER_X_RATE_LIMIT_REMAINING)

	if remaining != "" {
		if remainingInt, err := strconv.Atoi(remaining); err == nil {
			log.Printf("Requests para %s restantes: %d\n", OPENSKY_NETWORK_API_STATES_ALL, remainingInt)
		}
	} else {
		log.Printf("Cabeçalho %s não encontrado", HEADER_X_RATE_LIMIT_REMAINING)
	}
}

func (a *AircraftOpenSky) convertToAircraft(states [][]interface{}) []AircraftData {
	var aircrafts []AircraftData

	for _, state := range states {
		if len(state) < 17 {
			continue
		}

		aircraft := AircraftData{}

		// Extrair os campos do array de estados
		if state[0] != nil {
			aircraft.Icao24 = fmt.Sprintf("%v", state[0])
		}
		if state[1] != nil {
			aircraft.Callsign = fmt.Sprintf("%v", state[1])
		}
		if state[2] != nil {
			aircraft.OriginCountry = fmt.Sprintf("%v", state[2])
		}
		if state[3] != nil {
			timePos, ok := state[3].(float64)
			if ok {
				aircraft.TimePosition = int64(timePos)
			}
		}
		if state[4] != nil {
			lastCont, ok := state[4].(float64)
			if ok {
				aircraft.LastContact = int64(lastCont)
			}
		}
		if state[5] != nil {
			lon, ok := state[5].(float64)
			if ok {
				aircraft.Longitude = lon
			}
		}
		if state[6] != nil {
			lat, ok := state[6].(float64)
			if ok {
				aircraft.Latitude = lat
			}
		}
		if state[7] != nil {
			alt, ok := state[7].(float64)
			if ok {
				aircraft.BaroAltitude = alt
			}
		}
		if state[8] != nil {
			onGround, ok := state[8].(bool)
			if ok {
				aircraft.OnGround = onGround
			}
		}
		if state[9] != nil {
			vel, ok := state[9].(float64)
			if ok {
				aircraft.Velocity = vel
			}
		}
		if state[10] != nil {
			track, ok := state[10].(float64)
			if ok {
				aircraft.TrueTrack = track
			}
		}
		if state[11] != nil {
			vRate, ok := state[11].(float64)
			if ok {
				aircraft.VerticalRate = vRate
			}
		}
		// Ignoramos o campo sensors por enquanto
		if state[13] != nil {
			geoAlt, ok := state[13].(float64)
			if ok {
				aircraft.GeoAltitude = geoAlt
			}
		}
		if state[14] != nil {
			aircraft.Squawk = fmt.Sprintf("%v", state[14])
		}
		if state[15] != nil {
			spi, ok := state[15].(bool)
			if ok {
				aircraft.Spi = spi
			}
		}
		if state[16] != nil {
			posSource, ok := state[16].(float64)
			if ok {
				aircraft.PositionSource = int(posSource)
			}
		}
		if len(state) > 17 && state[17] != nil {
			cat, ok := state[17].(float64)
			if ok {
				aircraft.Category = int(cat)
			}
		}

		// Só adiciona aeronaves com coordenadas válidas
		if aircraft.Latitude != 0 && aircraft.Longitude != 0 {
			aircrafts = append(aircrafts, aircraft)
		}
	}

	return aircrafts
}
