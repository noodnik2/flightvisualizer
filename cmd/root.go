package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Flight Visualizer",
	Long:  `Generates KML visualizations of flight data retrieved from FlightAware's AeroAPI`,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const cmdFlagRootVerbose = "verbose"

func init() {
	rootCmd.PersistentFlags().BoolP(cmdFlagRootVerbose, "v", false, "Enables 'verbose' operation")
}
