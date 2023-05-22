# Flight Visualizer

Uses data retrieved using [AeroAPI] to create a [KML] visualization of a flight,
enabling it to be "played back" in a renderer such as [Google Earth].

All that's needed to identify (a) flight(s) to [AeroAPI] is an aircraft's "tail number," and optionally
a "cutoff time" prior to which to consider.  If no "cutoff time" is given, the current time will be used.

## Packaging

Flight Visualizer is packaged into a standalone CLI command named `fviz`.

## Setup / Configuration

In order to retrieve flight information using [AeroAPI], you'll need to obtain and install an
API Key.  See the [AeroAPI] documentation to learn more about this.

Once you have your API Key:
1. Copy the supplied [`./.env.local.template`](./.env.local.template) template file to `.env.local` in the folder from which you run `fviz`
2. Set the value of the `AEROAPI_API_KEY` configuration parameter in `.env.local`

Once you've done this, you'll be good to go.  Just remember to protect this `.env.local` file going
forward, as it contains a secret.

## Example Invocations

Here's a simple example of invoking `fviz` in which a KML visualization of previously saved
flight tracks is created and immediately "launched."  

Notes:

- No calls to [AeroAPI] are made when using previously saved flights
- Launching the KML file will invoke the application registered for `.kml` files in your operating system
  (such as [Google Earth Pro](https://www.google.com/earth/versions/#earth-pro))

```shell
$ fviz tracks N12345 --fromSavedResponses savedfiles/flights_N605WM-1683866868-adhoc-2431p_track.json --launch
```

Another example below leverages more of the available options:

```shell
$ dist/fviz tracks N605WM --verbose --cutoffTime 2023-05-11T22:40:00-07:00 --kind camera --outputDir kml --saveResponses --flightCount 3
2023/05/19 15:55:07 INFO: requesting from endpoint(/flights/N605WM?&end=2023-05-11T22:40:00-07:00)
2023/05/19 15:55:08 INFO: creating file: savedfiles/230519155507Z-flights_N605WM_end_2023-05-11T22_40_00-07_00.json
2023/05/19 15:55:08 INFO: requesting from endpoint(/flights/N605WM-1683869494-adhoc-1188p/track)
2023/05/19 15:55:08 INFO: creating file: savedfiles/230519155507Z-flights_N605WM-1683869494-adhoc-1188p_track.json
2023/05/19 15:55:08 INFO: requesting from endpoint(/flights/N605WM-1683869013-adhoc-1394p/track)
2023/05/19 15:55:09 INFO: creating file: savedfiles/230519155507Z-flights_N605WM-1683869013-adhoc-1394p_track.json
2023/05/19 15:55:09 INFO: requesting from endpoint(/flights/N605WM-1683866868-adhoc-2431p/track)
2023/05/19 15:55:09 INFO: creating file: savedfiles/230519155507Z-flights_N605WM-1683866868-adhoc-2431p_track.json
2023/05/19 15:55:09 INFO: writing 3 camera KML document(s)
2023/05/19 15:55:09 INFO: creating file: savedfiles/230519155507Z-N605WM-230512053129-5Z-camera.kmz
2023/05/19 15:55:09 INFO: creating file: savedfiles/230519155507Z-N605WM-230512052328-2Z-camera.kmz
2023/05/19 15:55:09 INFO: creating file: savedfiles/230519155507Z-N605WM-230512050306-2Z-camera.kmz
```

This example creates "camera" visualizations for the three flights prior to the stated "cutoff time",
saving them and also the raw response files obtained from [AeroAPI] in the `kml` directory, using
"verbose" operation in order to reveal more about the internal functioning.

[AeroAPI]: https://flightaware.com/commercial/aeroapi
[KML]: https://developers.google.com/kml
[Google Earth]: https://earth.google.com