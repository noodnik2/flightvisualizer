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
	CutoffTime  time.Time
	FlightCount int
}

func (tc *TracksConverter) ConvertForTailNumber(aeroApi aeroapi.Api, tracker kml.TrackGenerator, tailNumber string) ([]*kml.Track, error) {

	flightIds, getIdsErr := aeroApi.GetFlightIds(tailNumber, tc.CutoffTime)
	if getIdsErr != nil {
		return nil, getIdsErr
	}

	var kmlTracks []*kml.Track
	nFlights := len(flightIds)
	var errorList []error
	for i := 0; i < nFlights; i++ {
		kmlTrack, convertErr := ConvertForFlightId(aeroApi, tracker, flightIds[i])
		if convertErr != nil {
			errorList = append(errorList, convertErr)
			continue
		}
		kmlTracks = append(kmlTracks, kmlTrack)
		if tc.FlightCount != 0 && len(kmlTracks) == tc.FlightCount {
			// presumes the user's preferred flights are first in the list
			break
		}
	}
	if errorList != nil {
		verboseMessagePrinter := func(mt string) {
			if tc.Verbose {
				for _, err := range errorList {
					log.Printf("%s: %s\n", mt, err)
				}
			} else {
				log.Printf("NOTE: not all tracks were generated; use 'verbose' for more detail\n")
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

func ConvertForFlightId(aeroApi aeroapi.Api, tracker kml.TrackGenerator, flightId string) (*kml.Track, error) {
	track, getTrackErr := aeroApi.GetTrackForFlightId(flightId)
	if getTrackErr != nil {
		return nil, getTrackErr
	}
	kmlTrack, kmlTrackErr := tracker.Generate(track)
	if kmlTrackErr != nil {
		return nil, kmlTrackErr
	}
	return kmlTrack, nil
}
