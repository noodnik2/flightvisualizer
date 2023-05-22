package builders

import (
    gokml "github.com/twpayne/go-kml/v3"
)

// KmlProduct contains the top-level KML model element and the assets it references
type KmlProduct struct {
    Root   gokml.Element
    Assets map[string]any
}

const feetPerMeter = 3.28084

// AeroAlt2Meters converts altitude values emitted by AeroAPI,
// which are expressed in units of 100 feet, into meters
func aeroAlt2Meters(altD100ft float64) float64 {
    feetAgl := altD100ft * 100
    return feetAgl / feetPerMeter
}
