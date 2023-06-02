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
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(editConfigCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: fmt.Sprintf("Facilitates configuration of the %s", rootCmd.Short),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return showConfigCmd.RunE(showConfigCmd, nil)
		}
		return errors.New("invalid argument(s)")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {

		configFilename := internal.GetConfigFilename()
		var config internal.Config
		if loadConfigErr := configurator.LoadConfig(configFilename, &config); loadConfigErr != nil {
			log.Printf("NOTE: %v\n", loadConfigErr)
		}
		items, getInfoErr := configurator.GetConfigEnvItems(config)
		if getInfoErr != nil {
			log.Fatalf("ERROR: %v\n", getInfoErr)
		}
		log.Printf("from '%s':\n", configFilename)
		for _, configItem := range items {
			var val any
			if configItem.Secret != "" {
				val = "*******"
			} else {
				val = configItem.Val
			}
			fmt.Printf("%s: %v\n", configItem.Name, val)
		}
		return nil
	},
}

var editConfigCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edits current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {

		configFilename := internal.GetConfigFilename()
		configDir := filepath.Dir(configFilename)
		if _, err := os.Stat(configDir); err != nil {
			log.Fatalf("NOTE: config file directory(%s) not found; please create it\n", configDir)
		}

		log.Printf("editing '%s':\n", configFilename)

		var config internal.Config
		if loadConfigErr := configurator.LoadConfig(configFilename, &config); loadConfigErr != nil {
			log.Printf("NOTE: %v\n", loadConfigErr)
		}

		if editErr := configurator.EditConfig(&config); editErr != nil {
			log.Fatalf("ERROR: %v\n", editErr)
		}

		if saveErr := configurator.SaveConfig(configFilename, config); saveErr != nil {
			log.Fatalf("ERROR: %v\n", saveErr)
		}

		return nil
	},
}
