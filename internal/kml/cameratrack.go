package kml

import (
	"time"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type CameraTrackBuilder struct {
	AddBankAngle bool
	VerboseFlag  bool
}

func (ctb *CameraTrackBuilder) Build(positions []aeroapi.Position) *KmlProduct {
	var frames []gokml.Element
	nPositions := len(positions)
	aeroapiMathUtil := &aeroapi.Math{}
	flyToMode := gokml.GxFlyToModeBounce // initial "bounce" into tour
	var startTime time.Time

	for i := 0; i < nPositions-1; i++ {

		thisPosition := positions[i]
		nextPosition := positions[i+1]

		if startTime.IsZero() {
			startTime = thisPosition.Timestamp
		}

		var bankAngle float64
		if ctb.AddBankAngle {
			bankAngle = float64(aeroapiMathUtil.GetBankAngle(thisPosition, nextPosition))
		}

		deltaT := nextPosition.Timestamp.Sub(thisPosition.Timestamp)
		const cameraHeightFromWheels = 2
		frames = append(frames, gokml.GxFlyTo(
			gokml.GxDuration(deltaT),
			gokml.GxFlyToMode(flyToMode),
			gokml.Camera(
				gokml.TimeSpan(
					gokml.Begin(startTime),
					gokml.End(nextPosition.Timestamp),
				),
				gokml.Longitude(thisPosition.Longitude),
				gokml.Latitude(thisPosition.Latitude),
				gokml.Altitude(aeroAlt2Meters(thisPosition.AltAglD100)+cameraHeightFromWheels),
				gokml.Heading(thisPosition.Heading),
				gokml.Tilt(80),
				gokml.Roll(-bankAngle),
				gokml.AltitudeMode(gokml.AltitudeModeRelativeToGround),
			)))
		flyToMode = gokml.GxFlyToModeSmooth
	}
	root := gokml.GxTour(
		gokml.Name("Camera View"),
		gokml.Description("First-person view of the flight"),
		gokml.GxPlaylist(frames...),
	)

	return &KmlProduct{Root: root}
}
