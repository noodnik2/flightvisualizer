package builders

import (
	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type PlacemarkBuilder struct{}

func (*PlacemarkBuilder) Name() string {
	return "Placemark"
}

func (*PlacemarkBuilder) Build(aeroTrackPositions []aeroapi.Position) (*KmlProduct, error) {

	// - https://github.com/twpayne/go-kml/blob/a1a42dcf7ccb20a4b7b88b5bd61178cc14e050fc/kml_test.go#L870
	var positions []gokml.Element

	// - https://developers.google.com/kml/documentation/kmlreference#gxtrack
	for _, position := range aeroTrackPositions {
		positions = append(positions, gokml.When(position.Timestamp))
	}

	// - https://developers.google.com/kml/documentation/kmlreference#gxcoord
	for _, position := range aeroTrackPositions {
		positions = append(positions, gokml.GxCoord(gokml.Coordinate{
			Lon: position.Longitude,
			Lat: position.Latitude,
			Alt: aeroAlt2Meters(position.AltMslD100),
		}))
	}

	track := gokml.GxTrack(positions...)
	placemark := gokml.Placemark(track)
	root := gokml.Folder(
		gokml.Name("Placemark Track"),
		gokml.Description("Flight path track across the ground in a single Placemark"),
		placemark,
	)

	return &KmlProduct{Root: root}, nil
}
