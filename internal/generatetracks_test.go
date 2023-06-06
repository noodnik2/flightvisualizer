package internal

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracksCommandArgs_GenerateTracks(t *testing.T) {

	testCases := []struct {
		name             string
		layers           []string
		expectedEnsemble []string
		expectedErrors   []string
	}{
		{
			name: "no layers",
		},
		{
			name:             "single layer",
			layers:           []string{TracksLayerCamera},
			expectedEnsemble: []string{TracksLayerCamera},
		},
		{
			name:             "all layers, random order",
			layers:           []string{TracksLayerPath, TracksLayerVector, TracksLayerPlacemark, TracksLayerCamera},
			expectedEnsemble: []string{TracksLayerCamera, TracksLayerPath, TracksLayerPlacemark, TracksLayerVector},
		},
		{
			name:             "all layers - with duplicates",
			layers:           []string{TracksLayerPlacemark, TracksLayerCamera, TracksLayerPath, TracksLayerPlacemark, TracksLayerPath, TracksLayerVector},
			expectedEnsemble: []string{TracksLayerCamera, TracksLayerPath, TracksLayerPlacemark, TracksLayerVector},
		},
		{
			name:           "unrecognized layer",
			layers:         []string{TracksLayerPlacemark, TracksLayerCamera, "unrecognized", TracksLayerPlacemark, TracksLayerPath, TracksLayerVector},
			expectedErrors: []string{"unrecognized kmlLayer(unrecognized)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			tca := TracksCommandArgs{}
			generator, err := tca.newKmlTrackGenerator(tc.layers)
			if tc.expectedErrors != nil {
				requirer.Error(err)
				for _, expectedErr := range tc.expectedErrors {
					requirer.Contains(err.Error(), expectedErr)
				}
				return
			}
			requirer.NoError(err)
			requirer.Equal(len(tc.expectedEnsemble), len(generator.Builders))
			requirer.Equal(strings.Join(tc.expectedEnsemble, "-"), generator.Name)
		})
	}
}

func TestTracksCommandArgs_NewTrackFactory(t *testing.T) {

	testCases := []struct {
		name              string
		artifactsFilename string
		expectedErrors    []string
		expectedFnName    string
	}{
		{
			name:           "no artifacts file",
			expectedFnName: "multiTrackRemoteFactory",
		},
		{
			name:              "flight ids file",
			artifactsFilename: "some_dir/fvf_file.json",
			expectedFnName:    "multiTrackArtifactFactory",
		},
		{
			name:              "track file",
			artifactsFilename: "fvt_file.json",
			expectedFnName:    "singleTrackArtifactFactory",
		},
		{
			name:              "unrecognized artifact",
			artifactsFilename: "unknown_artifact.json",
			expectedErrors:    []string{"unrecognized", "unknown_artifact.json"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			tca := TracksCommandArgs{
				FromArtifacts: tc.artifactsFilename,
			}
			trackFactory, tfErr := tca.newTrackFactory()
			if tc.expectedErrors != nil {
				requirer.Error(tfErr)
				for _, expectedErr := range tc.expectedErrors {
					requirer.Contains(tfErr.Error(), expectedErr)
				}
				return
			}
			requirer.NoError(tfErr)
			// Yes, I know - it's a bit banal to look at the function names as "proof" things are working
			// TODO refactor to allow better tests....
			trackFactoryFnName := runtime.FuncForPC(reflect.ValueOf(trackFactory).Pointer()).Name()
			requirer.Contains(trackFactoryFnName, tc.expectedFnName)
		})
	}
}

func TestTracksCommandArgs_SaveKmlTracks(t *testing.T) {

	testCases := []struct {
		name string
	}{
		{
			name: "no input",
		},
		// TODO add meaningful tests
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			tca := TracksCommandArgs{}
			tracks, err := tca.saveKmlTracks(nil, "")
			requirer.NoError(err)
			requirer.Empty(tracks)
		})
	}
}
