package cmd

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	iaeroapi "github.com/noodnik2/kmlflight/internal/aeroapi"
	"github.com/noodnik2/kmlflight/internal/kml"
	ios "github.com/noodnik2/kmlflight/internal/os"
	"github.com/noodnik2/kmlflight/internal/persistence"
	"github.com/noodnik2/kmlflight/pkg/aeroapi"
)

const cmdFlagTracksMock = "mock"
const cmdFlagTracksSaveResponses = "saveResponses"
const cmdFlagTracksNoBanking = "noBanking"
const cmdFlagTracksLaunch = "launchFirstKml"
const cmdFlagTracksKind = "kind"
const cmdFlagTracksOutputDir = "outputDir"
const cmdFlagTracksCutoffTime = "cutoffTime"
const cmdFlagTracksFlightCount = "flightCount"

const cmdFlagTracksKindStd = "std"
const cmdFlagTracksKindPlacemark = "placemark"
const cmdFlagTracksKindCamera = "camera"

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP(cmdFlagTracksMock, "m", false, "Use mock backend")
	versionCmd.Flags().BoolP(cmdFlagTracksLaunch, "l", false, "Launch the KML file (from the most recent flight) once created")
	versionCmd.Flags().BoolP(cmdFlagTracksSaveResponses, "s", false, "Save responses from AeroAPI requests")
	versionCmd.Flags().BoolP(cmdFlagTracksNoBanking, "b", false, "Disable banking heuristic calculations")
	versionCmd.Flags().String(cmdFlagTracksKind, cmdFlagTracksKindCamera, "Kind of tour to create")
	versionCmd.Flags().String(cmdFlagTracksOutputDir, ".", "Directory to receive file(s) created")
	versionCmd.Flags().String(cmdFlagTracksCutoffTime, "", "Cut off time for flight(s) to consider")
	versionCmd.Flags().Int(cmdFlagTracksFlightCount, 0, "Count of (most recent) flights to consider (0=unlimited)")
}

type tracksCommandArgs struct {
	tailNumber       string
	cutoffTime       *time.Time
	kmlKind          string
	outputDir        string
	flightCount      int
	verboseOperation bool
	saveResponses    bool
	launchFirstKml   bool
	useMockData      bool
	noBanking        bool
}

