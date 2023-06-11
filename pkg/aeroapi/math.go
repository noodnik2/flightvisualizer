package aeroapi

import (
	"log"
	"math"
	"time"

	"github.com/twpayne/go-kml/v3"
	"github.com/twpayne/go-kml/v3/sphere"
)

type Math struct {
	Debug bool
}

type f = float64
type Degrees float64
type Radians float64

// GetBankAngle calculates and reports, using a loose heuristic, a reasonable "bank angle"
// that could be used by a general aviation aircraft to achieve the observed change in heading
// between one reported position and another.
func (u *Math) GetBankAngle(fromPosition, toPosition Position) Degrees {

	isLeftHandTurn := !isRightHandTurn(fromPosition.Heading, toPosition.Heading)
	deltaH := toPosition.Heading - fromPosition.Heading
	if isLeftHandTurn {
		// turning left so magnitude of heading change is
		// current heading (larger) less next heading
		deltaH = -deltaH
	}
	if deltaH < 0 {
		// rationalize negative values
		deltaH += 360
	}
	if isLeftHandTurn {
		// express left hand turns as a negative offset from current heading
		deltaH = -deltaH
	}

	deltaT := toPosition.Timestamp.Sub(fromPosition.Timestamp)
	turnRate := deltaH * f(time.Second) / f(deltaT)

	// use prior heuristic:
	// see https://github.com/noodnik2/MSFS2020-PilotPathRecorder/blob/8bda00905b8566d103e32d0e76b01941ce066c92/FS2020PlanePath/FlightDataGenerator.cs#L74
	newRawPlaneBankAngle := ((fromPosition.GsKnots / 10) + 7) * turnRate / 3
	bankAngle := rationalizeBankAngle(Degrees(newRawPlaneBankAngle))

	if u.Debug {
		log.Printf("heading(%f), bank angle(%f), deltaT(%v) deltaH(%f), turnRate(%v), groundspeed(%f)\n",
			fromPosition.Heading, bankAngle, deltaT, deltaH, turnRate, fromPosition.GsKnots)
	}

	return bankAngle
}

// GetGeoGsKnots calculates and reports the apparent average ground speed used
// navigate the straight line distance between two geolocations
func (u *Math) GetGeoGsKnots(fromPosition, toPosition Position) float64 {

	// get distance
	const earthRadiusKm = 6371
	earth := sphere.T{R: earthRadiusKm}
	kilometers := earth.HaversineDistance(
		kml.Coordinate{Lon: fromPosition.Longitude, Lat: fromPosition.Latitude},
		kml.Coordinate{Lon: toPosition.Longitude, Lat: toPosition.Latitude})

	const kilometersPerNauticalMile = 1.852
	nauticalMiles := kilometers / kilometersPerNauticalMile

	// get time
	deltaT := toPosition.Timestamp.Sub(fromPosition.Timestamp)
	deltaTHours := f(deltaT) / f(time.Hour)

	// get ground speed
	return nauticalMiles / deltaTHours
}

// GetGeoBearing calculates and reports the apparent compass bearing
// (0 <= bearing < 360) needed to arrive at a new geolocation
func (u *Math) GetGeoBearing(fromPosition, toPosition Position) Degrees {
	var earth sphere.T
	initialBearing := earth.InitialBearingTo(
		kml.Coordinate{Lon: fromPosition.Longitude, Lat: fromPosition.Latitude},
		kml.Coordinate{Lon: toPosition.Longitude, Lat: toPosition.Latitude})
	if initialBearing < 0 {
		initialBearing += 360
	}
	return Degrees(math.Mod(initialBearing, 360))
}

func isRightHandTurn(origHdg, destHdg float64) bool {
	switchDir := math.Abs(destHdg-origHdg) > 180
	destGreaterThanOrig := destHdg > origHdg
	return destGreaterThanOrig != switchDir
}

func rationalizeBankAngle(bankAngle Degrees) Degrees {
	if bankAngle < -60 {
		return -60
	}
	if bankAngle > 60 {
		return 60
	}
	return bankAngle
}
