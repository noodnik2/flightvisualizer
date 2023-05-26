package internal

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	iaeroapi "github.com/noodnik2/flightvisualizer/internal/aeroapi"
	"github.com/noodnik2/flightvisualizer/internal/kml"
	"github.com/noodnik2/flightvisualizer/internal/kml/builders"
	ios "github.com/noodnik2/flightvisualizer/internal/os"
	"github.com/noodnik2/flightvisualizer/internal/persistence"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

const TracksLayerCamera = "camera"
const TracksLayerPath = "path"
const TracksLayerPlacemark = "placemark"
const TracksLayerVector = "vector"

var TracksLayersSupported = []string{TracksLayerCamera, TracksLayerPath, TracksLayerPlacemark, TracksLayerVector}

type TracksCommandArgs struct {
	LaunchFirstKml   bool
	NoBanking        bool
	SaveResponses    bool
	VerboseOperation bool
	FromArtifacts    string
	ArtifactsDir     string
	KmlLayers        string
	TailNumber       string
	FlightCount      int
	CutoffTime       time.Time
}

func GenerateTracks(cmdArgs TracksCommandArgs) error {

	// fileBasedAeroApi is a file-based instance of our AeroApi library API,
	// used to save or retrieve saved AeroAPI responses
	fileBasedAeroApi := &aeroapi.FileAeroApi{ArtifactsDir: cmdArgs.ArtifactsDir}

	aeroApi := &aeroapi.RetrieverSaverApiImpl{}
	var cutoffTime time.Time
	var tailNumber string
	var flightCount int

	if cmdArgs.FromArtifacts != "" {
		// reading AeroAPI data from saved artifact files
		aeroApi.Retriever = fileBasedAeroApi
		fileBasedAeroApi.FlightIdsFileName = cmdArgs.FromArtifacts
	} else {
		// reading AeroAPI data from live AeroAPI REST API calls
		var aeroApiErr error
		if aeroApi.Retriever, aeroApiErr = getAeroApiHttpRetriever(cmdArgs.VerboseOperation); aeroApiErr != nil {
			log.Fatalf("ERROR: unable to access AeroAPI: %v", aeroApiErr)
			//notreached
		}
		if cmdArgs.SaveResponses {
			aeroApi.Saver = fileBasedAeroApi
		}
		// these user-supplied values are only relevant when we call the external AeroAPI
		cutoffTime = cmdArgs.CutoffTime
		tailNumber = cmdArgs.TailNumber
		flightCount = cmdArgs.FlightCount
	}

	// construct builder(s) for selected "layer(s)"
	sortedKmlLayers := strings.Split(cmdArgs.KmlLayers, ",")
	sort.Slice(sortedKmlLayers, func(i, j int) bool {
		// order builders for deterministic ordering of KML layers
		return sortedKmlLayers[i] < sortedKmlLayers[j]
	})

	var kmlBuilders []builders.GxKmlBuilder
	for _, kmlLayer := range sortedKmlLayers {
		var kmlBuilder builders.GxKmlBuilder
		switch kmlLayer {
		case TracksLayerCamera:
			kmlBuilder = &builders.CameraBuilder{
				AddBankAngle: !cmdArgs.NoBanking,
				VerboseFlag:  cmdArgs.VerboseOperation,
			}
		case TracksLayerPath:
			kmlBuilder = &builders.PathBuilder{
				Extrude: true,
				Color:   color.RGBA{R: 217, G: 51, B: 255},
			}
		case TracksLayerPlacemark:
			kmlBuilder = &builders.PlacemarkBuilder{}
		case TracksLayerVector:
			kmlBuilder = &builders.VectorBuilder{}
		default:
			return fmt.Errorf("unrecognized kmlLayer(%s); supported: %v", kmlLayer,
				strings.Join(TracksLayersSupported, ","))
		}
		kmlBuilders = append(kmlBuilders, kmlBuilder)
	}

	// construct & invoke track converter using builder(s)
	tc := iaeroapi.TracksConverter{
		Verbose:     cmdArgs.VerboseOperation,
		TailNumber:  tailNumber,
		CutoffTime:  cutoffTime,
		FlightCount: flightCount,
		Api:         aeroApi,
	}

	tracker := &kml.GxTracker{Builders: kmlBuilders}
	aeroKmls, createKmlErr := tc.Convert(tracker)
	if createKmlErr != nil {
		log.Fatalf("ERROR: couldn't create KMLs: %v", createKmlErr)
		//notreached
	}

	sortedKmlLayersStr := strings.Join(sortedKmlLayers, "-")

	nKmlDocs := len(aeroKmls)
	if cmdArgs.VerboseOperation || nKmlDocs > 1 {
		log.Printf("INFO: writing %d %s KML document(s)", nKmlDocs, sortedKmlLayersStr)
	}

	// save the KML document(s) produced along with their asset(s) as `.kmz` file(s)
	var firstKmlFilename string
	for _, aeroKml := range aeroKmls {
		kmzSaver := &persistence.KmzSaver{
			Saver:  fileBasedAeroApi,
			Assets: aeroKml.KmlAssets,
		}
		flightTimeRange := getTsFromTo(*aeroKml.StartTime, *aeroKml.EndTime)
		kmlFilename := filepath.Join(
			fileBasedAeroApi.ArtifactsDir,
			fmt.Sprintf("fvk_%s_%s_%s.kmz", cmdArgs.TailNumber, flightTimeRange, sortedKmlLayersStr),
		)

		if writeErr := kmzSaver.Save(kmlFilename, aeroKml.KmlDoc); writeErr != nil {
			log.Fatalf("ERROR: couldn't write output artifact(%s): %v", kmlFilename, writeErr)
			//notreached
		}

		if firstKmlFilename == "" {
			firstKmlFilename = kmlFilename
		}
	}

	// if indicated, "launch" the (first of the) generated KML visualization(s)
	if cmdArgs.LaunchFirstKml && firstKmlFilename != "" {
		log.Printf("INFO: Launching '%s'", firstKmlFilename)
		if openErr := ios.LaunchFile(firstKmlFilename); openErr != nil {
			log.Fatalf("ERROR: error returned from launching(%s): %v", firstKmlFilename, openErr)
			//notreached
		}
	}

	return nil
}

func getAeroApiHttpRetriever(isVerbose bool) (*aeroapi.HttpAeroApi, error) {
	httpAeroApi, dotFileErr := aeroapi.HttpApiFromDotFiles(".env.local")
	if dotFileErr != nil {
		return nil, dotFileErr
	}
	if isVerbose {
		httpAeroApi.Verbose = isVerbose
	}
	return httpAeroApi, nil
}

const fnPrefixTimestampFormat = "20060102150405Z"

// GetTsFromTo returns a string representation of a time range using fnPrefixTimestampFormat
// to format the "from" time, and a subsequence of that for the "to" time, with leading common
// prefix removed.  Example:
//
// { 2023010203040506Z, 2023010203050506Z } => "23010203040506Z-50506Z" ('5' differs with '4' in tsBase)
func getTsFromTo(from, to time.Time) string {
	fromFmt := from.Format(fnPrefixTimestampFormat)[2:]
	toFmt := to.Format(fnPrefixTimestampFormat)[2:]

	i := 0
	for i < len(fromFmt) && i < len(toFmt) && fromFmt[i] == toFmt[i] {
		i++
	}
	return fmt.Sprintf("%s-%s", fromFmt, toFmt[i:])
}
