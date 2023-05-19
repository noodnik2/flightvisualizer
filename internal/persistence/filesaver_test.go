package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFilenameCleaner(t *testing.T) {

	type testCaseDef struct {
		name             string
		timestamp        *time.Time
		suffix           string
		endpoint         string
		expectedFragment string
		expectedFilename string
	}

	testTimestamp := time.Date(2023, 5, 12, 14, 41, 0, 0, time.UTC)
	testCases := []testCaseDef{
		{
			name:             "simple case",
			endpoint:         "http://localhost:2324/test?arg=one",
			expectedFragment: "http_localhost_2324_test_arg_one",
			expectedFilename: "http_localhost_2324_test_arg_one",
		},
		{
			name:             "simple case path only",
			endpoint:         "/_test?arg=one ",
			expectedFragment: "test_arg_one",
			expectedFilename: "test_arg_one",
		},
		{
			name:             "simple case with suffix",
			suffix:           ".json",
			endpoint:         "http://localhost:2324/test?arg=one",
			expectedFragment: "http_localhost_2324_test_arg_one",
			expectedFilename: "http_localhost_2324_test_arg_one.json",
		},
		{
			name:             "simple case with timestamp",
			timestamp:        &testTimestamp,
			suffix:           ".txt",
			endpoint:         "http://my.domain/test?arg=one",
			expectedFilename: "230512144100Z-http_my.domain_test_arg_one.txt",
		},
		{
			name:             "truncated case with timestamp and long suffix",
			timestamp:        &testTimestamp,
			suffix:           ".verylongsuffix",
			endpoint:         `https://foreign.service.org:2324/one/two?spice='nice'&that="one two"`,
			expectedFragment: "https_foreign.service.org_2324_one_two_spice_nice_that_one+two",
			expectedFilename: "230512144100Z-https_foreign.service.org_2324_one.verylongsuffix",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			rs := FileSaver{FilenameSuffix: tc.suffix, Timestamp: tc.timestamp}
			fragment := rs.fnFragmentFromEndpoint(tc.endpoint)
			if tc.expectedFragment != "" {
				requirer.Equal(tc.expectedFragment, fragment)
			}
			filename := rs.fnFromFragment(fragment)
			requirer.Equal(tc.expectedFilename, filename)
		})
	}

}
