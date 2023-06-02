package aeroapi

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/noodnik2/flightvisualizer/pkg/persistence"
)

func TestFileAeroApi_Save(t *testing.T) {

	testCases := []struct{ name string }{
		{name: "one and only for now"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)

			const saveFilename = "saveFilename"
			const saveContents = "this is a string to be saved"

			const loadFilename = "loadFilename"
			const loadContents = "this is the string to be loaded"

			var actuallySavedContents []byte
			var actuallySavedFilename string
			var actuallyLoadedFilename string

			fileAeroApi := &FileAeroApi{
				FileLoader: persistence.FileLoader{Reader: func(filePath string) ([]byte, error) {
					actuallyLoadedFilename = filePath
					return []byte(loadContents), nil
				}},
				FileSaver: persistence.FileSaver{Writer: func(filePath string, contents []byte) error {
					actuallySavedFilename = filePath
					actuallySavedContents = contents
					return nil
				}},
			}

			saveErr := fileAeroApi.Save(saveFilename, []byte(saveContents))
			requirer.NoError(saveErr)

			actuallyLoadedContents, readErr := fileAeroApi.Load(loadFilename)
			requirer.NoError(readErr)

			requirer.Equal(loadFilename, actuallyLoadedFilename)
			requirer.Equal(loadContents, string(actuallyLoadedContents))

			requirer.Equal(saveFilename, actuallySavedFilename)
			requirer.Equal(saveContents, string(actuallySavedContents))

		})
	}

}

func TestFileAeroApi_GetUris(t *testing.T) {

	testCases := []struct {
		name                 string
		artifactsDir         string
		flightIdsFileName    string
		flightId             string
		tailNumber           string
		cutoffTime           time.Time
		expectedFlightIdsUri string
		expectedTrackUri     string
	}{
		{
			name:                 "without artifacts dir",
			flightId:             "aFid",
			tailNumber:           "aT#",
			expectedFlightIdsUri: makeFlightIdsArtifactFilename("aT#"),
			expectedTrackUri:     makeTrackArtifactFilename("aFid"),
		},
		{
			name:                 "without cutoff time",
			artifactsDir:         "bDir",
			flightId:             "bFid",
			tailNumber:           "bT#",
			expectedFlightIdsUri: filepath.Join("bDir", makeFlightIdsArtifactFilename("bT#")),
			expectedTrackUri:     filepath.Join("bDir", makeTrackArtifactFilename("bFid")),
		},
		{
			name:                 "with UTC cutoff time",
			artifactsDir:         "cDir",
			flightId:             "cFid",
			tailNumber:           "cT#",
			cutoffTime:           time.Date(2023, 5, 24, 14, 2, 3, 4, time.UTC),
			expectedFlightIdsUri: filepath.Join("cDir", makeFlightIdsArtifactFilename("cT#_cutoff-20230524T140203Z")),
			expectedTrackUri:     filepath.Join("cDir", makeTrackArtifactFilename("cFid")),
		},
		{
			name:                 "with non-UTC cutoff time",
			artifactsDir:         "dDir",
			flightId:             "dFid",
			tailNumber:           "dT#",
			cutoffTime:           time.Date(2023, 5, 24, 14, 2, 3, 4, time.FixedZone("PDT", -7*60*60)),
			expectedFlightIdsUri: filepath.Join("dDir", makeFlightIdsArtifactFilename("dT#_cutoff-20230524T140203-0700")),
			expectedTrackUri:     filepath.Join("dDir", makeTrackArtifactFilename("dFid")),
		},
		{
			name:                 "with FlightIdsFileName",
			artifactsDir:         "eDir",
			flightIdsFileName:    makeFlightIdsArtifactFilename("f99"),
			flightId:             "eFid",
			expectedFlightIdsUri: filepath.Join("eDir", makeFlightIdsArtifactFilename("f99")),
			expectedTrackUri:     filepath.Join("eDir", makeTrackArtifactFilename("eFid")),
		},
		{
			name:                 "with track filename",
			artifactsDir:         "fDir",
			flightIdsFileName:    makeTrackArtifactFilename("f99"),
			flightId:             "fFid",
			expectedFlightIdsUri: "[f99]",
			expectedTrackUri:     filepath.Join("fDir", makeTrackArtifactFilename("fFid")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)

			fileAeroApi := &FileAeroApi{ArtifactsDir: tc.artifactsDir, FlightIdsFileName: tc.flightIdsFileName}
			requirer.Equal(tc.expectedFlightIdsUri, fileAeroApi.GetFlightIdsRef(tc.tailNumber, tc.cutoffTime))
			requirer.Equal(tc.expectedTrackUri, fileAeroApi.GetTrackForFlightUri(tc.flightId))
		})
	}
}
