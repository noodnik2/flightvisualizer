package cmd

import (
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "github.com/spf13/cobra"

    "github.com/noodnik2/flightvisualizer/internal"
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

var cmdFlagTracksLayersDefault = []string{internal.TracksLayerCamera, internal.TracksLayerPath, internal.TracksLayerVector}

func init() {
    rootCmd.AddCommand(versionCmd)
    versionCmd.Flags().StringP(cmdFlagTracksArtifactsDir, "a", "artifacts", "Directory to save or load artifacts")
    versionCmd.Flags().BoolP(cmdFlagTracksNoBanking, "b", false, "Disable banking heuristic calculations")
    versionCmd.Flags().IntP(cmdFlagTracksFlightCount, "c", 0, "Count of (most recent) flights to consider (0=unlimited)")
    versionCmd.Flags().StringP(cmdFlagTracksFromArtifacts, "f", "", "Use saved responses instead of querying AeroAPI")
    versionCmd.Flags().StringP(cmdFlagTracksLayers, "l", strings.Join(cmdFlagTracksLayersDefault, ","), "Layer(s) of the KML depiction to create")
    versionCmd.Flags().StringP(cmdFlagTracksTailNumber, "n", "", "Tail number identifier")
    versionCmd.Flags().BoolP(cmdFlagTracksLaunch, "o", false, "Open the KML visualization of the most recent flight retrieved")
    versionCmd.Flags().BoolP(cmdFlagTracksSaveArtifacts, "s", false, "Save responses from AeroAPI requests")
    versionCmd.Flags().StringP(cmdFlagTracksCutoffTime, "t", "", "Cut off time for flight(s) to consider")
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

        return internal.GenerateTracks(cmdArgs)
    },
}

func parseArgs(cmd *cobra.Command) (cmdArgs internal.TracksCommandArgs, err error) {
    if cmdArgs.VerboseOperation, err = cmd.Flags().GetBool(cmdFlagRootVerbose); err != nil {
        return
    }
    if cmdArgs.LaunchFirstKml, err = cmd.Flags().GetBool(cmdFlagTracksLaunch); err != nil {
        return
    }
    if cmdArgs.NoBanking, err = cmd.Flags().GetBool(cmdFlagTracksNoBanking); err != nil {
        return
    }
    if cmdArgs.SaveResponses, err = cmd.Flags().GetBool(cmdFlagTracksSaveArtifacts); err != nil {
        return
    }

    if cmdArgs.TailNumber, err = cmd.Flags().GetString(cmdFlagTracksTailNumber); err != nil {
        return
    }
    if cmdArgs.FromArtifacts, err = cmd.Flags().GetString(cmdFlagTracksFromArtifacts); err != nil {
        return
    }
    if cmdArgs.ArtifactsDir, err = cmd.Flags().GetString(cmdFlagTracksArtifactsDir); err != nil {
        return
    }
    if cmdArgs.KmlLayers, err = cmd.Flags().GetString(cmdFlagTracksLayers); err != nil {
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
        cmdArgs.CutoffTime = toTime
    }

    if cmdArgs.FlightCount, err = cmd.Flags().GetInt(cmdFlagTracksFlightCount); err != nil {
        return
    }

    if _, artifactsDirExistsErr := os.Stat(cmdArgs.ArtifactsDir); os.IsNotExist(artifactsDirExistsErr) {
        err = fmt.Errorf("artifacts directory(%s) not found; either create it or use the '%s' option to change: %w",
            cmdArgs.ArtifactsDir, cmdFlagTracksArtifactsDir, artifactsDirExistsErr)
        return
    }

    if cmdArgs.TailNumber == "" && cmdArgs.FromArtifacts == "" {
        err = fmt.Errorf("required option missing; one of {'%s', '%s'} required", cmdFlagTracksTailNumber, cmdFlagTracksFromArtifacts)
        return
    }

    // warn user of implications of option combinations by invoking knowledge of downstream semantics
    if cmdArgs.FromArtifacts != "" {
        if cmdArgs.SaveResponses { // no reason to save artifacts when we're reading from artifacts
            incompatibleOptions(cmdFlagTracksSaveArtifacts, cmdFlagTracksFromArtifacts)
        }
        if cmdArgs.TailNumber != "" { // tail number is inherent to saved artifact being used
            incompatibleOptions(cmdFlagTracksTailNumber, cmdFlagTracksFromArtifacts)
        }
        if !cmdArgs.CutoffTime.IsZero() { // cutoff time is inherent to saved artifact being used
            incompatibleOptions(cmdFlagTracksCutoffTime, cmdFlagTracksFromArtifacts)
        }
        if cmdArgs.FlightCount != 0 { // flight count is unused when reading from saved artifact
            incompatibleOptions(cmdFlagTracksFlightCount, cmdFlagTracksFromArtifacts)
        }
    }

    return
}

func incompatibleOptions(option1, option2 string) {
    log.Printf("NOTE: ignoring '%s' option; incompatible with '%s'", option1, option2)
}
