package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/noodnik2/flightvisualizer/internal"
)

var Version = "v0.0.2"

var rootCmd = &cobra.Command{
	Use:     os.Args[0],
	Short:   "Flight Visualizer CLI",
	Long:    `Generates visualizations of flight data`,
	Version: fmt.Sprintf("%s %s", Version, internal.GetBuildVcsVersion()),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		debugFlag, debugFlagErr := cmd.Flags().GetBool(cmdFlagRootDebug)
		if debugFlagErr != nil {
			return debugFlagErr
		}

		if debugFlag {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				cmd.SilenceUsage = true
				return fmt.Errorf("no build info available")
			}
			log.Printf("GetBuildInfo(%v)\n", info)
			return nil
		}

		return errors.New("invalid syntax")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// err has already been presented to the user, if needed
		os.Exit(1)
	}
}

const cmdFlagRootVerbose = "verbose"
const cmdFlagRootDebug = "debug"

func init() {
	rootCmd.PersistentFlags().BoolP(cmdFlagRootVerbose, "v", false, "Enables 'verbose' operation")
	rootCmd.PersistentFlags().BoolP(cmdFlagRootDebug, "d", false, "Enables 'debug' operation")
}
