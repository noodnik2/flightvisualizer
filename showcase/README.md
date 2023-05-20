# [flightvisualizer] Showcase

This folder "show cases" some examples created by the CLI, including several examples
demonstrating some problems to be fixed.

## Examples in Google Earth Pro

When played in the Google Earth Pro application, the following output, produced with the 
depicted `fviz` CLI command, exhibits the described exemplary characteristics:

### [230518165047Z-N12345-230511231752-01949Z-camera.kmz](230518165047Z-N12345-230511231752-01949Z-camera.kmz)
```shell
$ fviz tracks N12345 --mock --kind camera
```
- At time offset ~`2:15`, this file demonstrates a repeating problem where the "smooth track" suddenly
    starts playing backwards.  A similar effect is demonstrated each time the plane flies over this point
    on its final to the runway.
### [230518171747Z-N12345-230511231752-01949Z-std.kmz](230518171747Z-N12345-230511231752-01949Z-std.kmz)
```shell
$ fviz tracks N12345 --mock --kind std
```
- In addition to the "Camera" view present in the example file described above, several additional "layers"
  are embedded into this "standard" kind of track output, including an "extruded" flight path and floating
  "vectors" depicting some aeronautical performance data present in the original response from AeroAPI.

[flightvisualizer]: https://github.com/noodnik2/flightvisualizer