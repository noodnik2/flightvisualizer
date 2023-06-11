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
		expectedErrors       []string
	}{
		{
			name:                 "without artifacts dir",
			flightId:             "aFid",
			tailNumber:           "aT#",
			expectedFlightIdsUri: MakeFlightIdsArtifactFilename("aT#"),
			expectedTrackUri:     MakeTrackArtifactFilename("aFid"),
		},
		{
			name:                 "without cutoff time",
			artifactsDir:         "bDir",
			flightId:             "bFid",
			tailNumber:           "bT#",
			expectedFlightIdsUri: filepath.Join("bDir", MakeFlightIdsArtifactFilename("bT#")),
			expectedTrackUri:     filepath.Join("bDir", MakeTrackArtifactFilename("bFid")),
		},
		{
			name:                 "with UTC cutoff time",
			artifactsDir:         "cDir",
			flightId:             "cFid",
			tailNumber:           "cT#",
			cutoffTime:           time.Date(2023, 5, 24, 14, 2, 3, 4, time.UTC),
			expectedFlightIdsUri: filepath.Join("cDir", MakeFlightIdsArtifactFilename("cT#_cutoff-20230524T140203Z")),
			expectedTrackUri:     filepath.Join("cDir", MakeTrackArtifactFilename("cFid")),
		},
		{
			name:                 "with non-UTC cutoff time",
			artifactsDir:         "dDir",
			flightId:             "dFid",
			tailNumber:           "dT#",
			cutoffTime:           time.Date(2023, 5, 24, 14, 2, 3, 4, time.FixedZone("PDT", -7*60*60)),
			expectedFlightIdsUri: filepath.Join("dDir", MakeFlightIdsArtifactFilename("dT#_cutoff-20230524T140203-0700")),
			expectedTrackUri:     filepath.Join("dDir", MakeTrackArtifactFilename("dFid")),
		},
		{
			name:                 "with FlightIdsFileName",
			artifactsDir:         "eDir",
			flightIdsFileName:    MakeFlightIdsArtifactFilename("f99"),
			flightId:             "eFid",
			expectedFlightIdsUri: filepath.Join("eDir", MakeFlightIdsArtifactFilename("f99")),
			expectedTrackUri:     filepath.Join("eDir", MakeTrackArtifactFilename("eFid")),
		},
		{
			name:                 "with unrecognized (e.g., track) filename",
			artifactsDir:         "fDir",
			flightIdsFileName:    MakeTrackArtifactFilename("f99"),
			flightId:             "fFid",
			expectedFlightIdsUri: "[f99]",
			expectedTrackUri:     filepath.Join("fDir", MakeTrackArtifactFilename("fFid")),
			expectedErrors:       []string{"unrecognized"}, // this function no longer handles track files
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)

			fileAeroApi := &FileAeroApi{ArtifactsDir: tc.artifactsDir, FlightIdsFileName: tc.flightIdsFileName}
			fidsRef, getFidsRefErr := fileAeroApi.GetFlightIdsRef(tc.tailNumber, tc.cutoffTime)
			if tc.expectedErrors != nil {
				requirer.Error(getFidsRefErr)
				for _, expectedErr := range tc.expectedErrors {
					requirer.Contains(getFidsRefErr.Error(), expectedErr)
				}
			} else {
				requirer.NoError(getFidsRefErr)
				requirer.Equal(tc.expectedFlightIdsUri, fidsRef)
				requirer.Equal(tc.expectedTrackUri, fileAeroApi.GetTrackForFlightRef(tc.flightId))
			}
		})
	}
}
