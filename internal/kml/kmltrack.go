package kml

import (
	"bytes"
	"fmt"
	"time"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

// KmlTrack contains the fully-rendered KML document representing a flight,
// assets referenced by that KML document, and some relevant metadata
type KmlTrack struct {
	KmlDoc    []byte
	KmlAssets map[string]any
	StartTime *time.Time
	EndTime   *time.Time
}

// KmlTracker can generate a KmlTrack from raw flight position data
type KmlTracker interface {
	Generate(*aeroapi.Track) (*KmlTrack, error)
}

// KmlProduct contains the top-level KML model element and the assets it references
type KmlProduct struct {
	Root   gokml.Element
	Assets map[string]any
}

// GxKmlBuilder can build a KmlProduct from a list of raw AeroAPI positions
type GxKmlBuilder interface {
	Build(positions []aeroapi.Position) *KmlProduct
}

// GxTracker is a named set of GxKmlBuilder instances
type GxTracker struct {
	Name     string
	Builders []GxKmlBuilder
}

func (gxt *GxTracker) Generate(aeroTrack *aeroapi.Track) (*KmlTrack, error) {

	positions := aeroTrack.Positions
	nPositions := len(positions)
	var fromTime, toTime *time.Time
	if nPositions > 0 {
		fromTime = &positions[0].Timestamp
		toTime = &positions[nPositions-1].Timestamp
	}

	mainDocument := gokml.Document(
		gokml.Name(gxt.Name),
		gokml.Description(fmt.Sprintf("%s depiction of AeroAPI Flight %s", gxt.Name, aeroTrack.FlightId)),
	)

	kmlAssets := make(map[string]any)
	for _, kb := range gxt.Builders {
		kmlThing := kb.Build(positions)
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
	kmlTrack := KmlTrack{
		KmlDoc:    kmlBuilder.Bytes(),
		KmlAssets: kmlAssets,
		StartTime: fromTime,
		EndTime:   toTime,
	}
	return &kmlTrack, nil
}

const feetPerMeter = 3.28084

// AeroAlt2Meters converts altitude values emitted by AeroAPI,
// which are expressed in units of 100 feet, into meters
func aeroAlt2Meters(altD100ft float64) float64 {
	feetAgl := altD100ft * 100
	return feetAgl / feetPerMeter
}
