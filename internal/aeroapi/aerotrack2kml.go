package aeroapi

import (
	"time"

	"github.com/noodnik2/flightvisualizer/internal/kml"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type TracksConverter struct {
	TailNumber  string
	CutoffTime  *time.Time
	FlightCount int
	aeroapi.Api
}

func (tc *TracksConverter) Convert(tracker kml.KmlTracker) ([]*kml.KmlTrack, error) {

	flightIds, getIdsErr := tc.Api.GetFlightIds(tc.TailNumber, tc.CutoffTime, tc.FlightCount)
	if getIdsErr != nil {
		return nil, getIdsErr
	}

	var kmlTracks []*kml.KmlTrack
	for _, flightId := range flightIds {
		track, getTrackErr := tc.Api.GetTrackForFlightId(flightId)
		if getTrackErr != nil {
			return nil, getTrackErr
		}
		kmlTrack, kmlTrackErr := tracker.Generate(track)
		if kmlTrackErr != nil {
			return nil, kmlTrackErr
		}
		kmlTracks = append(kmlTracks, kmlTrack)
	}

	return kmlTracks, nil
}
