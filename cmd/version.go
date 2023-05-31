package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "pre-release"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf(`Prints the version of %s`, rootCmd.Short),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s\n", Version)
		return nil
	},
}
