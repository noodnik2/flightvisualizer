package aeroapi

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/noodnik2/flightvisualizer/pkg/persistence"
)

type FileAeroApi struct {
	ArtifactsDir      string
	FlightIdsFileName string
	persistence.FileLoader
	persistence.FileSaver
}

const trackArtifactFilenamePrefix = "fvt_"
const trackArtifactFilenameSuffix = ".json"
const flightIdsArtifactFilenamePrefix = "fvf_"
const flightIdsArtifactFilenameSuffix = ".json"

func MakeTrackArtifactFilename(flightId string) string {
	return trackArtifactFilenamePrefix + flightId + trackArtifactFilenameSuffix
}

func IsTrackArtifactFilename(fn string) bool {
	base := filepath.Base(fn)
	return strings.HasPrefix(base, trackArtifactFilenamePrefix) && strings.HasSuffix(base, trackArtifactFilenameSuffix)
}

func MakeFlightIdsArtifactFilename(queryId string) string {
	return flightIdsArtifactFilenamePrefix + queryId + flightIdsArtifactFilenameSuffix
}

func IsFlightIdsArtifactFilename(fn string) bool {
	base := filepath.Base(fn)
	return strings.HasPrefix(base, flightIdsArtifactFilenamePrefix) && strings.HasSuffix(base, flightIdsArtifactFilenameSuffix)
}

func (c *FileAeroApi) GetFlightIdsRef(tailNumber string, cutoffTime time.Time) (string, error) {
	var fileName string
	if c.FlightIdsFileName != "" {
		fileName = c.FlightIdsFileName
	} else {
		var queryId string
		if cutoffTime.IsZero() {
			queryId = tailNumber
		} else {
			queryId = fmt.Sprintf("%s_cutoff-%s", tailNumber, cutoffTime.Format("20060102T150405Z0700"))
		}
		fileName = MakeFlightIdsArtifactFilename(queryId)
	}
	if !IsFlightIdsArtifactFilename(filepath.Base(fileName)) {
		return "", fmt.Errorf("unrecognized flight ids filename(%s)", fileName)
	}
	if filepath.Dir(fileName) == "." {
		// use the artifacts directory if not specified
		return filepath.Join(c.ArtifactsDir, fileName), nil
	}
	return fileName, nil
}

func (c *FileAeroApi) GetTrackForFlightRef(flightId string) string {
	artifactDir := filepath.Dir(c.FlightIdsFileName)
	if artifactDir == "." {
		// use the artifacts directory if not specified
		artifactDir = c.ArtifactsDir
	}
	return filepath.Join(artifactDir, MakeTrackArtifactFilename(flightId))
}
