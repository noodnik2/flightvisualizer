# Log Notes

## May 18th

```shell
$ dist/fviz tracks --tailNumber N4468N --cutoffTime 2023-04-23T13:00:00-04:00 -s --layers camera,path,vector --flightCount 3
```

## May 22nd

```shell
$ dist/fviz tracks --tailNumber N6189Q --cutoffTime 2023-05-18T20:40:00-04:00 -s --layers camera,path,vector --flightCount 3
2023/05/22 13:14:15 INFO: requesting from endpoint(/flights/N6189Q?&end=2023-05-18T20:40:00-04:00)
2023/05/22 13:14:16 INFO: requesting from endpoint(/flights/N6189Q-1684452005-adhoc-1107p/track)
2023/05/22 13:14:16 INFO: requesting from endpoint(/flights/N6189Q-1684113448-adhoc-1454p/track)
2023/05/22 13:14:17 INFO: requesting from endpoint(/flights/N6189Q-1684108895-adhoc-1665p/track)
```

```shell
$ dist/fviz tracks --tailNumber N6189Q --cutoffTime 2023-05-18T20:40:00-04:00 -s -v --layers camera,path --flightCount 1
2023/05/22 15:52:01 INFO: requesting from endpoint(/flights/N6189Q?&end=2023-05-18T20:40:00-04:00)
2023/05/22 15:52:02 INFO: requesting from endpoint(/flights/N6189Q-1684452005-adhoc-1107p/track)
2023/05/22 15:52:02 INFO: writing 1 camera,path KML document(s)
```