package aeroapi

import (
	"fmt"
	"log"
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

func makeTrackArtifactFilename(flightId string) string {
	return trackArtifactFilenamePrefix + flightId + trackArtifactFilenameSuffix
}

func makeFlightIdsArtifactFilename(queryId string) string {
	return flightIdsArtifactFilenamePrefix + queryId + flightIdsArtifactFilenameSuffix
}

func (c *FileAeroApi) GetFlightIdsRef(tailNumber string, cutoffTime time.Time) string {
	var fileName string
	if c.FlightIdsFileName != "" {
		fileName = c.FlightIdsFileName
		// handle the case of invoking a saved track file directly; use extracted flight ID
		baseFilename := filepath.Base(fileName)
		if strings.HasPrefix(baseFilename, trackArtifactFilenamePrefix) && strings.HasSuffix(baseFilename, trackArtifactFilenameSuffix) {
			// the "flight ID" is simply the part in between this prefix and suffix
			return fmt.Sprintf("[%s]", baseFilename[4:len(baseFilename)-5])
		}
	} else {
		var queryId string
		if cutoffTime.IsZero() {
			queryId = tailNumber
		} else {
			queryId = fmt.Sprintf("%s_cutoff-%s", tailNumber, cutoffTime.Format("20060102T150405Z0700"))
		}
		fileName = makeFlightIdsArtifactFilename(queryId)
	}
	if filepath.Dir(fileName) == "." {
		// use the artifacts directory if not specified
		return filepath.Join(c.ArtifactsDir, fileName)
	}
	return fileName
}

func (c *FileAeroApi) GetTrackForFlightUri(flightId string) string {
	artifactDir := filepath.Dir(c.FlightIdsFileName)
	if artifactDir == "." {
		// use the artifacts directory if not specified
		artifactDir = c.ArtifactsDir
	}
	return filepath.Join(artifactDir, makeTrackArtifactFilename(flightId))
}

func (c *FileAeroApi) Load(fileName string) ([]byte, error) {
	log.Printf("INFO: reading from file(%s)\n", fileName)
	return c.FileLoader.Load(fileName)
}

func (c *FileAeroApi) Save(fileName string, contents []byte) error {
	log.Printf("INFO: saving to file(%s)\n", fileName)
	return c.FileSaver.Save(fileName, contents)
}
