package kml

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/internal/kml/builders"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
	"github.com/noodnik2/flightvisualizer/testfixtures"
)

func TestNewKmlTrack(t *testing.T) {

	type testCaseDef struct {
		tracker        TrackGenerator
		input          *aeroapi.Track
		expectAssets   bool
		expectedErrors []string
	}

	testCases := []testCaseDef{
		{
			tracker: &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.PlacemarkBuilder{}}},
			input:   newMockTestAeroApiTrack(),
		},
		{
			tracker: &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.CameraBuilder{}}},
			input:   newMockTestAeroApiTrack(),
		},
		{
			tracker: &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.PathBuilder{Color: color.RGBA{R: 217, G: 51, B: 255}}}},
			input:   newMockTestAeroApiTrack(),
		},
		{
			tracker:      &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.VectorBuilder{}}},
			input:        newMockTestAeroApiTrack(),
			expectAssets: true,
		},
		{
			tracker:        &TrackBuilderEnsemble{},
			input:          &aeroapi.Track{FlightId: "xyz321"},
			expectedErrors: []string{"xyz321", "no builders"},
		},
		{
			tracker:        &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.PlacemarkBuilder{}}},
			input:          &aeroapi.Track{FlightId: "abc123"},
			expectedErrors: []string{"abc123", "no positions"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.tracker), func(t *testing.T) {
			requirer := require.New(t)
			kmlTrack, newKmlTrackErr := tc.tracker.Generate(tc.input)
			if tc.expectedErrors != nil {
				requirer.Error(newKmlTrackErr)
				for _, expectedErr := range tc.expectedErrors {
					requirer.Contains(newKmlTrackErr.Error(), expectedErr)
				}
				return
			}
			requirer.NoError(newKmlTrackErr)
			requirer.NotEmpty(kmlTrack)
			requirer.NotNil(kmlTrack.StartTime)
			requirer.NotNil(kmlTrack.EndTime)
			if tc.expectAssets {
				requirer.NotEmpty(kmlTrack.KmlAssets)
			} else {
				requirer.Empty(kmlTrack.KmlAssets)
			}
		})
	}

}

func newMockTestAeroApiTrack() *aeroapi.Track {
	track, _ := aeroapi.TrackFromJson([]byte(testfixtures.NewMockTestAeroApiTrackResponse()))
	return track
}
