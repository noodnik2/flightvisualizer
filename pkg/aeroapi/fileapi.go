package aeroapi

import (
    "fmt"
    "log"
    "path/filepath"
    "time"

    "github.com/noodnik2/flightvisualizer/pkg/persistence"
)

type FileAeroApi struct {
    Verbose      bool
    ArtifactsDir string
    persistence.FileLoader
    persistence.FileSaver
}

func (c *FileAeroApi) GetFlightIdsUri(tailNumber string, cutoffTime *time.Time) string {
    var base string
    if cutoffTime == nil {
        base = fmt.Sprintf("fvf_%s.json", tailNumber)
    } else {
        base = fmt.Sprintf("fvf_%s_cutoff_%s.json", tailNumber, cutoffTime.Format("20060102T150405Z0700"))
    }
    return filepath.Join(c.ArtifactsDir, base)
}

func (c *FileAeroApi) GetTrackForFlightUri(flightId string) string {
    return filepath.Join(c.ArtifactsDir, fmt.Sprintf("fvt_%s.json", flightId))
}

func (c *FileAeroApi) Load(fileName string) ([]byte, error) {
    if c.Verbose {
        log.Printf("INFO: requesting from file(%s)\n", fileName)
    }
    return c.FileLoader.Load(fileName)
}

func (c *FileAeroApi) Save(fileName string, contents []byte) error {
    if c.Verbose {
        log.Printf("INFO: saving to file(%s)\n", fileName)
    }
    return c.FileSaver.Save(fileName, contents)
}

//func (c *FileAeroApi) resolvePathIfNeeded(fileName string) string {
//    if filepath.Dir(fileName) != "" {
//        return fileName
//    }
//    return filepath.Join(c.ArtifactsDir, fileName)
//}
