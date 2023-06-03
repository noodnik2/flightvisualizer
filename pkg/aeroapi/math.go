package aeroapi

import (
	"math"
	"time"
)

type Math struct {
	Verbose bool
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

	// TODO consider adding a "debug" flag in order to re-enable this sometimes-helpful output
	//if u.Verbose {
	//	log.Printf("heading(%f), bank angle(%f), deltaT(%v) deltaH(%f), turnRate(%v), groundspeed(%f)\n",
	//		fromPosition.Heading, bankAngle, deltaT, deltaH, turnRate, fromPosition.GsKnots)
	//}

	return bankAngle
}

// GetGeoGsKnots calculates and reports the apparent average ground speed used
// navigate the straight line distance between two geolocations
func (u *Math) GetGeoGsKnots(fromPosition, toPosition Position) float64 {
	from := geoCoord{latitude: fromPosition.Latitude, longitude: fromPosition.Longitude}
	to := geoCoord{latitude: toPosition.Latitude, longitude: toPosition.Longitude}
	nm := from.nmTo(to)
	deltaT := toPosition.Timestamp.Sub(fromPosition.Timestamp)

	deltaTHours := f(deltaT) / f(time.Hour)
	return nm / deltaTHours
}

// GetGeoBearing calculates and reports the apparent compass bearing
// (0 <= bearing < 360) needed to arrive at a new geolocation
func (u *Math) GetGeoBearing(fromPosition, toPosition Position) Degrees {
	from := geoCoord{latitude: fromPosition.Latitude, longitude: fromPosition.Longitude}
	to := geoCoord{latitude: toPosition.Latitude, longitude: toPosition.Longitude}
	return from.bearingTo(to)
}

// TODO use https://github.com/twpayne/go-kml/blob/master/sphere/sphere.go instead
func (from geoCoord) bearingTo(to geoCoord) Degrees {

	rLatFrom := degToRad(Degrees(from.latitude))
	rLonFrom := degToRad(Degrees(from.longitude))
	rLatTo := degToRad(Degrees(to.latitude))
	rLonTo := degToRad(Degrees(to.longitude))

	deltaRlon := rLonTo - rLonFrom

	rBearing := math.Atan2(
		math.Sin(f(deltaRlon))*math.Cos(f(rLatTo)),
		math.Cos(f(rLatFrom))*math.Sin(f(rLatTo))-math.Sin(f(rLatFrom))*math.Cos(f(rLatTo))*math.Cos(f(deltaRlon)),
	)

	dBearing := radToDeg(Radians(rBearing))
	if dBearing < 0 {
		dBearing += 360
	}
	dBearing = Degrees(math.Mod(f(dBearing), 360))
	return dBearing
}

// TODO use https://github.com/twpayne/go-kml/blob/master/sphere/sphere.go instead
func (from geoCoord) nmTo(to geoCoord) float64 {

	latFrom := degToRad(Degrees(from.latitude))
	lonFrom := degToRad(Degrees(from.longitude))
	latTo := degToRad(Degrees(to.latitude))
	lonTo := degToRad(Degrees(to.longitude))

	deltaLat := latTo - latFrom
	deltaLon := lonTo - lonFrom

	// https://www.movable-type.co.uk/scripts/latlong.html
	a := math.Pow(math.Sin(f(deltaLat)/2), 2) +
		math.Cos(f(latFrom))*math.Cos(f(latTo))*math.Pow(math.Sin(f(deltaLon)/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	const earthRadiusKm = 6371
	const kilometersPerNauticalMile = 1.852

	return earthRadiusKm * c / kilometersPerNauticalMile
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

func degToRad(deg Degrees) Radians {
	return Radians(deg * math.Pi / 180)
}

func radToDeg(rad Radians) Degrees {
	return Degrees(rad * 180 / math.Pi)
}

type geoCoord struct {
	latitude  float64
	longitude float64
}
