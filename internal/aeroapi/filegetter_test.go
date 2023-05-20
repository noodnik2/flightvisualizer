package aeroapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

func TestMockAeroApiGetterFactory_NewRequester(t *testing.T) {

	testCases := []struct {
		name                  string
		flightsIdFilename     string
		flightsIdFileJson     string
		expectedAssetRequests []string
	}{
		{
			name:              "with file set prefix",
			flightsIdFilename: "230519155507Z-flights_flightsRequest.json",
			flightsIdFileJson: `{
			  "flights": [
				{"fa_flight_id": "flight_id_1"},
				{"fa_flight_id": "flight_id_2"},
				{"fa_flight_id": "flight_id_3"}
			  ]
			}`,
			expectedAssetRequests: []string{
				"assets/folder/230519155507Z-flights_flightsRequest.json",
				"assets/folder/230519155507Z-flights_flightsRequest.json",
				"assets/folder/230519155507Z-flights_flight_id_1_track.json",
				"assets/folder/230519155507Z-flights_flight_id_2_track.json",
				"assets/folder/230519155507Z-flights_flight_id_3_track.json",
			},
		},
		{
			name:              "without file set prefix",
			flightsIdFilename: "noprefix_Request.json",
			flightsIdFileJson: `{
			  "flights": [
				{"fa_flight_id": "fid_a"},
				{"fa_flight_id": "fid_b"}
			  ]
			}`,
			expectedAssetRequests: []string{
				"assets/folder/noprefix_Request.json",
				"assets/folder/noprefix_Request.json",
				"assets/folder/fid_a_track.json",
				"assets/folder/fid_b_track.json",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)

			var assetFilenames []string
			mg := AeroApiFileGetterFactory{
				Verbose:     false,
				AssetFolder: "assets/folder",
				AssetReader: func(assetFilename string) (response []byte, err error) {
					assetFilenames = append(assetFilenames, assetFilename)
					return []byte(tc.flightsIdFileJson), nil
				},
			}
			requester, newRequesterErr := mg.NewRequester(tc.flightsIdFilename)
			requirer.NoError(newRequesterErr)
			flightsResponse, flightsRequestErr := requester("/flights/id")
			requirer.NoError(flightsRequestErr)
			requirer.NotEmpty(flightsResponse)
			json, newRequesterErr := aeroapi.FlightsFromJson(flightsResponse)
			requirer.NoError(newRequesterErr)
			const expectedFlightIdsRequestCount = 2 // one from prep to get flight ids, one from client
			requirer.Equal(expectedFlightIdsRequestCount, len(assetFilenames))
			expectedFlightCount := len(tc.expectedAssetRequests) - expectedFlightIdsRequestCount
			requirer.Equal(expectedFlightCount, len(json.Flights))
			for _, flight := range json.Flights {
				trackResponse, trackRequestErr := requester("/track/" + flight.FlightId)
				requirer.NoError(trackRequestErr)
				requirer.NotEmpty(trackResponse)
			}
			expectedTotalRequestCount := expectedFlightIdsRequestCount + expectedFlightCount
			requirer.Equal(expectedTotalRequestCount, len(assetFilenames))
			for i := 0; i < expectedTotalRequestCount; i++ {
				requirer.Equal(tc.expectedAssetRequests[i], assetFilenames[i])
			}
		})
	}
}
