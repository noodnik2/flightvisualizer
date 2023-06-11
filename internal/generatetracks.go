package internal

import (
	"errors"
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
	persistence2 "github.com/noodnik2/flightvisualizer/pkg/persistence"
)

const (
	TracksLayerCamera          = "camera"
	TracksLayerPath            = "path"
	TracksLayerPlacemark       = "placemark"
	TracksLayerVector          = "vector"
	kmlArtifactsFilenamePrefix = "fvk_"
)

type sourceType int

const (
	sourceTypeUnrecognized        sourceType = iota
	sourceTypeMultiTrackRemote               // pull a remote "flight ids" document (e.g., from AeroAPI server)
	sourceTypeSingleTrackArtifact            // use a recorded "track" artifact as the source document
	sourceTypeMultiTrackArtifact             // use a recorded "flight ids" artifact as the source document
)

var TracksLayersSupported = []string{TracksLayerCamera, TracksLayerPath, TracksLayerPlacemark, TracksLayerVector}

type TracksCommandArgs struct {
	Config           Config
	LaunchFirstKml   bool
	NoBanking        bool
	SaveResponses    bool
	VerboseOperation bool
	DebugOperation   bool
	FromArtifacts    string
	ArtifactsDir     string
	KmlLayers        string
	TailNumber       string
	FlightCount      int
	CutoffTime       time.Time
}

func (tca TracksCommandArgs) GenerateTracks() error {

	kmlGenerator, getKmlGeneratorErr := tca.newKmlTrackGenerator(strings.Split(tca.KmlLayers, ","))
	if getKmlGeneratorErr != nil {
		return getKmlGeneratorErr
	}

	trackFactory, trackFactoryErr := tca.newTrackFactory()
	if trackFactoryErr != nil {
		return fmt.Errorf("no KML track factory could be created: %v", trackFactoryErr)
	}
	kmlTracks, generateKmlTracksErr := trackFactory(kmlGenerator)
	if generateKmlTracksErr != nil {
		return generateKmlTracksErr
	}

	firstKmlFilename, saveKmlErr := tca.saveKmlTracks(kmlTracks, kmlGenerator.Name)
	if saveKmlErr != nil {
		return saveKmlErr
	}

	// if indicated, "launch" the (first of the) generated KML visualization(s)
	if tca.LaunchFirstKml && firstKmlFilename != "" {
		log.Printf("INFO: Launching '%s'\n", firstKmlFilename)
		if openErr := ios.LaunchFile(firstKmlFilename); openErr != nil {
			return fmt.Errorf("error returned from launching(%s): %v", firstKmlFilename, openErr)
		}
	}

	return nil
}

func (tca TracksCommandArgs) saveKmlTracks(kmlTracks []*kml.Track, kmlLayersUi string) (string, error) {
	nKmlDocs := len(kmlTracks)
	if tca.IsVerbose() || nKmlDocs > 1 {
		log.Printf("INFO: writing %d %s KML document(s)\n", nKmlDocs, kmlLayersUi)
	}

	// save the KML document(s) produced along with their asset(s) as `.kmz` file(s)
	var firstKmlFilename string
	for _, aeroKml := range kmlTracks {
		kmzSaver := &persistence.KmzSaver{
			Saver:  &persistence2.FileSaver{},
			Assets: aeroKml.KmlAssets,
		}
		flightTimeRange := getTsFromTo(*aeroKml.StartTime, *aeroKml.EndTime)
		kmlFilename := filepath.Join(
			tca.getArtifactsDir(),
			fmt.Sprintf("%s%s_%s_%s.kmz", kmlArtifactsFilenamePrefix, tca.TailNumber, flightTimeRange, kmlLayersUi),
		)

		if writeErr := kmzSaver.Save(kmlFilename, aeroKml.KmlDoc); writeErr != nil {
			return "", fmt.Errorf("couldn't write output artifact(%s): %v", kmlFilename, writeErr)
		}

		if firstKmlFilename == "" {
			firstKmlFilename = kmlFilename
		}
	}
	return firstKmlFilename, nil
}

