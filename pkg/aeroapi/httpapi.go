package aeroapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpAeroApi struct {
	Verbose bool
	ApiKey  string
	ApiUrl  string
}

func (c *HttpAeroApi) GetFlightIdsRef(tailNumber string, cutoffTime time.Time) string {
	endpoint := fmt.Sprintf("/flights/%s", tailNumber)
	if !cutoffTime.IsZero() {
		endpoint += fmt.Sprintf("?&end=%s", cutoffTime.Format(time.RFC3339))
	}
	return endpoint
}

func (c *HttpAeroApi) GetTrackForFlightUri(flightId string) string {
	endpoint := fmt.Sprintf("/flights/%s/track", flightId)
	return endpoint
}

func (c *HttpAeroApi) Load(endpoint string) ([]byte, error) {
	const pathSep = "/"
	requestUrl := fmt.Sprintf("%s%s%s", strings.TrimRight(c.ApiUrl, pathSep), pathSep, strings.TrimLeft(endpoint, pathSep))
	if c.Verbose {
		log.Printf("INFO: requesting from endpoint(%s)\n", endpoint)
	}
	req, buildReqErr := http.NewRequest("GET", requestUrl, nil)
	if buildReqErr != nil {
		return nil, newApiError("create request", requestUrl, buildReqErr)
	}

	client := &http.Client{}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-apikey", c.ApiKey)
	resp, issueReqErr := client.Do(req)
	if issueReqErr != nil {
		return nil, newApiError("issue request", requestUrl, issueReqErr)
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Printf("WARNING: error closing request body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		responsePayload, _ := io.ReadAll(resp.Body)
		responseErr := fmt.Errorf("statusCode(%d), status(%s), body(%s)", resp.StatusCode, resp.Status, string(responsePayload))
		return nil, newApiError("get successful response", requestUrl, responseErr)
	}

	responsePayload, readResponseBodyErr := io.ReadAll(resp.Body)
	if readResponseBodyErr != nil {
		return nil, newApiError("read response body", requestUrl, readResponseBodyErr)
	}

	return responsePayload, nil
}

func newApiError(what, where string, err error) error {
	return fmt.Errorf("couldn't %s for %s: %w", what, where, err)
}
