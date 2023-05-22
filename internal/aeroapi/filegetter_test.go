package aeroapi

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/noodnik2/flightvisualizer/internal/persistence"
//	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
//)
//
//func TestMockAeroApiGetterFactory_NewRequester(t *testing.T) {
//
//	testCases := []struct {
//		name                  string
//		flightsIdFnRef        string
//		flightsIdFileJson     string
//		expectedAssetRequests []string
//	}{
//		{
//			name:           "with artifacts folder",
//			flightsIdFnRef: "artifacts/noprefix_Request.json",
//			flightsIdFileJson: `{
//			  "flights": [
//				{"fa_flight_id": "fid_a"},
//				{"fa_flight_id": "fid_b"}
//			  ]
//			}`,
//			expectedAssetRequests: []string{
//				"artifacts/noprefix_Request.json",
//				"artifacts/fid_a_trk.json",
//				"artifacts/fid_b_trk.json",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			requirer := require.New(t)
//
//			var assetFilenames []string
//			mg := AeroApiFileGetterFactory{
//				Verbose: true,
//				FileContext: persistence.FileContext{
//					Reader: func(assetFilepath string) (response []byte, err error) {
//						assetFilenames = append(assetFilenames, assetFilepath)
//						return []byte(tc.flightsIdFileJson), nil
//					},
//				},
//			}
//			requester, newRequesterErr := mg.NewRequester(tc.flightsIdFnRef, func(fid string) string {
//				return fmt.Sprintf("artifacts/%s_trk.json", fid)
//			})
//			requirer.NoError(newRequesterErr)
//			flightsResponse, flightsRequestErr := requester("/flights/id")
//			requirer.NoError(flightsRequestErr)
//			requirer.NotEmpty(flightsResponse)
//			json, newRequesterErr := aeroapi.FlightsFromJson(flightsResponse)
//			requirer.NoError(newRequesterErr)
//			const expectedFlightIdsRequestCount = 1 // one from prep to get flight ids
//			requirer.Equal(expectedFlightIdsRequestCount, len(assetFilenames))
//			expectedFlightCount := len(tc.expectedAssetRequests) - expectedFlightIdsRequestCount
//			requirer.Equal(expectedFlightCount, len(json.Flights))
//			for _, flight := range json.Flights {
//				trackResponse, trackRequestErr := requester("/track/" + flight.FlightId)
//				requirer.NoError(trackRequestErr)
//				requirer.NotEmpty(trackResponse)
//			}
//			expectedTotalRequestCount := expectedFlightIdsRequestCount + expectedFlightCount
//			requirer.Equal(expectedTotalRequestCount, len(assetFilenames))
//			for i := 0; i < expectedTotalRequestCount; i++ {
//				requirer.Equal(tc.expectedAssetRequests[i], assetFilenames[i])
//			}
//		})
//	}
//}
