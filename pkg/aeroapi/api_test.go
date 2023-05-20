package aeroapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/testfixtures"
)

func TestGetFlightIdsRequestGeneration(t *testing.T) {

	var requestedEndpoints []string
	getter := func(requestEndpoint string) ([]byte, error) {
		requestedEndpoints = append(requestedEndpoints, requestEndpoint)
		return nil, nil
	}

	_, _ = (&Api{Getter: getter}).GetFlightIds("tailNumber", nil, 0)

	requirer := require.New(t)
	requirer.Equal([]string{"/flights/tailNumber"}, requestedEndpoints)

}

func TestGetFlightIdsResponseProcessing(t *testing.T) {

	type testCaseDef struct {
		name              string
		getter            GetRequester
		collector         *responseCollector
		assertions        func(*require.Assertions, *testCaseDef)
		cutoffTime        *time.Time
		flightCount       int
		expectedFlightIds []string
		expectedErrors    []string
	}

	// TODO add test cases for cutoffTime & flightCount

	testFlightId1 := newTestFlightIds(1)[0]
	testCaseTemplate := `{"flights": [{"fa_flight_id": "%s"}]}`
	testCases := []testCaseDef{
		{
			name:              "single flight from JSON",
			getter:            newTestStringGetter(fmt.Sprintf(testCaseTemplate, testFlightId1)),
			expectedFlightIds: []string{testFlightId1},
		},
		{
			name:              "multiple flights from struct",
			getter:            newTestFlightsGetter(3),
			expectedFlightIds: newTestFlightIds(3),
		},
		{
			name:              "multiple flights from struct limited to last 1",
			getter:            newTestFlightsGetter(3),
			expectedFlightIds: newTestFlightIds(1),
			flightCount:       1,
		},
		{
			name:              "empty response",
			getter:            newTestStringGetter("{}"),
			expectedFlightIds: nil,
		},
		{
			name:              "error response",
			getter:            newTestErrorGetter(),
			expectedFlightIds: nil,
			expectedErrors:    []string{"couldn't get", "test error"},
		},
		{
			name:      "saved response",
			getter:    newTestTracksGetter(6),
			collector: newTestResponseCollector(),
			assertions: func(requirer *require.Assertions, tc *testCaseDef) {
				// verify the expected saved response
				responses := *tc.collector.responses
				requirer.Equal(1, len(responses))
				requirer.Equal("/flights/irrelevant", responses[0].name)
				sum := crc32.ChecksumIEEE(responses[0].contents)
				requirer.Equal(uint32(0xff2cac5), sum)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &Api{Getter: tc.getter}
			if tc.collector != nil {
				api.Saver = tc.collector.saver
			}
			flightIds, err := api.GetFlightIds("irrelevant", tc.cutoffTime, tc.flightCount)
			requirer := require.New(t)
			if len(tc.expectedErrors) > 0 {
				requirer.Error(err)
				requirer.Nil(tc.expectedFlightIds)
				for _, eet := range tc.expectedErrors {
					requirer.Contains(err.Error(), eet)
				}
			} else {
				requirer.NoError(err)
				requirer.Equal(tc.expectedFlightIds, flightIds)
				if tc.assertions != nil {
					tc.assertions(requirer, &tc)
				}
			}
		})
	}

}

