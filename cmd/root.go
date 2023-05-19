package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kmlflight",
	Short: "KML track log generator",
	Long: `Generates a KML file representing a track log obtained from FlightAware's AeroAPI.`,
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


