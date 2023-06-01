package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
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

func GetBuildVersion() string {

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	goVersion := info.GoVersion
	path := info.Path
	mainVersion := info.Main.Version
	var goArch, goOs, vcs, vcsRevision, vcsTime, vcsModified string

	for _, kv := range info.Settings {
		switch kv.Key {
		case "GOARCH":
			goArch = kv.Value
		case "GOOS":
			goOs = kv.Value
		case "vcs":
			vcs = kv.Value
		case "vcs.revision":
			vcsRevision = kv.Value
		case "vcs.time":
			vcsTime = kv.Value
		case "vcs.modified":
			vcsModified = kv.Value
		}
	}

	return fmt.Sprintf("%s:%s:%s/%s:%s:%s:%s:%s:%s", vcs, vcsRevision[:7], vcsTime, vcsModified,
		goVersion, goArch, goOs, mainVersion, path)
}