func TestGetTrackForFlightId(t *testing.T) {

	type testCaseDef struct {
		name                  string
		getter                GetRequester
		collector             *responseCollector
		expectedPositionCount int
		expectedErrors        []string
		assertions            func(*require.Assertions, *Track, *testCaseDef)
	}

	testCases := []testCaseDef{
		{
			name:                  "single empty track from JSON",
			getter:                newTestStringGetter(`{"positions": [{}]}`),
			expectedPositionCount: 1,
		},
		{
			name:                  "small realistic track from JSON",
			getter:                newTestStringGetter(testfixtures.NewMockTestAeroApiTrackResponse()),
			expectedPositionCount: 19,
		},
		{
			name:                  "multiple tracks from struct",
			getter:                newTestTracksGetter(6),
			expectedPositionCount: 6,
			assertions: func(requirer *require.Assertions, response *Track, _ *testCaseDef) {
				// check an arbitrary member of the returned collection
				requirer.Equal(newTestPosition(4), response.Positions[4])
			},
		},
		{
			name:                  "save track response",
			getter:                newTestTracksGetter(6),
			collector:             newTestResponseCollector(),
			expectedPositionCount: 6,
			assertions: func(requirer *require.Assertions, _ *Track, tc *testCaseDef) {
				// verify the expected saved response
				responses := *tc.collector.responses
				requirer.Equal(1, len(responses))
				requirer.Equal("/flights/irrelevant/track", responses[0].name)
				sum := crc32.ChecksumIEEE(responses[0].contents)
				requirer.Equal(uint32(0xff2cac5), sum)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &Api{Getter: tc.getter}
			if tc.collector != nil {
				api.Saver = tc.collector.saver
			}
			track, err := api.GetTrackForFlightId("irrelevant")
			requirer := require.New(t)
			if len(tc.expectedErrors) > 0 {
				requirer.Error(err)
				requirer.Nil(track)
				for _, eet := range tc.expectedErrors {
					requirer.Contains(err.Error(), eet)
				}
			} else {
				requirer.NoError(err)
				requirer.NotNil(track.Positions)
				requirer.Equal(tc.expectedPositionCount, len(track.Positions))
				if tc.assertions != nil {
					tc.assertions(requirer, track, &tc)
				}
			}
		})
	}
}

func newTestErrorGetter() GetRequester {
	return func(url string) ([]byte, error) {
		return nil, errors.New("test error")
	}
}

func newTestStringGetter(jsonText string) GetRequester {
	return func(string) ([]byte, error) {
		return []byte(jsonText), nil
	}
}

const testFlightIdTemplate = "N1234H-5678901234-stams-%05d"

func newTestFlightIds(n int) []string {
	var flightIds []string
	for i := 1; i <= n; i++ {
		flightIds = append(flightIds, fmt.Sprintf(testFlightIdTemplate, i))
	}
	return flightIds
}

func newTestFlightsGetter(n int) GetRequester {
	var fr FlightsResponse
	for _, id := range newTestFlightIds(n) {
		fr.Flights = append(fr.Flights, Flight{FlightId: id})
	}
	return func(string) ([]byte, error) {
		return json.Marshal(&fr)
	}
}

func newTestTracksGetter(n int) GetRequester {
	tr := Track{FlightId: fmt.Sprintf("track%d", n)}
	for i := 0; i < n; i++ {
		tr.Positions = append(tr.Positions, newTestPosition(i))
	}
	return func(string) ([]byte, error) {
		return json.Marshal(&tr)
	}
}

type savedResponse struct {
	name     string
	contents []byte
}

type responseCollector struct {
	responses *[]savedResponse
	saver     ResponseSaver
}

func newTestResponseCollector() *responseCollector {
	responses := &[]savedResponse{}
	return &responseCollector{
		responses: responses,
		saver: func(name string, contents []byte) (string, error) {
			*responses = append(*responses, savedResponse{
				name:     name,
				contents: contents,
			})
			return "", nil
		},
	}
}

func newTestPosition(offset int) Position {
	ts := newTestTime().Add(time.Duration(offset) * time.Minute)
	return Position{
		Timestamp:  ts,
		GsKnots:    float64(offset),
		Heading:    float64(offset),
		AltAglD100: float64(offset),
		Latitude:   float64(offset),
		Longitude:  float64(offset),
	}
}

func newTestTime() *time.Time {
	testTime := time.Date(2023, 5, 11, 11, 11, 23, 0, time.UTC)
	return &testTime
}
