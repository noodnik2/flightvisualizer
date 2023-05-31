# Log Notes

The original intent of keeping this "scratchpad" sort of file is for leaving temporary notes,
reminders & general internal information related to progress and development of Flight Visualizer.

## Resources
### Symbols
- [Google Map Icons](http://kml4earth.appspot.com/icons.html)
### Forums
- [Google Earth Support](https://support.google.com/earth)
- [Google Maps Support](https://support.google.com/maps)
### Documentation
- [1010MHz Riddle](https://mode-s.org/decode/)

## On the Radar: Related Projects & People
- [traffic](https://github.com/xoolive/traffic)
  - [Project lead](https://www.xoolive.org/)
  - [traffic.js](https://github.com/xoolive/traffic.js)
  - [A journey through aviation data](https://aviationbook.netlify.app)
- [OpenFlights](https://github.com/jpatokal/openflights)
- [Flight Stream](https://callumprentice.github.io/apps/flight_stream/index.html)
- [FlightWise](https://flightwise.com) - appears to be outdated
- [airplanejs](https://github.com/ADSBexchange/airplanejs) - plots nearby ADS-B on a live map
- [atmdata](https://atmdata.github.io/sources/)

### Possible Data Sources?
Some of the identified sources below also support "live" streams, not just historical tracks of past flights.
- [FlightAware's "FireHose"](https://flightaware.com/commercial/firehose/)
- [OpenSky Network](https://opensky-network.org/)
  - [GitHub](https://github.com/openskynetwork/)
  - e.g., start with [`GET /flights/aircraft`](https://openskynetwork.github.io/opensky-api/rest.html#flights-by-aircraft)
  - Check out: ["State Vectors"](https://openskynetwork.github.io/opensky-api/#state-vectors))
- [FlightRadar24](https://www.flightradar24.com/) - supports feeds, real-time, 7-days history (`$+`)
  - [Live map](https://www.flightradar24.com/)
  - Looks like they provide [App Integration](https://www.flightradar24.com/commercial-services/app-integration)
    collaboration on a custom project basis, not as "open" as the others
    - E.g., see [FlightRadarAPI](https://pypi.org/project/FlightRadarAPI/)
- [ADS-B Exchange](https://www.adsbexchange.com/) - supports feeds
  - [Live map](https://globe.adsbexchange.com/)
- [Aviation Edge](https://aviation-edge.com/) - real-time, historical, database download (`$$`)
- [Aviation Stack](https://aviationstack.com/) - real-time, free 100 requests / month (`$$`)
- [Flight Labs](https://www.goflightlabs.com/) - real-time, historical (`$$$`)
- [AeroDataBox](https://aerodatabox.com/) - no flight tracking (`$`)
- [Airlabs](https://airlabs.co/) - Reseller of flight data?
  - Description of their "free" plan seems misleading
  - Are they a reseller? (`$$`?)
  - Seems like a wannabe startup