var versionCmd = &cobra.Command{
	Use:   "tracks",
	Short: "Creates KML track log(s) from a flight by querying FLightAware AeroAPI",
	Long:  `Uses the parameters supplied to invoke AeroAPI and build the KML track log from its response.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdArgs, parseErr := parseArgs(cmd, args)
		if parseErr != nil {
			return fmt.Errorf("invalid syntax: %w", parseErr)
		}

		var aeroApiGetter aeroapi.GetRequester
		var aeroApiErr error
		if cmdArgs.useMockData {
			aeroApiGetter, aeroApiErr = getMockGetter(cmdArgs.verboseOperation)
			log.Printf("INFO: using mock data")
		} else {
			aeroApiGetter, aeroApiErr = getAeroApiHttpGetter(cmdArgs.verboseOperation)
		}
		if aeroApiErr != nil {
			log.Fatalf("unable to access AeroAPI: %v", aeroApiErr)
			//notreached
		}

		aeroApi := &aeroapi.Api{Getter: aeroApiGetter}

		commandStartTime := time.Now()
		if cmdArgs.saveResponses {
			saver := &persistence.FileSaver{
				FilenameSuffix: ".json",
				Timestamp:      &commandStartTime,
			}
			aeroApi.Saver = saver.SaveEndpointResponse
		}

		var tracker kml.KmlTracker
		switch cmdArgs.kmlKind {
		case cmdFlagTracksKindStd:
			tracker = &kml.GxTracker{
				Name: "Standard",
				Builders: []kml.GxKmlBuilder{
					&kml.PathBuilder{
						Extrude: true,
						Color:   color.RGBA{R: 217, G: 51, B: 255},
					},
					&kml.CameraTrackBuilder{
						AddBankAngle: !cmdArgs.noBanking,
						VerboseFlag:  cmdArgs.verboseOperation,
					},
					&kml.VectorBuilder{},
				},
			}
		case cmdFlagTracksKindPlacemark:
			tracker = &kml.GxTracker{
				Name:     "Map Tour",
				Builders: []kml.GxKmlBuilder{&kml.PlacemarkBuilder{}},
			}
		case cmdFlagTracksKindCamera:
			tracker = &kml.GxTracker{
				Name: "Birds-eye",
				Builders: []kml.GxKmlBuilder{
					&kml.CameraTrackBuilder{
						AddBankAngle: !cmdArgs.noBanking,
						VerboseFlag:  cmdArgs.verboseOperation,
					},
				},
			}
		default:
			return fmt.Errorf("unrecognized kmlKind(%s); supported: %v", cmdArgs.kmlKind,
				[]string{cmdFlagTracksKindCamera, cmdFlagTracksKindPlacemark})
		}

		tc := iaeroapi.TracksConverter{
			TailNumber:  cmdArgs.tailNumber,
			CutoffTime:  cmdArgs.cutoffTime,
			FlightCount: cmdArgs.flightCount,
		}

		aeroKmls, createKmlErr := tc.Convert(aeroApi, tracker)
		if createKmlErr != nil {
			log.Fatalf("couldn't create KMLs: %v", createKmlErr)
			//notreached
		}

		nKmlDocs := len(aeroKmls)
		if cmdArgs.verboseOperation || nKmlDocs > 1 {
			log.Printf("INFO: writing %d %s KML document(s)", nKmlDocs, cmdArgs.kmlKind)
		}

		var firstKmlFilename string
		for _, aeroKml := range aeroKmls {
			flightTime := aeroKml.StartTime.Format(persistence.FnPrefixTimestampFormat)
			flightEndUpdate := getTsUpdate(aeroKml.EndTime.Format(persistence.FnPrefixTimestampFormat), flightTime)

			lft := len(flightTime)
			saveFn := fmt.Sprintf("%s-%s-%s-%s", cmdArgs.tailNumber, flightTime[2:lft-1], flightEndUpdate[2:], cmdArgs.kmlKind)
			kmlFilename, writeErr := (&persistence.KmzSaver{
				FileSaver: persistence.FileSaver{
					Timestamp:      &commandStartTime,
					FilenameSuffix: ".kmz",
				},
				Assets: aeroKml.KmlAssets,
			}).SaveNewKmz(saveFn, aeroKml.KmlDoc)
			if writeErr != nil {
				log.Fatalf("couldn't write to stdout: %v", writeErr)
				//notreached
			}

			if firstKmlFilename == "" {
				firstKmlFilename = kmlFilename
			}
		}

		if cmdArgs.launchFirstKml && firstKmlFilename != "" {
			if cmdArgs.verboseOperation {
				log.Printf("INFO: opening(%s) in %s\n", firstKmlFilename, runtime.GOOS)
			}
			if openErr := ios.LaunchFile(firstKmlFilename); openErr != nil {
				log.Fatalf("error returned from launching(%s): %v", firstKmlFilename, openErr)
				//notreached
			}
		}

		return nil
	},
}

func getAeroApiHttpGetter(isVerbose bool) (aeroapi.GetRequester, error) {
	httpAeroApi, dotFileErr := aeroapi.HttpApiFromDotFiles(".env.local")
	if dotFileErr != nil {
		return nil, dotFileErr
	}
	if isVerbose {
		httpAeroApi.Verbose = isVerbose
	}
	return httpAeroApi.Get, nil
}

func getMockGetter(isVerbose bool) (aeroapi.GetRequester, error) {
	mockGetter := func(endpoint string) ([]byte, error) {
		mockFilename := "testfixtures/pattern_practice_flight_id.json"
		if strings.Contains(endpoint, "/track") {
			mockFilename = "testfixtures/pattern_practice_track.json"
		}
		if isVerbose {
			log.Printf("INFO: satisfying request for(%s) with mock data from(%s)\n", endpoint, mockFilename)
		}
		return os.ReadFile(mockFilename)
	}
	return mockGetter, nil
}

func parseArgs(cmd *cobra.Command, args []string) (cmdArgs tracksCommandArgs, err error) {
	if len(args) != 1 {
		err = errors.New("please supply single argument: tail number")
		return
	}
	cmdArgs.tailNumber = args[0]

	var cutoffTimeString string
	if cutoffTimeString, err = cmd.Flags().GetString(cmdFlagTracksCutoffTime); err != nil {
		return
	}
	if cutoffTimeString != "" {
		var toTime time.Time
		if toTime, err = time.Parse(time.RFC3339, cutoffTimeString); err != nil {
			return
		}
		cmdArgs.cutoffTime = &toTime
	}

	if cmdArgs.outputDir, err = cmd.Flags().GetString(cmdFlagTracksOutputDir); err != nil {
		return
	}
	if cmdArgs.kmlKind, err = cmd.Flags().GetString(cmdFlagTracksKind); err != nil {
		return
	}
	if cmdArgs.flightCount, err = cmd.Flags().GetInt(cmdFlagTracksFlightCount); err != nil {
		return
	}
	if cmdArgs.verboseOperation, err = cmd.Flags().GetBool(cmdFlagRootVerbose); err != nil {
		return
	}
	if cmdArgs.saveResponses, err = cmd.Flags().GetBool(cmdFlagTracksSaveResponses); err != nil {
		return
	}
	if cmdArgs.launchFirstKml, err = cmd.Flags().GetBool(cmdFlagTracksLaunch); err != nil {
		return
	}
	if cmdArgs.useMockData, err = cmd.Flags().GetBool(cmdFlagTracksMock); err != nil {
		return
	}
	if cmdArgs.noBanking, err = cmd.Flags().GetBool(cmdFlagTracksNoBanking); err != nil {
		return
	}

	return
}

// getTsUpdate returns the suffix of tsUpdate which is not in common with tsBase
// example: { 23010203040506Z, 23010203050506Z } => 50506Z ('5' differs with '4' in tsBase)
func getTsUpdate(tsUpdate, tsBase string) string {
	i := 0
	for i < len(tsBase) && i < len(tsUpdate) && tsBase[i] == tsUpdate[i] {
		i++
	}
	return tsUpdate[i:]
}
