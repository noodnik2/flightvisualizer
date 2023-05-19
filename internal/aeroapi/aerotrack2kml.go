package aeroapi

import (
	"time"

	"github.com/noodnik2/kmlflight/internal/kml"
	"github.com/noodnik2/kmlflight/pkg/aeroapi"
)

type TracksConverter struct {
	TailNumber  string
	CutoffTime  *time.Time
	FlightCount int
}

func (tc *TracksConverter) Convert(aeroApi *aeroapi.Api, tracker kml.KmlTracker) ([]*kml.KmlTrack, error) {

	flightIds, getIdsErr := aeroApi.GetFlightIds(tc.TailNumber, tc.CutoffTime, tc.FlightCount)
	if getIdsErr != nil {
		return nil, getIdsErr
	}

	var kmlTracks []*kml.KmlTrack
	for _, flightId := range flightIds {
		track, getTrackErr := aeroApi.GetTrackForFlightId(flightId)
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
