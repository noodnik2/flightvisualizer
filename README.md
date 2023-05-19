# KmlFlight

Uses [AeroAPI]() to create a KML track log of a flight, enabling it to be
"played back" in a KML renderer such as "Google Earth."

## Usage

```shell
$ kmlflight N5322J 2023-05-09T12:00Z-2023-05-09T13:00Z
```

## Implementation

I decided to try out [Cobra](https://cobra.dev/) for building this app as a CLI.
I've never used this framework before, and it seemed worthwhile to try out & learn
due to its apparently popularity.

### Logic

1. Set up the (KML) output stream(s)
2. Look up the [FlightAware] flight identifiers for the specified tail number
   within the specified time range.
3. For each flight:
   1. Prepare the KML output according to options; e.g.:
      - Separate output file?
      - Any needed prologue, etc.
   2. Render the track log into the (KML) output stream
4. Close the (KML) output stream(s)

### Details

Refreshing my memory of KML, I slowly start to recall the relevant schema and
data model objects, e.g.:

- [NetworkLink](https://developers.google.com/kml/documentation/kmlreference#networklink)
- [gx:MultiTrack](https://developers.google.com/kml/documentation/kmlreference#gxmultitrack)

Looking more into the [go-kml](https://github.com/twpayne/go-kml) library I'm planning to use,
[this test](https://github.com/twpayne/go-kml/blob/a1a42dcf7ccb20a4b7b88b5bd61178cc14e050fc/kml_test.go#L870)
seemed relevant and useful for reference.

## Actual Examples

### May 11th 5:40pm PDT
```shell
$ dist/kmlflight create N605WM 2023-05-11T11:00:00Z 2023-05-11T12:30:00Z
```

This request came back empty; however, FlightAware reflects
[a single flight within this time frame](https://flightaware.com/live/flight/N605WM/history/20230511/2305Z/KSQL/KHWD):

```text
Takeoff 04:10PM PDT
Landing 05:20PM PDT
```

Here's the `curl` for this which indeed reveals that the API returns 0 flights 
during this time period.  Something's amiss - _but, what??_

```text
curl -X GET "https://aeroapi.flightaware.com/aeroapi/flights/N605WM?start=2023-05-11T11%3A00%3A00Z&end=2023-05-11T12%3A30%3A00Z" \
 -H "Accept: application/json; charset=UTF-8" \
 -H "x-apikey: mykey"
```

Turns out my calculation of UTC from local was wrong!  Here's the re-submission,
letting `go` do the time zone math for me:

```text
$ dist/kmlflight create N605WM 2023-05-11T16:00:00-07:00 2023-05-11T17:30:00-07:00
```

It works!!  Though that retrieved two flights (??), separating them into single flights
and removing the trailing document marker (why??) produced a file which successfully
"played" in Google Earth by flying a pin around the track!

### Next

Flying a pin around the track isn't want I want though.  I want a "first person" perspective,
as though the viewer is actually flying in the plane.  Perhaps this implies different "modes"
need to be introduced.  In "first-person" mode (for example), this might generate a KML file
using directives such as `<LookAt>` and `<FlyTo>`, etc.

#### Related

- [Tours in KML: Animating Camera](https://mapsplatform.googleblog.com/2009/04/tours-in-kml-animating-camera-and.html)
- [Touring in KML](https://developers.google.com/kml/documentation/touring)
- [How to View a Route in 3D](https://www.plotaroute.com/tip/4/how-to-view-a-route-in-3d-in-google-earth)

## May 14th

```shell
$ dist/kmlflight create N605WM -v -m --kind birdseye --flightCount 1 
2023/05/14 09:53:00 NOTE: using mock data
2023/05/14 09:53:00 NOTE: satisfying request for(/flights/N605WM) with mock data from(testfixtures/aeroapi-flight-id.json)
2023/05/14 09:53:00 NOTE: satisfying request for(/flights/N5322J-1683690340-adhoc-1256p/track) with mock data from(testfixtures/aeroapi-flight-id-track.json)
2023/05/14 09:53:00 writing 1 birdseye KML document(s)
2023/05/14 09:53:00 saving to(N605WM-230510034535-153Z)
2023/05/14 09:53:00 saving file: 230514095300Z-N605WM-230510034535-153Z.kml
$ open 230514095300Z-N605WM-230510034535-153Z.kml
```

The above "sort of" works well.  Problem is it starts in outer space then plunges into the ocean for a fast swim,
then suddenly pops out already having taken off, flying upwind!  Very strange.

## See Also

### FlightAware
- [FlightAware Search](https://flightaware.com) - Search by Flight or Route
- [Mobile Tail# Search](https://flightaware.com/mobile/)

### KML 
- [KML Reference](https://developers.google.com/kml/documentation/kmlreference)
- [MultiTrackKmlGenerator.cs](https://github.com/noodnik2/MSFS2020-PilotPathRecorder/blob/8bda00905b8566d103e32d0e76b01941ce066c92/FS2020PlanePath/MultiTrackKmlGenerator.cs#L45)
- [Google Earth: Debugging guide](https://developers.google.com/earth-engine/guides/debugging)
- [Relevant Example](https://mapsplatform.googleblog.com/2010/07/making-tracks-new-kml-extensions-in.html)
- [More interesting reading](https://www.endpointdev.com/blog/2013/04/creating-smooth-flight-paths-in-google/)
- [aeroapi-py](../aeroapi-py)