func (tca TracksCommandArgs) newKmlTrackGenerator(kmlLayers []string) (*kml.TrackBuilderEnsemble, error) {

	// order layer builder(s) for deterministic output
	sort.Strings(kmlLayers)

	builtLayers := make([]string, 0, len(kmlLayers))
	var kmlBuilders []builders.KmlTrackBuilder
	for _, kmlLayer := range kmlLayers {
		if len(builtLayers) > 0 && kmlLayer == builtLayers[len(builtLayers)-1] {
			// ignore duplicates
			continue
		}
		var kmlBuilder builders.KmlTrackBuilder
		switch kmlLayer {
		case TracksLayerCamera:
			kmlBuilder = &builders.CameraBuilder{
				AddBankAngle: !tca.NoBanking,
				DebugFlag:    tca.DebugOperation,
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
			return nil, fmt.Errorf("unrecognized kmlLayer(%s); supported: %v", kmlLayer,
				strings.Join(TracksLayersSupported, ","))
		}
		builtLayers = append(builtLayers, kmlLayer)
		kmlBuilders = append(kmlBuilders, kmlBuilder)
	}

	ensemble := &kml.TrackBuilderEnsemble{
		Name:     strings.Join(builtLayers, "-"),
		Builders: kmlBuilders,
	}
	return ensemble, nil
}

type kmlTrackFactory func(kml.TrackGenerator) ([]*kml.Track, error)

func (tca TracksCommandArgs) newTrackFactory() (kmlTrackFactory, error) {

	st, stErr := tca.getSourceType()
	if stErr != nil {
		return nil, stErr
	}

	switch st {
	case sourceTypeSingleTrackArtifact:
		return singleTrackArtifactFactory(tca), nil

	case sourceTypeMultiTrackArtifact:
		return multiTrackArtifactFactory(tca), nil

	case sourceTypeMultiTrackRemote:
		return multiTrackRemoteFactory(tca), nil
	}

	return nil, errors.New("can't determine source type")
}

func multiTrackRemoteFactory(tca TracksCommandArgs) kmlTrackFactory {
	return func(tracker kml.TrackGenerator) ([]*kml.Track, error) {

		var artifactSaver aeroapi.ArtifactSaver
		if tca.SaveResponses {
			artifactSaver = &aeroapi.FileAeroApi{ArtifactsDir: tca.getArtifactsDir()}
		}

		verbose := tca.IsVerbose()
		aeroApi := &aeroapi.RetrieverSaverApiImpl{
			// reading AeroAPI data from live AeroAPI REST API calls
			Retriever: &aeroapi.HttpAeroApi{
				Verbose: verbose,
				ApiKey:  tca.Config.AeroApiKey,
				ApiUrl:  tca.Config.AeroApiUrl,
			},

			Saver: artifactSaver,
		}
		tc := iaeroapi.TracksConverter{
			Verbose:     verbose,
			FlightCount: tca.FlightCount,
			TailNumber:  tca.TailNumber,
			CutoffTime:  tca.CutoffTime,
		}
		return tc.Convert(aeroApi, tracker)
	}
}

func multiTrackArtifactFactory(tca TracksCommandArgs) kmlTrackFactory {
	return func(tracker kml.TrackGenerator) ([]*kml.Track, error) {
		if tca.SaveResponses {
			// there's no good reason to save data already coming from local files
			log.Printf("NOTE: inappropriate 'save responses' option ignored\n")
		}
		aeroApi := &aeroapi.RetrieverSaverApiImpl{
			// reading AeroAPI data from saved artifact files
			Retriever: &aeroapi.FileAeroApi{
				ArtifactsDir:      tca.getArtifactsDir(),
				FlightIdsFileName: tca.FromArtifacts,
			},
		}
		tc := iaeroapi.TracksConverter{
			Verbose:     tca.IsVerbose(),
			FlightCount: tca.FlightCount,
			TailNumber:  tca.TailNumber,
			CutoffTime:  tca.CutoffTime,
		}
		return tc.Convert(aeroApi, tracker)
	}
}

func singleTrackArtifactFactory(tca TracksCommandArgs) kmlTrackFactory {
	return func(tracker kml.TrackGenerator) ([]*kml.Track, error) {
		track, getTfaErr := tca.getTrackFromArtifact()
		if getTfaErr != nil {
			return nil, getTfaErr
		}
		kmlTrack, err := tracker.Generate(track)
		if err != nil {
			return nil, err
		}
		return []*kml.Track{kmlTrack}, nil
	}
}

func (tca TracksCommandArgs) getSourceType() (sourceType, error) {

	if tca.FromArtifacts == "" {
		return sourceTypeMultiTrackRemote, nil
	}

	if aeroapi.IsTrackArtifactFilename(tca.FromArtifacts) {
		return sourceTypeSingleTrackArtifact, nil
	}

	if aeroapi.IsFlightIdsArtifactFilename(tca.FromArtifacts) {
		return sourceTypeMultiTrackArtifact, nil
	}

	return sourceTypeUnrecognized, fmt.Errorf("unrecognized artifact(%s)", tca.FromArtifacts)
}

func (tca TracksCommandArgs) getTrackFromArtifact() (*aeroapi.Track, error) {
	contents, loadErr := (&persistence2.FileLoader{}).Load(tca.FromArtifacts)
	if loadErr != nil {
		return nil, loadErr
	}
	return aeroapi.TrackFromJson(contents)
}

func (tca TracksCommandArgs) getArtifactsDir() string {
	if tca.ArtifactsDir != "" {
		return tca.ArtifactsDir
	}
	return tca.Config.ArtifactsDir
}

func (tca TracksCommandArgs) IsVerbose() bool {
	return tca.VerboseOperation || tca.Config.Verbose
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
