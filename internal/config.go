package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

type Config struct {
	AeroApiUrl   string `env:"AEROAPI_API_URL,default=https://aeroapi.flightaware.com/aeroapi"`
	ArtifactsDir string `env:"ARTIFACTS_DIR,default=."`
	Verbose      bool   `env:"VERBOSE,default=false"`
	// "required" fields should come at the end; otherwise, the defaults (above) won't be applied when
	// the required values aren't found (that error isn't fatal so we want the defaults to be applied)
	AeroApiKey string `env:"AEROAPI_API_KEY,required" secret:"mask"`
}

const (
	userConfigFilenameEnvVar = "FVIZ_CONFIG_FILE"
	configFile               = ".config/fviz"
)

func GetConfigFilename(verbose bool) string {
	if userConfigFilename := os.Getenv(userConfigFilenameEnvVar); userConfigFilename != "" {
		return userConfigFilename
	}
	const homeDirEnvVarName = "HOME"
	_, ok := os.LookupEnv(homeDirEnvVarName)
	if !ok {
		log.Printf("NOTE: '%s' environment variable not found\n", homeDirEnvVarName)
	}
	homeDir := os.Getenv(homeDirEnvVarName)
	configFilename := filepath.Join(homeDir, configFile)
	if verbose {
		log.Printf("INFO: config file location is '%s'\n", configFilename)
	}
	return configFilename
}

func GetBuildVcsVersion() string {

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	var vcsRevision, vcsTime, vcsModified string

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			vcsRevision = kv.Value
		case "vcs.time":
			vcsTime = kv.Value
		case "vcs.modified":
			if kv.Value == "true" {
				vcsModified = " (modified)"
			}
		}
	}

	return fmt.Sprintf("%s %s%s", vcsTime, vcsRevision[:7], vcsModified)
}
