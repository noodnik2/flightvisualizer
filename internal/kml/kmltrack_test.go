package kml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/internal/kml/builders"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
	"github.com/noodnik2/flightvisualizer/testfixtures"
)

func TestNewKmlTrack(t *testing.T) {

	type testCaseDef struct {
		tracker TrackGenerator
	}

	testCases := []testCaseDef{
		{
			tracker: &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.PlacemarkBuilder{}}},
		},
		{
			tracker: &TrackBuilderEnsemble{Builders: []builders.KmlTrackBuilder{&builders.CameraBuilder{}}},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.tracker), func(t *testing.T) {
			requirer := require.New(t)
			kmlTrack, newKmlTrackErr := tc.tracker.Generate(newMockTestAeroApiTrack())
			// TODO add a meaningful test
			requirer.NoError(newKmlTrackErr)
			requirer.NotEmpty(kmlTrack)
		})
	}

}

func newMockTestAeroApiTrack() *aeroapi.Track {
	track, _ := aeroapi.TrackFromJson([]byte(testfixtures.NewMockTestAeroApiTrackResponse()))
	return track
}
