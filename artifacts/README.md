# Flight Visualizer Artifacts

This folder provides Flight Visualizer with:
- A default output location - e.g. for newly generated KML, or saved responses
- A default input location - e.g., for re-using previously saved external responses  

## Included Examples

Several "pre-loaded" artifacts are supplied in this folder, enabling demonstration of
functionality without the need to invoke the external AeroAPI service, as demonstrated below. 

Assuming both the saved "flight" (e.g., `fvf_N335SP_cutoff-20230523T220000Z.json`) and
its referenced "track" (`fvt_N335SP-1684868329-adhoc-760p.json`) artifact file(s) are
present, you can generate and launch a visualization of a saved artifact in this folder
as demonstrated in the following example:

```shell
$ fviz tracks --fromArtifacts fvf_N335SP_cutoff-20230523T220000Z.json --launch
2023/05/25 16:43:58 INFO: reading from file(artifacts/fvf_N335SP_cutoff-20230523T220000Z.json)
2023/05/25 16:43:58 INFO: reading from file(artifacts/fvt_N335SP-1684874159-adhoc-1864p.json)
2023/05/25 16:43:58 INFO: saving to file(artifacts/fvk__230523203550Z-22102Z_camera-path-vector.kmz)
2023/05/25 16:43:58 INFO: Launching 'artifacts/fvk__230523203550Z-22102Z_camera-path-vector.kmz' 
```

Assuming [Google Earth Pro] is installed properly, the flight visualization should immediately
load and start displaying.  Alternatively, the KML visualization (output) `.kmz` file can be
loaded directly into the renderer.  For example, on a MacOs, this can be done from the command-
line - i.e.:

```shell
$ open artifacts/fvk__230523203550Z-22102Z_camera-path-vector.kmz
```

Or by loading the file directly the file into the renderer (e.g., any version of [Google Earth]):

[Google Earth Pro]: https://www.google.com/earth/versions/#earth-pro
[Google Earth]: https://www.google.com/earth/versions
