# Wishlist

Below are some ideas and visions for future enhancements to FlightVisualizer. 

## Support for Multiple Sources & Targets

Initially, Flight Visualizer is "hard-coded" to support flight track data from a single
source - [AeroAPI] - and generate visualizations using a single target - [Google Earth KML].

The application could relatively easily be extended to support _multiple sources_ of flight
tracking data and/or _multiple targets_ for orchestrating flight visualizations (see some
potential candidates identified below).

In addition to expanding the number of sources and target output formats, conformance to
data exchange standards (such as [AIXM](https://www.aixm.aero/) and [FIXM](https://fixm.aero/))
is a goal to increase interoperability with, and usefulness of Flight Visualizer. 

## Support for "Live Visualization"

Flight Visualizer operates now only as an offline batch processor that simply transforms a static
set of data (e.g., responses previously received from [AeroAPI]) into another format (e.g., [KML]).

This vision is to extend the application so that it can do "on the fly" transformations of an incoming
stream of flight track data, feeding it in "real-time" to a rendering engine (such as [Google Earth] or
[Microsoft Flight Simulator]), to enable visualizations of flight(s) in progress. 

See the ["Live Cam"](https://github.com/noodnik2/MSFS2020-PilotPathRecorder/blob/master/README-kmlcam.md)
and [related video demo](https://github.com/noodnik2/MSFS2020-PilotPathRecorder/blob/master/README-kmlcam-QandA.md)
of how this is already working for a real-time view in [Google Earth] of a flight being flown in 
[Microsoft Flight Simulator].

## Candidate Sources of Flight Track Data

Some potential candidates for additional _sources_ of flight data include:

### Flight Planning
Visualization of a planned flight can be both entertaining and a prudent step to take prior to takeoff; and
is an obvious use case behind products such as [ForeFlight 3D Views](https://youtu.be/Fl7ubeiEB2o).

If implemented, this vision would enable FlightVisualizer to accept flight plan track data (such as 
waypoints, planned altitudes / headings / airspeed, etc.) to give the user a pre-flight visualization
and familiarization with a planned or proposed flight.

Support for popular flight plan data formats such as Garmin's [FlightPlan](https://www8.garmin.com/xmlschemas/FlightPlanv1.xsd)
or [other popular formats](https://www.littlenavmap.org/manuals/littlenavmap/release/latest/en/FLIGHTPLANFMT.html)
would be needed.

### Flight Track Logs

These products and websites appear to offer alternate sources of flight track data which could be supported
by Flight Visualizer:

- [Aireon](https://aireon.com)
- [FlightRadar24](https://www.flightradar24.com)
- [ADS-B Exchange](https://www.adsbexchange.com/)
- [OpenSky Network](https://opensky-network.org)
- [FlightAware's "FireHose"](https://flightaware.com/commercial/firehose/)
- [Microsoft Flight Simulator] or [X-Plane]
  - Yes, that's right!  This would enable flight simulator fans to render alternate visualizations
    (see the discussion below about an existing application which
    [_already does this in MSFS2020_](https://github.com/noodnik2/MSFS2020-PilotPathRecorder/blob/master/README-kmlcam.md))

## Candidate Visualization Targets

In addition to [KML], other output formats and technologies should be embraced in order to expand
the type(s) of visualization(s) supported by Flight Visualizer.  Some examples include:

- [GeoJSON](https://geojson.org/)
- Other data formats supported by geospatial (GIS) or Virtual Reality (VR) applications according to future use cases,
  potentially including:
  - Geography Markup Language ([GML])
  - GPS Exchange Format ([GPX])

[AeroAPI]: https://flightaware.com/aeroapi
[Google Earth]: https://www.google.com/earth/index.html
[KML]:https://en.wikipedia.org/wiki/Keyhole_Markup_Language
[Google Earth KML]: https://developers.google.com/kml
[Microsoft Flight Simulator]: https://www.flightsimulator.com/
[X-Plane]: https://www.x-plane.com/
[GML]: https://www.ogc.org/standard/gml/
[GPX]: https://en.wikipedia.org/wiki/GPS_Exchange_Format