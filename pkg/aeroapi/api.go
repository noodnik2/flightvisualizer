package aeroapi

import (
	"encoding/json"
	"fmt"
	"time"
)

type GetRequester func(url string) ([]byte, error)
type ResponseSaver func(string, []byte) (string, error)

type FlightsResponse struct {
	Flights []Flight
}

type Flight struct {
	FlightId string `json:"fa_flight_id"`
}

type Track struct {
	FlightId  string
	Positions []Position `json:"positions"`
}

type Position struct {
	AltAglD100 float64   `json:"altitude"`    // feet / 100 (AGL)
	GsKnots    float64   `json:"groundspeed"` // knots
	Heading    float64   `json:"heading"`     // 0..359
	Latitude   float64   `json:"latitude"`    // -90..90
	Longitude  float64   `json:"longitude"`   // -180..180
	Timestamp  time.Time `json:"timestamp"`
}

type Api struct {
	Getter GetRequester
	Saver  ResponseSaver
}

// GetFlightIds returns the AeroAPI identifier(s) of the flight(s) specified by the parameters
// cutoffTime (optional) - most recent time for a flight to be considered
// flightCount (optional) - max number of (most recent) flights to consider
func (a *Api) GetFlightIds(tailNumber string, cutoffTime *time.Time, flightCount int) ([]string, error) {
	endpoint := fmt.Sprintf("/flights/%s", tailNumber)
	if cutoffTime != nil {
		endpoint += fmt.Sprintf("?&end=%s", cutoffTime.Format(time.RFC3339))
	}
	responseBytes, getErr := a.Getter(endpoint)
	if getErr != nil {
		return nil, newFlightApiError("get", endpoint, getErr)
	}

	if a.Saver != nil {
		if _, getSaveErr := a.Saver(endpoint, responseBytes); getSaveErr != nil {
			return nil, newFlightApiError("save get flight ids response", endpoint, getSaveErr)
		}
	}

	var flights FlightsResponse
	if unmarshallErr := json.Unmarshal(responseBytes, &flights); unmarshallErr != nil {
		return nil, newFlightApiError("unmarshal", endpoint, unmarshallErr)
	}

	var flightIds []string
	for i, f := range flights.Flights {
		if flightCount != 0 && i >= flightCount {
			// Assuming that first entries returned are most recent; as the AeroAPI doc says: "approximately 14 days
			// of recent and scheduled flight information is returned, ordered by scheduled_out (or scheduled_off
			// if scheduled_out is missing) descending"
			break
		}
		flightIds = append(flightIds, f.FlightId)
	}
	return flightIds, nil
}

// GetTrackForFlightId retrieves the track for the given flight given its AeroAPI identifier
func (a *Api) GetTrackForFlightId(flightId string) (*Track, error) {
	endpoint := fmt.Sprintf("/flights/%s/track", flightId)
	responseBytes, getErr := a.Getter(endpoint)
	if getErr != nil {
		return nil, newFlightApiError("get", endpoint, getErr)
	}

	if a.Saver != nil {
		if _, getSaveErr := a.Saver(endpoint, responseBytes); getSaveErr != nil {
			return nil, newFlightApiError("save get track response", endpoint, getSaveErr)
		}
	}

	track, unmarshallErr := TrackFromJson(responseBytes)
	if unmarshallErr != nil {
		return nil, newFlightApiError("unmarshal", endpoint, unmarshallErr)
	}

	track.FlightId = flightId
	return track, nil
}

func TrackFromJson(aeroApiTrackJson []byte) (*Track, error) {
	var track Track
	if unmarshallErr := json.Unmarshal(aeroApiTrackJson, &track); unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &track, nil
}

func newFlightApiError(what, where string, err error) error {
	return fmt.Errorf("couldn't %s for %s: %w", what, where, err)
}
