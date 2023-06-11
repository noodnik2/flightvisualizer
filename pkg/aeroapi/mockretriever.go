package aeroapi

import (
	"fmt"
	"time"
)

type MockArtifactRetriever struct {
	Err                error
	Contents           []byte
	RequestedEndpoints []string
}

func (*MockArtifactRetriever) GetFlightIdsRef(tailNumber string, _ time.Time) (string, error) {
	return "/fl/" + tailNumber, nil
}

func (*MockArtifactRetriever) GetTrackForFlightRef(flightId string) string {
	return fmt.Sprintf("/fli/%s/track", flightId)
}

func (r *MockArtifactRetriever) Load(requestEndpoint string) ([]byte, error) {
	r.RequestedEndpoints = append(r.RequestedEndpoints, requestEndpoint)
	return r.Contents, r.Err
}
