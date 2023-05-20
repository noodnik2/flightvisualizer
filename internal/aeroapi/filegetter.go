package aeroapi

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/noodnik2/flightvisualizer/internal/persistence"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type AeroApiFileGetterFactory struct {
	Verbose     bool
	AssetReader func(name string) ([]byte, error)
	AssetFolder string
}

func (mg *AeroApiFileGetterFactory) NewRequester(flightIdsFilename string) (aeroapi.GetRequester, error) {

	assetReader := mg.AssetReader
	if assetReader == nil {
		assetReader = os.ReadFile
	}

	flightIdsJsonBytes, flightIdReaderErr := assetReader(joinPathIfNeeded(mg.AssetFolder, flightIdsFilename))
	if flightIdReaderErr != nil {
		return nil, flightIdReaderErr
	}

	flights, flightsErr := aeroapi.FlightsFromJson(flightIdsJsonBytes)
	if flightsErr != nil {
		return nil, flightsErr
	}

	var trackFilenames []string
	filenameSetPrefix := GetFilenameSetPrefix(flightIdsFilename)
	for _, flight := range flights.Flights {
		trackFilename := fmt.Sprintf("%s%s_track.json", filenameSetPrefix, flight.FlightId)
		trackFilenames = append(trackFilenames, trackFilename)
	}

	var trackIndex int

	localGetter := func(endpoint string) ([]byte, error) {
		assetFilename := flightIdsFilename
		if strings.Contains(endpoint, "/track") {
			if trackIndex >= len(trackFilenames) {
				return nil, fmt.Errorf("no more mock tracks available")
			}
			assetFilename = trackFilenames[trackIndex]
			trackIndex++
		}
		if mg.Verbose {
			log.Printf("INFO: satisfying request for(%s) with local data from(%s)\n", endpoint, assetFilename)
		}
		return assetReader(joinPathIfNeeded(mg.AssetFolder, assetFilename))
	}
	return localGetter, nil
}

var invocationPrefixRegexp *regexp.Regexp

func init() {
	prefixDigitCount := len(persistence.FnPrefixTimestampFormat[2 : len(persistence.FnPrefixTimestampFormat)-1])
	// prefixDigitCount is the length of the timestamp prefix used for the output files; we expect the first
	// two first characters of the full timestamp are stripped, and the final 'Z' is accounted for with len-1
	invocationPrefixRegexp = regexp.MustCompile(fmt.Sprintf(`([0-9]{%d}Z-flights_).*`, prefixDigitCount))
}

func GetFilenameSetPrefix(filenameFromSet string) string {
	prefix := invocationPrefixRegexp.FindStringSubmatch(filepath.Base(filenameFromSet))
	if len(prefix) > 1 {
		return prefix[1]
	}
	return ""
}

func joinPathIfNeeded(leftPart, rightPart string) string {
	if filepath.IsAbs(rightPart) {
		return rightPart
	}
	return filepath.Join(leftPart, rightPart)
}
