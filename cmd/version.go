package cmd

import (
	"fmt"
	"github.com/noodnik2/flightvisualizer/internal"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf(`Prints the version of %s`, rootCmd.Short),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := internal.GetBuildVersion()
		if version == "" {
			version = "unknown"
		}
		fmt.Printf("%s\n", version)
		return nil
	},
}
