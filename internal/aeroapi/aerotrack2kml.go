package aeroapi

import (
	"errors"
	"log"
	"time"

	"github.com/noodnik2/flightvisualizer/internal/kml"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type TracksConverter struct {
	Verbose     bool
	TailNumber  string
	CutoffTime  time.Time
	FlightCount int
	aeroapi.Api
}

func (tc *TracksConverter) Convert(tracker kml.KmlTracker) ([]*kml.KmlTrack, error) {

	flightIds, getIdsErr := tc.Api.GetFlightIds(tc.TailNumber, tc.CutoffTime, tc.FlightCount)
	if getIdsErr != nil {
		return nil, getIdsErr
	}

	var kmlTracks []*kml.KmlTrack
	nFlights := len(flightIds)
	var errorList []error
	for i := 0; i < nFlights; i++ {
		track, getTrackErr := tc.Api.GetTrackForFlightId(flightIds[i])
		if getTrackErr != nil {
			errorList = append(errorList, getTrackErr)
			continue
		}
		kmlTrack, kmlTrackErr := tracker.Generate(track)
		if kmlTrackErr != nil {
			errorList = append(errorList, kmlTrackErr)
			continue
		}
		kmlTracks = append(kmlTracks, kmlTrack)
	}
	if errorList != nil {
		verboseMessagePrinter := func(mt string) {
			if tc.Verbose {
				for _, err := range errorList {
					log.Printf("%s: %s\n", mt, err)
				}
			} else {
				log.Printf("NOTE: not all tracks were generated; use 'verbose' for more detail")
			}
		}
		nGenerated := len(kmlTracks)
		if nGenerated == 0 {
			verboseMessagePrinter("ERROR")
			return nil, errors.New("error(s) encountered generating KML visualization(s)")
		}
		verboseMessagePrinter("INFO")
	}

	return kmlTracks, nil
}
