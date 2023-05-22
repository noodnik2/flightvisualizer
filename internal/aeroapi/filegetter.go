package aeroapi

//
//import (
//	"fmt"
//	"log"
//	"strings"
//
//	"github.com/noodnik2/flightvisualizer/internal/persistence"
//	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
//)
//
//type AeroApiFileGetterFactory struct {
//	Verbose     bool
//	FileContext persistence.FileContext
//}
//
//type FlightIdToTrackPath func(string) string
//
//func (mg *AeroApiFileGetterFactory) NewRequester(flightIdsFnPath string, x FlightIdToTrackPath) (aeroapi.GetRequester, error) {
//
//	flightIdsFnRef := flightIdsFnPath
//	flightIdsFilePath, flightIdsJsonBytes, flightIdReaderErr := mg.FileContext.LoadFromFnRef(flightIdsFnRef)
//	if flightIdReaderErr != nil {
//		return nil, fmt.Errorf("from path(%s): %w", flightIdsFilePath, flightIdReaderErr)
//	}
//
//	flights, flightsErr := aeroapi.FlightsFromJson(flightIdsJsonBytes)
//	if flightsErr != nil {
//		return nil, flightsErr
//	}
//
//	var flightIds []string
//	for _, flight := range flights.Flights {
//		flightIds = append(flightIds, flight.FlightId)
//	}
//
//	var flightIndex int
//
//	localGetter := func(endpoint string) ([]byte, error) {
//		if !strings.Contains(endpoint, "/track") {
//			// if it's not the "/track" endpoint, then it must be the "/flights/id"
//			// endpoint for which we've already got the response; so, just return it
//			return flightIdsJsonBytes, nil
//		}
//
//		if flightIndex >= len(flightIds) {
//			return nil, fmt.Errorf("no more mock tracks available")
//		}
//		flightId := flightIds[flightIndex]
//		flightIndex++
//
//		trackFileContext := persistence.FileContext{
//			Reader: mg.FileContext.Reader,
//		}
//		trackFilepath, trackJsonBytes, trackLoadErr := trackFileContext.LoadFromFnRef(x(flightId))
//		if trackLoadErr != nil {
//			return nil, fmt.Errorf("could not load track file(%s): %w", trackFilepath, trackLoadErr)
//		}
//
//		if mg.Verbose {
//			log.Printf("INFO: satisfying request for(%s) with local data from(%s)\n", endpoint, trackFilepath)
//		}
//		return trackJsonBytes, nil
//	}
//
//	return localGetter, nil
//}
