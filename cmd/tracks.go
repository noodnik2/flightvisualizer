package cmd

import (
    "fmt"
    "image/color"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "sort"
    "strings"
    "time"

    "github.com/spf13/cobra"

    iaeroapi "github.com/noodnik2/flightvisualizer/internal/aeroapi"
    "github.com/noodnik2/flightvisualizer/internal/kml"
    "github.com/noodnik2/flightvisualizer/internal/kml/builders"
    ios "github.com/noodnik2/flightvisualizer/internal/os"
    "github.com/noodnik2/flightvisualizer/internal/persistence"
    "github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

const cmdFlagTracksTailNumber = "tailNumber"
const cmdFlagTracksFromArtifacts = "fromArtifacts"
const cmdFlagTracksSaveArtifacts = "saveArtifacts"
const cmdFlagTracksNoBanking = "noBanking"
const cmdFlagTracksLaunch = "launch"
const cmdFlagTracksLayers = "layers"
const cmdFlagTracksArtifactsDir = "artifactsDir"
const cmdFlagTracksCutoffTime = "cutoffTime"
const cmdFlagTracksFlightCount = "flightCount"

const cmdFlagTracksLayerCamera = "camera"
const cmdFlagTracksLayerPath = "path"
const cmdFlagTracksLayerPlacemark = "placemark"
const cmdFlagTracksLayerVector = "vector"

var cmdFlagTracksLayersDefault = []string{cmdFlagTracksLayerCamera, cmdFlagTracksLayerPath, cmdFlagTracksLayerVector}
var cmdFlagTracksLayersSupported = []string{cmdFlagTracksLayerCamera, cmdFlagTracksLayerPath, cmdFlagTracksLayerPlacemark, cmdFlagTracksLayerVector}

func init() {
    rootCmd.AddCommand(versionCmd)
    versionCmd.Flags().BoolP(cmdFlagTracksLaunch, "l", false, "Launch the KML file (depicting the track of the most recent flight retrieved) once created")
    versionCmd.Flags().BoolP(cmdFlagTracksNoBanking, "b", false, "Disable banking heuristic calculations")
    versionCmd.Flags().BoolP(cmdFlagTracksSaveArtifacts, "s", false, "Save responses from AeroAPI requests")
    versionCmd.Flags().StringP(cmdFlagTracksFromArtifacts, "r", "", "Use saved responses instead of querying AeroAPI")
    versionCmd.Flags().String(cmdFlagTracksTailNumber, "", "Tail number identifier")
    versionCmd.Flags().String(cmdFlagTracksLayers, strings.Join(cmdFlagTracksLayersDefault, ","), "Layer(s) of the KML depiction to create")
    versionCmd.Flags().String(cmdFlagTracksArtifactsDir, "artifacts", "Directory to save or load artifacts")
    versionCmd.Flags().String(cmdFlagTracksCutoffTime, "", "Cut off time for flight(s) to consider")
    versionCmd.Flags().Int(cmdFlagTracksFlightCount, 0, "Count of (most recent) flights to consider (0=unlimited)")
    _ = versionCmd.MarkFlagRequired(cmdFlagTracksTailNumber)
}

type tracksCommandArgs struct {
    launchFirstKml     bool
    noBanking          bool
    saveResponses      bool
    verboseOperation   bool
    artifactsDir       string
    fromSavedResponses string
    kmlLayers          string
    tailNumber         string
    flightCount        int
    cutoffTime         *time.Time
}

var versionCmd = &cobra.Command{
    Use:   "tracks",
    Short: "Tracks",
    Long:  `Generates KML visualizations of flight track logs retrieved from FlightAware's AeroAPI`,
    RunE: func(cmd *cobra.Command, args []string) error {

        cmdArgs, parseErr := parseArgs(cmd)
        if parseErr != nil {
            return fmt.Errorf("invalid syntax: %w", parseErr)
        }

        // fileBasedAeroApi is a file-based instance of our AeroApi library API,
        // used to save or retrieve saved AeroAPI responses
        fileBasedAeroApi := &aeroapi.FileAeroApi{
            Verbose:      cmdArgs.verboseOperation,
            ArtifactsDir: cmdArgs.artifactsDir,
        }

        aeroApi := &aeroapi.RetrieverSaverApiImpl{}
        if cmdArgs.fromSavedResponses != "" {
            aeroApi.Retriever = fileBasedAeroApi
            if cmdArgs.saveResponses {
                // there's no reason to save data read from files
                log.Printf("NOTE: ignoring '%s' option; incompatible with '%s'",
                    cmdFlagTracksSaveArtifacts, cmdFlagTracksFromArtifacts)
            }
        } else {
            var aeroApiErr error
            if aeroApi.Retriever, aeroApiErr = getAeroApiHttpRetriever(cmdArgs.verboseOperation); aeroApiErr != nil {
                log.Fatalf("unable to access AeroAPI: %v", aeroApiErr)
                //notreached
            }
            if cmdArgs.saveResponses {
                aeroApi.Saver = fileBasedAeroApi
            }
        }

        sortedKmlLayers := strings.Split(cmdArgs.kmlLayers, ",")
        sort.Slice(sortedKmlLayers, func(i, j int) bool {
            // order builders for deterministic ordering of KML layers
            return sortedKmlLayers[i] < sortedKmlLayers[j]
        })

        var kmlBuilders []kml.GxKmlBuilder
        for _, kmlLayer := range sortedKmlLayers {
            var kmlBuilder kml.GxKmlBuilder
            switch kmlLayer {
            case cmdFlagTracksLayerCamera:
                kmlBuilder = &builders.CameraBuilder{
                    AddBankAngle: !cmdArgs.noBanking,
                    VerboseFlag:  cmdArgs.verboseOperation,
                }
            case cmdFlagTracksLayerPath:
                kmlBuilder = &builders.PathBuilder{
                    Extrude: true,
                    Color:   color.RGBA{R: 217, G: 51, B: 255},
                }
            case cmdFlagTracksLayerPlacemark:
                kmlBuilder = &builders.PlacemarkBuilder{}
            case cmdFlagTracksLayerVector:
                kmlBuilder = &builders.VectorBuilder{}
            default:
                return fmt.Errorf("unrecognized kmlLayer(%s); supported: %v", kmlLayer,
                    strings.Join(cmdFlagTracksLayersSupported, ","))
            }
            kmlBuilders = append(kmlBuilders, kmlBuilder)
        }

        tc := iaeroapi.TracksConverter{
            TailNumber:  cmdArgs.tailNumber,
            CutoffTime:  cmdArgs.cutoffTime,
            FlightCount: cmdArgs.flightCount,
            Api:         aeroApi,
        }

        tracker := &kml.GxTracker{Builders: kmlBuilders}
        aeroKmls, createKmlErr := tc.Convert(tracker)
        if createKmlErr != nil {
            log.Fatalf("couldn't create KMLs: %v", createKmlErr)
            //notreached
        }

        sortedKmlLayersStr := strings.Join(sortedKmlLayers, ",")

        nKmlDocs := len(aeroKmls)
        if cmdArgs.verboseOperation || nKmlDocs > 1 {
            log.Printf("INFO: writing %d %s KML document(s)", nKmlDocs, sortedKmlLayersStr)
        }

        var firstKmlFilename string
        for _, aeroKml := range aeroKmls {
            kmzSaver := &persistence.KmzSaver{
                Saver:  fileBasedAeroApi,
                Assets: aeroKml.KmlAssets,
            }
            flightTimeRange := getTsFromTo(*aeroKml.StartTime, *aeroKml.EndTime)
            kmlFilename := filepath.Join(
                fileBasedAeroApi.ArtifactsDir,
                fmt.Sprintf("fvk_%s-%s-%s.kmz", cmdArgs.tailNumber, flightTimeRange, sortedKmlLayersStr),
            )

            if writeErr := kmzSaver.Save(kmlFilename, aeroKml.KmlDoc); writeErr != nil {
                log.Fatalf("couldn't write output artifact(%s): %v", kmlFilename, writeErr)
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

func parseArgs(cmd *cobra.Command) (cmdArgs tracksCommandArgs, err error) {
    if cmdArgs.verboseOperation, err = cmd.Flags().GetBool(cmdFlagRootVerbose); err != nil {
        return
    }
    if cmdArgs.launchFirstKml, err = cmd.Flags().GetBool(cmdFlagTracksLaunch); err != nil {
        return
    }
    if cmdArgs.noBanking, err = cmd.Flags().GetBool(cmdFlagTracksNoBanking); err != nil {
        return
    }
    if cmdArgs.saveResponses, err = cmd.Flags().GetBool(cmdFlagTracksSaveArtifacts); err != nil {
        return
    }

    if cmdArgs.tailNumber, err = cmd.Flags().GetString(cmdFlagTracksTailNumber); err != nil {
        return
    }
    if cmdArgs.fromSavedResponses, err = cmd.Flags().GetString(cmdFlagTracksFromArtifacts); err != nil {
        return
    }
    if cmdArgs.artifactsDir, err = cmd.Flags().GetString(cmdFlagTracksArtifactsDir); err != nil {
        return
    }
    if cmdArgs.kmlLayers, err = cmd.Flags().GetString(cmdFlagTracksLayers); err != nil {
        return
    }
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

    if cmdArgs.flightCount, err = cmd.Flags().GetInt(cmdFlagTracksFlightCount); err != nil {
        return
    }

    if _, artifactsDirExistsErr := os.Stat(cmdArgs.artifactsDir); os.IsNotExist(artifactsDirExistsErr) {
        err = fmt.Errorf("artifacts directory(%s) not found; either create it or use the '%s' option to change: %w",
            cmdArgs.artifactsDir, cmdFlagTracksArtifactsDir, artifactsDirExistsErr)
        return
    }

    return
}

// GetTsFromTo returns a string representation of a time range using fnPrefixTimestampFormat
// to format the "from" time, and a subsequence of that for the "to" time, with leading common
// prefix removed.  Example:
//
// { 2023010203040506Z, 2023010203050506Z } => "23010203040506-50506Z" ('5' differs with '4' in tsBase)
func getTsFromTo(from, to time.Time) string {
    fromFmt := from.Format(fnPrefixTimestampFormat)[2:]
    toFmt := to.Format(fnPrefixTimestampFormat)[2:]

    i := 0
    for i < len(fromFmt) && i < len(toFmt) && fromFmt[i] == toFmt[i] {
        i++
    }
    return fmt.Sprintf("%s-%s", fromFmt, toFmt[i:])
}

const fnPrefixTimestampFormat = "20060102150405Z"
