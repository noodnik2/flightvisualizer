package aeroapi

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestNewHttpAeroApi(t *testing.T) {

    type testCaseDef struct {
        name         string
        status       int
        response     []byte
        expectedErrs []string
    }

    testCases := []testCaseDef{
        {
            name:     "successful case",
            response: []byte("test server response"),
        },
        {
            name:         "error case",
            status:       http.StatusForbidden,
            response:     nil,
            expectedErrs: []string{"Forbidden", "couldn't"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            requirer := require.New(t)
            svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if tc.status != 0 {
                    w.WriteHeader(tc.status)
                }
                _, err := w.Write(tc.response)
                requirer.NoError(err)
            }))
            defer svr.Close()

            api := &HttpAeroApi{ApiUrl: svr.URL}
            actualResponse, errGet := api.Load("irrelevant")
            requirer.Equal(tc.response, actualResponse)
            if len(tc.expectedErrs) > 0 {
                // expecting an error
                requirer.Error(errGet)
                for _, errText := range tc.expectedErrs {
                    requirer.Contains(errGet.Error(), errText)
                }
            } else {
                // not expecting an error
                requirer.NoError(errGet)
            }
        })
    }

}
