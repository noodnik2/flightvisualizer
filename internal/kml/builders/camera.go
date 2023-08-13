package builders

import (
	"time"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/pkg/aeroapi"
)

type CameraBuilder struct {
	AddBankAngle bool
	DebugFlag    bool
}

func (ctb *CameraBuilder) Name() string {
	return "Camera"
}

func (ctb *CameraBuilder) Build(positions []aeroapi.Position) (*KmlProduct, error) {
	var frames []gokml.Element
	nPositions := len(positions)
	aeroApiMathUtil := &aeroapi.Math{
		Debug: ctb.DebugFlag,
	}
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
			bankAngle = float64(aeroApiMathUtil.GetBankAngle(thisPosition, nextPosition))
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
				gokml.Altitude(aeroAlt2Meters(thisPosition.AltMslD100)+cameraHeightFromWheels),
				gokml.Heading(thisPosition.Heading),
				gokml.Tilt(80),
				gokml.Roll(-bankAngle),
				gokml.AltitudeMode(gokml.AltitudeModeAbsolute),
			)))
		flyToMode = gokml.GxFlyToModeSmooth
	}
	root := gokml.GxTour(
		gokml.Name("Camera View"),
		gokml.Description("First-person view of the flight"),
		gokml.GxPlaylist(frames...),
	)

	return &KmlProduct{Root: root}, nil
}
