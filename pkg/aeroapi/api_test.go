package aeroapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/pkg/persistence"
	"github.com/noodnik2/flightvisualizer/testfixtures"
)

func TestGetFlightIdsRequestGeneration(t *testing.T) {

	retriever := &MockArtifactRetriever{
		Contents: []byte(`{}`),
	}
	_, _ = (&RetrieverSaverApiImpl{Retriever: retriever}).GetFlightIds("tail#", time.Time{})

	requirer := require.New(t)
	requirer.Equal([]string{"/fl/tail#"}, retriever.RequestedEndpoints)

}

func TestGetFlightIdsResponseProcessing(t *testing.T) {

	type testCaseDef struct {
		name              string
		retriever         ArtifactRetriever
		assertions        func(*require.Assertions, *testResponseSaver)
		cutoffTime        time.Time
		flightCount       int
		expectedFlightIds []string
		expectedErrors    []string
	}

	testFlightId1 := newTestFlightIds(1)[0]
	testCaseTemplate := `{"flights": [{"fa_flight_id": "%s"}]}`
	testCases := []testCaseDef{
		{
			name:              "single flight from JSON",
			retriever:         &MockArtifactRetriever{Contents: []byte(fmt.Sprintf(testCaseTemplate, testFlightId1))},
			expectedFlightIds: []string{testFlightId1},
		},
		{
			name:              "multiple flights from struct",
			retriever:         &MockArtifactRetriever{Contents: newTestFlightsJsonContents(3)},
			expectedFlightIds: newTestFlightIds(3),
		},
		{
			// this test case used to be relevant; TODO remove after confirming no longer a valid code path
			name:              "multiple flights from struct no longer respects flightCount",
			retriever:         &MockArtifactRetriever{Contents: newTestFlightsJsonContents(3)},
			expectedFlightIds: newTestFlightIds(3),
			flightCount:       1,
		},
		{
			name:              "empty response",
			retriever:         &MockArtifactRetriever{Contents: []byte("{}")},
			expectedFlightIds: nil,
		},
		{
			name:              "error response",
			retriever:         &MockArtifactRetriever{Err: errors.New("test error")},
			expectedFlightIds: nil,
			expectedErrors:    []string{"couldn't get", "test error"},
		},
		{
			name:      "saved response",
			retriever: &MockArtifactRetriever{Contents: newTestTracksJsonContents(6)},
			assertions: func(requirer *require.Assertions, savedResponses *testResponseSaver) {
				responses := savedResponses.responses
				requirer.Equal(1, len(responses))
				requirer.Equal(filepath.Join("adir", MakeFlightIdsArtifactFilename("irrelevant")), responses[0].name)
				sum := crc32.ChecksumIEEE(responses[0].contents)
				requirer.Equal(uint32(0xff2cac5), sum)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &RetrieverSaverApiImpl{Retriever: tc.retriever}
			responseSaver := &testResponseSaver{}
			api.Saver = &FileAeroApi{
				ArtifactsDir: "adir",
				FileSaver:    persistence.FileSaver{Writer: responseSaver.Save},
			}
			flightIds, err := api.GetFlightIds("irrelevant", tc.cutoffTime)
			requirer := require.New(t)
			if len(tc.expectedErrors) > 0 {
				requirer.Error(err)
				requirer.Nil(tc.expectedFlightIds)
				requirer.Nil(flightIds)
				for _, eet := range tc.expectedErrors {
					requirer.Contains(err.Error(), eet)
				}
			} else {
				requirer.NoError(err)
				requirer.Equal(tc.expectedFlightIds, flightIds)
				if tc.assertions != nil {
					tc.assertions(requirer, responseSaver)
				}
			}
		})
	}

}

func TestGetTrackForFlightId(t *testing.T) {

	type testCaseDef struct {
		name                  string
		retriever             ArtifactRetriever
		expectedPositionCount int
		expectedErrors        []string
		assertions            func(*require.Assertions, *Track, *testResponseSaver)
	}

	testCases := []testCaseDef{
		{
			name:                  "single empty track from JSON",
			retriever:             &MockArtifactRetriever{Contents: []byte(`{"positions": [{}]}`)},
			expectedPositionCount: 1,
		},
		{
			name:                  "small realistic track from JSON",
			retriever:             &MockArtifactRetriever{Contents: []byte(testfixtures.NewMockTestAeroApiTrackResponse())},
			expectedPositionCount: 19,
		},
		{
			name:                  "multiple tracks from struct",
			retriever:             &MockArtifactRetriever{Contents: newTestTracksJsonContents(6)},
			expectedPositionCount: 6,
			assertions: func(requirer *require.Assertions, response *Track, _ *testResponseSaver) {
				// check an arbitrary member of the returned collection
				requirer.Equal(newTestPosition(4), response.Positions[4])
			},
		},
		{
			name:                  "save track response",
			retriever:             &MockArtifactRetriever{Contents: newTestTracksJsonContents(6)},
			expectedPositionCount: 6,
			assertions: func(requirer *require.Assertions, _ *Track, saver *testResponseSaver) {
				// verify the expected saved response
				responses := saver.responses
				requirer.Equal(1, len(responses))
				requirer.Equal(MakeTrackArtifactFilename("irrelevant"), responses[0].name)
				sum := crc32.ChecksumIEEE(responses[0].contents)
				requirer.Equal(uint32(0xff2cac5), sum)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &RetrieverSaverApiImpl{Retriever: tc.retriever}
			responseSaver := &testResponseSaver{}
			api.Saver = &FileAeroApi{
				FileSaver: persistence.FileSaver{Writer: responseSaver.Save},
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
					tc.assertions(requirer, track, responseSaver)
				}
			}
		})
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

func newTestFlightsJsonContents(n int) []byte {
	var fr FlightsResponse
	for _, id := range newTestFlightIds(n) {
		fr.Flights = append(fr.Flights, Flight{FlightId: id})
	}
	contents, err := json.Marshal(&fr)
	if err != nil {
		log.Fatal(err)
	}
	return contents
}

func newTestTracksJsonContents(n int) []byte {
	tr := Track{FlightId: fmt.Sprintf("track%d", n)}
	for i := 0; i < n; i++ {
		tr.Positions = append(tr.Positions, newTestPosition(i))
	}
	contents, err := json.Marshal(&tr)
	if err != nil {
		log.Fatal(err)
	}
	return contents
}

type savedResponse struct {
	name     string
	contents []byte
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

type testResponseSaver struct {
	responses []savedResponse
}

func (ts *testResponseSaver) Save(name string, contents []byte) error {
	ts.responses = append(ts.responses, savedResponse{
		name:     name,
		contents: contents,
	})
	return nil
}
