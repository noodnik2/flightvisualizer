package kml

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/internal/kml/builders"
	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

// Track contains the fully-rendered KML document representing a flight,
// assets referenced by that KML document, and some relevant metadata
type Track struct {
	KmlDoc    []byte
	KmlAssets map[string]any
	StartTime *time.Time
	EndTime   *time.Time
}

// TrackGenerator can generate a Track from raw flight position data
type TrackGenerator interface {
	Generate(*aeroapi.Track) (*Track, error)
}

// TrackBuilderEnsemble is a named set of KmlTrackBuilder instances
type TrackBuilderEnsemble struct {
	Builders []builders.KmlTrackBuilder
}

func (gxt *TrackBuilderEnsemble) Generate(aeroTrack *aeroapi.Track) (*Track, error) {

	positions := aeroTrack.Positions
	nPositions := len(positions)
	var fromTime, toTime *time.Time
	if nPositions > 0 {
		fromTime = &positions[0].Timestamp
		toTime = &positions[nPositions-1].Timestamp
	}

	var layerNames []string
	for _, kmlBuilder := range gxt.Builders {
		layerNames = append(layerNames, kmlBuilder.Name())
	}
	mainDocument := gokml.Document(
		gokml.Name(fmt.Sprintf("AeroAPI Flight %s", aeroTrack.FlightId)),
		gokml.Description(fmt.Sprintf("Layers: %s", strings.Join(layerNames, ", "))),
	)

	kmlAssets := make(map[string]any)
	for _, kb := range gxt.Builders {
		kmlThing, buildErr := kb.Build(positions)
		if buildErr != nil {
			fmt.Printf("NOTE: %s\n", buildErr)
			continue
		}
		mainDocument.Append(kmlThing.Root)
		for k, v := range kmlThing.Assets {
			kmlAssets[k] = v
		}
	}
	gxKMLElement := gokml.GxKML(mainDocument)

	var kmlBuilder bytes.Buffer
	if err := gxKMLElement.Write(&kmlBuilder); err != nil {
		return nil, err
	}
	kmlTrack := Track{
		KmlDoc:    kmlBuilder.Bytes(),
		KmlAssets: kmlAssets,
		StartTime: fromTime,
		EndTime:   toTime,
	}
	return &kmlTrack, nil
}
