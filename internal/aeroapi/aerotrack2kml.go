package aeroapi

import (
    "log"
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
    nFlights := len(flightIds)
    for i := 0; i < nFlights; i++ {
        track, getTrackErr := tc.Api.GetTrackForFlightId(flightIds[i])
        if getTrackErr != nil {
            if i > 0 && tc.FlightCount == 0 {
                // Logic:
                // 1. We already have at least one track generated;
                // 2. The user didn't specify how many flight tracks they want or expect;
                // 3. This could be a "from artifact" (replay) scenario;
                // So, just stop the loop at the first error, but don't fail the method
                log.Printf("NOTE: %d/%d flight(s) converted; error ignored: %v", i, nFlights, getTrackErr)
                break
            }
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
