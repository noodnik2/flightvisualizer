package internal

import (
	"os"
	"path/filepath"
)

type Config struct {
	AeroApiUrl   string `env:"AEROAPI_API_URL,default=https://aeroapi.flightaware.com/aeroapi"`
	ArtifactsDir string `env:"ARTIFACTS_DIR,default=."`
	Verbose      bool   `env:"VERBOSE,default=false"`
	// "required" fields should come at the end; otherwise, the defaults (above) won't be applied
	AeroApiKey string `env:"AEROAPI_API_KEY,required" secret:"mask"`
}

const configFile = ".config/fviz"

func GetConfigFilename() string {
	homeDir := os.Getenv("HOME")
	configFilename := filepath.Join(homeDir, configFile)
	return configFilename
}
