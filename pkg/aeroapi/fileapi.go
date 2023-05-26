package aeroapi

import (
    "fmt"
    "log"
    "path/filepath"
    "time"

    "github.com/noodnik2/flightvisualizer/pkg/persistence"
)

type FileAeroApi struct {
    ArtifactsDir      string
    FlightIdsFileName string
    persistence.FileLoader
    persistence.FileSaver
}

func (c *FileAeroApi) GetFlightIdsUri(tailNumber string, cutoffTime time.Time) string {
    var fileName string
    if c.FlightIdsFileName != "" {
        fileName = c.FlightIdsFileName
    } else {
        if cutoffTime.IsZero() {
            fileName = fmt.Sprintf("fvf_%s.json", tailNumber)
        } else {
            fileName = fmt.Sprintf("fvf_%s_cutoff-%s.json", tailNumber, cutoffTime.Format("20060102T150405Z0700"))
        }
    }
    if filepath.Dir(fileName) != "." {
        // don't touch it if it already has a directory part
        return fileName
    }
    return filepath.Join(c.ArtifactsDir, fileName)
}

func (c *FileAeroApi) GetTrackForFlightUri(flightId string) string {
    return filepath.Join(c.ArtifactsDir, fmt.Sprintf("fvt_%s.json", flightId))
}

func (c *FileAeroApi) Load(fileName string) ([]byte, error) {
    log.Printf("INFO: reading from file(%s)\n", fileName)
    return c.FileLoader.Load(fileName)
}

func (c *FileAeroApi) Save(fileName string, contents []byte) error {
    log.Printf("INFO: saving to file(%s)\n", fileName)
    return c.FileSaver.Save(fileName, contents)
}
