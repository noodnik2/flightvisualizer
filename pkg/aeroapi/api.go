package aeroapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/noodnik2/flightvisualizer/pkg/persistence"
)

type ResponseSaver func(string, []byte) (string, error)

type FlightsResponse struct {
	Flights []Flight `json:"flights"`
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

type Api interface {
	GetFlightIds(tailNumber string, cutoffTime time.Time) ([]string, error)
	GetTrackForFlightId(flightId string) (*Track, error)
}

type ArtifactLocator interface {
	// GetFlightIdsRef returns a reference enabling retrieval of the flight identifier(s) for the desired track(s).
	// A return value enclosed by square brackets can be interpreted as a comma-separated-list of flight identifier(s);
	// otherwise it should be interpreted as a URI to be invoked in order to obtain the desired list.
	GetFlightIdsRef(tailNumber string, cutoffTime time.Time) string
	// GetTrackForFlightUri returns a URI to be invoked in order to obtain the desired track data.
	GetTrackForFlightUri(flightId string) string
}

type ArtifactRetriever interface {
	ArtifactLocator
	persistence.Loader
}

type ArtifactSaver interface {
	ArtifactLocator
	persistence.Saver
}

type RetrieverSaverApiImpl struct {
	Retriever ArtifactRetriever
	Saver     ArtifactSaver
}

// GetFlightIds returns the AeroAPI identifier(s) of the flight(s) specified by the parameters
// cutoffTime (optional) - most recent time for a flight to be considered
func (a *RetrieverSaverApiImpl) GetFlightIds(tailNumber string, cutoffTime time.Time) ([]string, error) {
	endpoint := a.Retriever.GetFlightIdsRef(tailNumber, cutoffTime)
	if strings.HasPrefix(endpoint, "[") && strings.HasSuffix(endpoint, "]") {
		// see contract of 'GetFlightIdsRef' about "A return value enclosed by square brackets"...
		return strings.Split(endpoint[1:len(endpoint)-1], ","), nil
	}
	responseBytes, getErr := a.Retriever.Load(endpoint)
	if getErr != nil {
		return nil, newFlightApiError("get", endpoint, getErr)
	}

	if a.Saver != nil {
		saveUri := a.Saver.GetFlightIdsRef(tailNumber, cutoffTime)
		if getSaveErr := a.Saver.Save(saveUri, responseBytes); getSaveErr != nil {
			return nil, newFlightApiError("save get flight ids response", endpoint, getSaveErr)
		}
	}

	flights, flightsErr := FlightsFromJson(responseBytes)
	if flightsErr != nil {
		return nil, newFlightApiError("unmarshal", endpoint, flightsErr)
	}

	var flightIds []string
	for _, flight := range flights.Flights {
		flightIds = append(flightIds, flight.FlightId)
	}
	return flightIds, nil
}

// GetTrackForFlightId retrieves the track for the given flight given its AeroAPI identifier
func (a *RetrieverSaverApiImpl) GetTrackForFlightId(flightId string) (*Track, error) {
	endpoint := a.Retriever.GetTrackForFlightUri(flightId)
	responseBytes, getErr := a.Retriever.Load(endpoint)
	if getErr != nil {
		return nil, newFlightApiError("get", endpoint, getErr)
	}

	if a.Saver != nil {
		saveUri := a.Saver.GetTrackForFlightUri(flightId)
		if getSaveErr := a.Saver.Save(saveUri, responseBytes); getSaveErr != nil {
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

func FlightsFromJson(flightsBytes []byte) (*FlightsResponse, error) {
	var flights FlightsResponse
	if unmarshallErr := json.Unmarshal(flightsBytes, &flights); unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return &flights, nil
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
