package kml

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/kmlflight/pkg/aeroapi"
)

// VectorBuilder - builds a KML folder of Placemarks revealing significant
// details of the track coordinates received from AeroAPI, including:
//
// => Location - Placemark's location
// => Altitude - Placemark's altitude
// => Heading - direction of an arrow representing the Placemark
// => Groundspeed - reflected by magnitude / size of the arrow
//
// Additional sets of Placemarks are used to reveal secondary information
// calculated from the track data (e.g., "imputed" values).
type VectorBuilder struct{}

const vectorArrowRelPath = "internal/kml/images/blue_fast_arrow.png"

func (sb *VectorBuilder) Build(aeroTrackPositions []aeroapi.Position) *KmlProduct {

	vectorArrowAbsPath, fpErr := filepath.Abs(vectorArrowRelPath)
	if fpErr != nil {
		log.Fatalf("couldn't get absolute path of(%s): %v", vectorArrowRelPath, fpErr)
		//notreached
	}
	vectorArrowPngBytes, rfErr := os.ReadFile(vectorArrowAbsPath)
	if rfErr != nil {
		log.Fatalf("couldn't read file(%s): %v", vectorArrowAbsPath, rfErr)
		//notreached
	}

	vectorArrowHref := filepath.Base(vectorArrowAbsPath)
	thing := &KmlProduct{
		Assets: map[string]any{
			vectorArrowHref: vectorArrowPngBytes,
		},
	}

	aeroapiMathUtil := &aeroapi.Math{}
	var positionReferences []gokml.Element

	for i := 0; i < len(aeroTrackPositions)-1; i++ {
		thisPosition := aeroTrackPositions[i]
		nextPosition := aeroTrackPositions[i+1]

		// the "reported" placemarks plot info received directly from AeroAPI
		reportedDescription := getReportedDescription(thisPosition, nextPosition)
		reportedPlacemark := styledPlacemark{
			styleName:    fmt.Sprintf("heading-icon%d", i),
			balloonText:  fmt.Sprintf("<h1>Reported</h1>%s", reportedDescription),
			styleColor:   color.RGBA{R: 255, G: 255, B: 0, A: 255},
			heading:      thisPosition.Heading,
			gs:           thisPosition.GsKnots,
			iconImageUrl: vectorArrowHref,
			position:     thisPosition,
		}
		positionReferences = append(positionReferences, reportedPlacemark.getElements()...)

		// the "imputed" placemarks plot info calculated indirectly from AeroAPI data
		geoHeading := aeroapiMathUtil.GetGeoBearing(thisPosition, nextPosition)
		geoGsKnots := aeroapiMathUtil.GetGeoGsKnots(thisPosition, nextPosition)
		imputedPlacemark := styledPlacemark{
			styleName:    fmt.Sprintf("bearing-icon%d", i),
			balloonText:  fmt.Sprintf("<h1>Imputed</h1>%s%s", getImputedDescription(geoHeading, geoGsKnots), reportedDescription),
			styleColor:   color.RGBA{R: 0, G: 255, B: 255, A: 255},
			heading:      float64(geoHeading),
			gs:           geoGsKnots,
			iconImageUrl: vectorArrowHref,
			position:     thisPosition,
		}
		positionReferences = append(positionReferences, imputedPlacemark.getElements()...)
	}

	thing.Root = gokml.Folder(
		gokml.Name("Vector Track"),
		gokml.Description("Vectors along flight path reflecting performance data"),
	).
		Append(positionReferences...)

	return thing
}

type styledPlacemark struct {
	styleName    string
	balloonText  string
	styleColor   color.Color
	heading      float64
	gs           float64
	iconImageUrl string
	position     aeroapi.Position
}

func (sp *styledPlacemark) getElements() []gokml.Element {
	return []gokml.Element{
		gokml.Style(
			gokml.IconStyle(
				gokml.Color(sp.styleColor),
				gokml.Icon(gokml.Href(sp.iconImageUrl)),
				gokml.Heading(sp.heading-90),
				gokml.Scale(sp.gs/100),
			),
			gokml.BalloonStyle(gokml.Text(sp.balloonText)),
		).WithID(sp.styleName),
		gokml.Placemark(
			gokml.StyleURL("#"+sp.styleName),
			gokml.Point(
				gokml.AltitudeMode(gokml.AltitudeModeRelativeToGround),
				gokml.Coordinates(
					gokml.Coordinate{
						Lon: sp.position.Longitude,
						Lat: sp.position.Latitude,
						Alt: aeroAlt2Meters(sp.position.AltAglD100),
					},
				),
			),
		),
	}
}

func getImputedDescription(geoHeading aeroapi.Degrees, geoGsKnots float64) string {
	return fmt.Sprintf(`<h2>Imputed From Location Change</h2>
		<ul>
			<li>Heading: %.1fº</li>
			<li>Groundspeed: %.1fkt</li>
		</ul>`,
		geoHeading,
		geoGsKnots,
	)
}

func getReportedDescription(thisPosition aeroapi.Position, nextPosition aeroapi.Position) string {
	return fmt.Sprintf(`<h2>Reported by AeroAPI</h2>
		<h3>This Location</h3>
		<ul>
			<li>Time: %v</li>
			<li>Location: %v</li>
			<li>Altitude: %.0f'</li>
			<li>Heading: %.1fº</li>
			<li>Groundspeed: %.1fkt</li>
		</ul>
		<h3>Next Location</h3>
		<ul>
			<li>Time: %v</li>
			<li>Location: %v</li>
			<li>Altitude: %.0f'</li>
			<li>Heading: %.1fº</li>
			<li>Groundspeed: %.1fkt</li>
		</ul>`,
		thisPosition.Timestamp.Format(time.RFC3339),
		[]float64{thisPosition.Latitude, thisPosition.Longitude},
		thisPosition.AltAglD100*100,
		thisPosition.Heading,
		thisPosition.GsKnots,
		thisPosition.Timestamp.Format(time.RFC3339),
		[]float64{nextPosition.Latitude, nextPosition.Longitude},
		nextPosition.AltAglD100*100,
		thisPosition.Heading,
		thisPosition.GsKnots,
	)
}
