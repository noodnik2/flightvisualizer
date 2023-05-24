package aeroapi

import (
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
        cutoffTime           *time.Time
        expectedFlightIdsUri string
        expectedTrackUri     string
    }{
        {
            name:                 "without artifacts dir",
            flightId:             "aFid",
            tailNumber:           "aT#",
            expectedFlightIdsUri: "fvf_aT#.json",
            expectedTrackUri:     "fvt_aFid.json",
        },
        {
            name:                 "without cutoff time",
            artifactsDir:         "bDir",
            flightId:             "bFid",
            tailNumber:           "bT#",
            expectedFlightIdsUri: "bDir/fvf_bT#.json",
            expectedTrackUri:     "bDir/fvt_bFid.json",
        },
        {
            name:                 "with cutoff time",
            artifactsDir:         "cDir",
            flightId:             "cFid",
            tailNumber:           "cT#",
            cutoffTime:           &time.Time{},
            expectedFlightIdsUri: "cDir/fvf_cT#_cutoff_00010101T000000Z.json",
            expectedTrackUri:     "cDir/fvt_cFid.json",
        },
        {
            name:                 "with FlightIdsFileName",
            artifactsDir:         "dDir",
            flightIdsFileName:    "f99.json",
            flightId:             "dFid",
            expectedFlightIdsUri: "dDir/f99.json",
            expectedTrackUri:     "dDir/fvt_dFid.json",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            requirer := require.New(t)

            fileAeroApi := &FileAeroApi{ArtifactsDir: tc.artifactsDir, FlightIdsFileName: tc.flightIdsFileName}
            requirer.Equal(tc.expectedFlightIdsUri, fileAeroApi.GetFlightIdsUri(tc.tailNumber, tc.cutoffTime))
            requirer.Equal(tc.expectedTrackUri, fileAeroApi.GetTrackForFlightUri(tc.flightId))
        })
    }
}
