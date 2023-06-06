package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/noodnik2/configurator"
	"github.com/spf13/cobra"

	"github.com/noodnik2/flightvisualizer/internal"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(editConfigCmd)
}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Shows current configuration",
	Version: rootCmd.Version,
	RunE: func(cmd *cobra.Command, args []string) (returnErr error) {

		if cmd.Flags().NArg() != 0 {
			return errors.New("invalid syntax")
		}

		var configFilename string
		if configFilename, returnErr = getArgs(cmd); returnErr == nil {
			cmd.SilenceUsage = true
			returnErr = showConfig(configFilename)
		}

		return
	},
}

var editConfigCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edits current configuration",
	Version: rootCmd.Version,
	RunE: func(cmd *cobra.Command, args []string) (returnErr error) {

		var configFilename string
		if configFilename, returnErr = getArgs(cmd); returnErr == nil {
			cmd.SilenceUsage = true
			if returnErr = editConfig(configFilename); returnErr == nil {
				returnErr = showConfig(configFilename)
			}
		}

		return
	},
}

func getArgs(cmd *cobra.Command) (configFilename string, returnErr error) {
	var verbose bool
	if verbose, returnErr = cmd.Flags().GetBool(cmdFlagRootVerbose); returnErr == nil {
		configFilename = internal.GetConfigFilename(verbose)
	}
	return
}

func showConfig(configFilename string) (returnErr error) {

	var config internal.Config
	if loadConfigErr := configurator.LoadConfig(configFilename, &config); loadConfigErr != nil {
		log.Printf("NOTE: %v\n", loadConfigErr)
	}

	var items []configurator.ConfigEnvItem
	items, returnErr = configurator.GetConfigEnvItems(config)
	if returnErr != nil {
		return
	}

	log.Printf("INFO: current configuration:\n")
	for _, configItem := range items {
		var val any
		if configItem.Secret != "" {
			val = "*******"
		} else {
			val = configItem.Val
		}
		log.Printf("%s: %v\n", configItem.Name, val)
	}
	return
}

func editConfig(configFilename string) (returnErr error) {

	configDir := filepath.Dir(configFilename)
	if _, configDirErr := os.Stat(configDir); configDirErr != nil {
		return fmt.Errorf("config file directory(%s) not found", configDir)
	}

	var config internal.Config
	if loadConfigErr := configurator.LoadConfig(configFilename, &config); loadConfigErr != nil {
		log.Printf("NOTE: %v\n", loadConfigErr)
	}

	log.Printf("INFO: editing configuration:\n")
	if returnErr = configurator.EditConfig(&config); returnErr != nil {
		return
	}

	log.Printf("INFO: saving configuration\n")
	returnErr = configurator.SaveConfig(configFilename, config)
	return
}