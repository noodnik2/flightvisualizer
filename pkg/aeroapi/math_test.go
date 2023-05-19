package aeroapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetGeoGs(t *testing.T) {

	tv := func(ts string) time.Time {
		timeVal, _ := time.Parse(time.RFC3339, ts)
		return timeVal
	}

	testCases := []struct {
		name            string
		thisPosition    Position
		nextPosition    Position
		expectedGsKnots float64
	}{
		{
			name: "one",
			thisPosition: Position{
				Latitude:  37.65633,
				Longitude: -122.09545,
				Timestamp: tv("2023-05-11T23:27:29Z"),
			},
			nextPosition: Position{
				Latitude:  37.65244,
				Longitude: -122.09936,
				Timestamp: tv("2023-05-11T23:27:45Z"),
			},
			expectedGsKnots: 67.159,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			m := &Math{}
			gsKnots := m.GetGeoGsKnots(tc.thisPosition, tc.nextPosition)
			requirer.InDelta(tc.expectedGsKnots, gsKnots, 0.001)
		})
	}

}

func TestGetGeoBearing(t *testing.T) {

	testCases := []struct {
		name            string
		thisPosition    Position
		nextPosition    Position
		expectedBearing Degrees
	}{
		{
			name: "simple",
			thisPosition: Position{
				Latitude:  37.65633,
				Longitude: -122.09545,
			},
			nextPosition: Position{
				Latitude:  37.65244,
				Longitude: -122.09936,
			},
			expectedBearing: 218.5133,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			m := &Math{}
			geoBearing := m.GetGeoBearing(tc.thisPosition, tc.nextPosition)
			requirer.InDelta(f(tc.expectedBearing), f(geoBearing), 0.01)
		})
	}

}
