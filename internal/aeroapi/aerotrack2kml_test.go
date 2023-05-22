package aeroapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/internal/kml"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

func TestConvert(t *testing.T) {

	testCases := []struct {
		name string
		TracksConverter
		kml.KmlTracker
	}{
		{
			name: "simple",
			TracksConverter: TracksConverter{
				TailNumber:  "tail#",
				CutoffTime:  nil,
				FlightCount: 1,
				Api: &aeroapi.RetrieverSaverApiImpl{
					Retriever: &aeroapi.MockArtifactRetriever{
						Contents: []byte(`{"flights": [{}], "positions": [{}]}`),
					},
					Saver: nil,
				},
			},
			KmlTracker: &TestKmlTracker{
				KmlTrack: kml.KmlTrack{
					KmlDoc:    []byte{},
					KmlAssets: make(map[string]any),
					StartTime: &time.Time{},
					EndTime:   &time.Time{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			convert, err := tc.Convert(tc.KmlTracker)
			requirer.NoError(err)
			requirer.NotNil(convert)
			requirer.Equal(1, len(convert))
		})
	}

}

type TestKmlTracker struct {
	kml.KmlTrack
}

func (tkt *TestKmlTracker) Generate(*aeroapi.Track) (*kml.KmlTrack, error) {
	return &tkt.KmlTrack, nil
}
