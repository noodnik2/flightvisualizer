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

### Proposed Implementation Plan using "AeroAPI"

A use case this morning was to track my wife's flight to D.C.  Using the following APIs I was able
to get periodic updates:

1. Find the `fa_flight_id` value using the `/airports/{id}/flights/departures` API.
   - For this you need to provide the (ICAO or IATA) airport ID and the `{id}` code for the airline,
     and provide a reasonable time range containing the time of the flight's departure.
   - I used a 30-minute time range capturing the actual time of departure somewhere in its middle.
   - Both of these input values can easily be found using FlightAware's web and mobile apps, 
     though if needed a call to `/operators` and `/operators/{id}/flights` could be used,
     though that would be expensive given the number of calls that would be necessary.
   - E.g.:
     - ```shell
       $ curl -X GET "https://aeroapi.flightaware.com/aeroapi/airports/KSFO/flights/departures?airline=ASA&start=2023-07-19T15%3A00%3A00Z&end=2023-07-19T15%3A30%3A00Z&max_pages=2" \
       -H "Accept: application/json; charset=UTF-8" \
       -H "x-apikey:  <AeroApiKey>"
       ```
2. Use the `/flights/{fa_flight_id}/position` API to periodically retrieve the flight's position
   and speed; e.g.:
   - ```shell
     curl -X GET "https://aeroapi.flightaware.com/aeroapi/flights/ASA8-1689607273-airline-443p/position" \
     -H "Accept: application/json; charset=UTF-8" \
     -H "x-apikey: <AeroApiKey>"
     ```
   - Which returns at jpath `$.last_position`:
     - ```json
       "last_position": {
         "fa_flight_id": "ASA8-1689607273-airline-443p",
         "altitude": 370,
         "altitude_change": "-",
         "groundspeed": 512,
         "heading": 101,
         "latitude": 40.14427,
         "longitude": -85.72974,
         "timestamp": "2023-07-19T18:57:05Z",
         "update_type": "A"
       },
       ```
3. The last step could be repeated to generate and feed KML updates to the destination based upon
   the last position and speed using a configured spline function, etc.

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