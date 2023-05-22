package builders

import (
    "image/color"

    gokml "github.com/twpayne/go-kml/v3"

    "github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

// PathBuilder builds the visible "path" track, and optionally its extrusion to the ground
type PathBuilder struct {
    Color   color.Color
    Extrude bool
}

func (pb *PathBuilder) Name() string {
    return "Path"
}

func (pb *PathBuilder) Build(aeroTrackPositions []aeroapi.Position) *KmlProduct {

    var coordinates []gokml.Coordinate
    for _, position := range aeroTrackPositions {
        coordinates = append(coordinates, gokml.Coordinate{
            Lon: position.Longitude,
            Lat: position.Latitude,
            Alt: aeroAlt2Meters(position.AltAglD100),
        })
    }

    lc := func(a uint8) color.RGBA {
        r, g, b, _ := pb.Color.RGBA()
        return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a}
    }

    flightStyle := gokml.Style(
        gokml.LineStyle(
            gokml.Color(lc(127)),
            gokml.Width(3),
        ),
        gokml.PolyStyle(gokml.Color(lc(63))),
    ).WithID("FlightStyle")

    lineString := gokml.LineString(
        gokml.AltitudeMode(gokml.AltitudeModeRelativeToGround),
        gokml.Extrude(pb.Extrude),
        gokml.Coordinates(coordinates...),
    )

    flightLine := gokml.Placemark(
        gokml.StyleURL("#FlightStyle"),
        lineString,
    )

    mainFolder := gokml.Folder(
        gokml.Name("Path Track"),
        gokml.Description("Visible flight path, optionally extruded to the ground"),
        flightStyle,
        flightLine,
    )

    return &KmlProduct{Root: mainFolder}
}
