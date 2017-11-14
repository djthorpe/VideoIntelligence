# VideoIntelligence

Test Code for Google's Video Intelligence API. In order to download the latest version
of the Video Intelligence API, you need to run the following script first:

```
[bash] build/updateapis.sh
```

This will create a "videointelligence" folder with the relevant API's included. Then you can
use the service object in your own software. For example,

```go
service, err := service.NewServiceFromServiceAccountJSON(serviceAccountPath,debug)
if err != nil {
    // Handle error
}
// Following will return the name of the operation
operation, err := service.Annotate(uri,service.ANNOTATE_LABEL)
for {
    // Wait until the operation has completed
    status, err := service.Status(operation)
    if status.Done {
        break
    }
    time.Sleep(1 * time.Second)
}
// Retrieve the labels
labels, err := service.AnnotationLabels(name)
```


There is an example in `vi-analyse.go` which is a command-line tool for analysing
videos and outputs an ASCII table of annotations. You'll need to include the 
Service Account JSON in your home directory, with the name
`.yt-video-intelligence.json` but you can also reference another service account
using the `-sa` flag.

Here is what typical output looks like:

```
[bash] go run vi-analyse.go -shot -explicit -label gs://cloud-ml-sandbox/video/chicago.mp4
Percent Complete=0%
Percent Complete=0%
Percent Complete=33%
Percent Complete=33%
Percent Complete=66%
Percent Complete=100%
+------------------+-----------+-----------------------------+------------+------------+--------------------------+
|       TYPE       |  ENTITY   |         DESCRIPTION         |   START    |    END     |        CONFIDENCE        |
+------------------+-----------+-----------------------------+------------+------------+--------------------------+
| shot             | <nil>     | <nil>                       | 0s         | 38.752016s | <nil>                    |
| shot_label       | /m/06gfj  | road                        | 0s         | 38.752016s |                0.8423761 |
| shot_label       | /m/017r8p | pedestrian > person         | 0s         | 38.752016s |               0.60235006 |
| shot_label       | /m/0btp2  | traffic                     | 0s         | 38.752016s |                 0.941059 |
| shot_label       | /m/01bjv  | bus > vehicle               | 0s         | 38.752016s |               0.86475587 |
| shot_label       | /m/01l7t2 | downtown > city             | 0s         | 38.752016s |               0.98352104 |
| shot_label       | /m/0pg52  | taxi > vehicle              | 0s         | 38.752016s |                0.7113388 |
| shot_label       | /m/02_286 | new york city               | 0s         | 38.752016s |                0.8469136 |
| shot_label       | /m/033j3c | lane > road                 | 0s         | 38.752016s |                 0.522443 |
| shot_label       | /m/01c8br | street > road               | 0s         | 38.752016s |                 0.973952 |
| shot_label       | /m/07bsy  | transport                   | 0s         | 38.752016s |               0.66574144 |
| shot_label       | /m/01n32  | city > geographical feature | 0s         | 38.752016s |                0.9664619 |
| shot_label       | /m/014xcs | pedestrian crossing > road  | 0s         | 38.752016s |                0.4099896 |
| shot_label       | /m/01j0ry | traffic congestion > event  | 0s         | 38.752016s |                0.5472057 |
| shot_label       | /m/039jbq | urban area > city           | 0s         | 38.752016s |                0.9469805 |
| shot_label       | /m/056mk  | metropolis > city           | 0s         | 38.752016s |               0.70255727 |
| shot_label       | /m/012f08 | motor vehicle > vehicle     | 0s         | 38.752016s |                 0.665753 |
| shot_label       | /m/0k4j   | car > vehicle               | 0s         | 38.752016s |                0.9374764 |
| shot_label       | /m/07yv9  | vehicle                     | 0s         | 38.752016s |                0.9199582 |
| shot_label       | /m/01prls | land vehicle > vehicle      | 0s         | 38.752016s |               0.42833233 |
| shot_label       | /m/017kvv | infrastructure              | 0s         | 38.752016s |                0.5684082 |
| shot_label       | /m/0j_s4  | metropolitan area > city    | 0s         | 38.752016s |               0.94977784 |
| shot_label       | /m/0cgh4  | building                    | 0s         | 38.752016s |               0.42035094 |
| segment_label    | /m/012f08 | motor vehicle > vehicle     | 0s         | 38.752016s |                 0.665753 |
| segment_label    | /m/056mk  | metropolis > city           | 0s         | 38.752016s |               0.70255727 |
| segment_label    | /m/039jbq | urban area > city           | 0s         | 38.752016s |                0.9469805 |
| segment_label    | /m/01j0ry | traffic congestion > event  | 0s         | 38.752016s |                0.5472057 |
| segment_label    | /m/014xcs | pedestrian crossing > road  | 0s         | 38.752016s |                0.4099896 |
| segment_label    | /m/0j_s4  | metropolitan area > city    | 0s         | 38.752016s |               0.94977784 |
| segment_label    | /m/0cgh4  | building                    | 0s         | 38.752016s |               0.42035094 |
| segment_label    | /m/017kvv | infrastructure              | 0s         | 38.752016s |                0.5684082 |
| segment_label    | /m/07yv9  | vehicle                     | 0s         | 38.752016s |                0.9199582 |
| segment_label    | /m/01prls | land vehicle > vehicle      | 0s         | 38.752016s |               0.42833233 |
| segment_label    | /m/0k4j   | car > vehicle               | 0s         | 38.752016s |                0.9374764 |
| segment_label    | /m/0pg52  | taxi > vehicle              | 0s         | 38.752016s |                0.7113388 |
| segment_label    | /m/01l7t2 | downtown > city             | 0s         | 38.752016s |               0.98352104 |
| segment_label    | /m/01bjv  | bus > vehicle               | 0s         | 38.752016s |               0.86475587 |
| segment_label    | /m/0btp2  | traffic                     | 0s         | 38.752016s |                 0.941059 |
| segment_label    | /m/017r8p | pedestrian > person         | 0s         | 38.752016s |               0.60235006 |
| segment_label    | /m/06gfj  | road                        | 0s         | 38.752016s |                0.8423761 |
| segment_label    | /m/07bsy  | transport                   | 0s         | 38.752016s |               0.66574144 |
| segment_label    | /m/01c8br | street > road               | 0s         | 38.752016s |                 0.973952 |
| segment_label    | /m/01n32  | city > geographical feature | 0s         | 38.752016s |                0.9664619 |
| segment_label    | /m/033j3c | lane > road                 | 0s         | 38.752016s |                 0.522443 |
| segment_label    | /m/02_286 | new york city               | 0s         | 38.752016s |                0.8469136 |
| explicit_content | <nil>     | <nil>                       | 570.808ms  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 1.381775s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 2.468091s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 3.426006s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 4.356055s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 5.21124s   | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 6.019059s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 7.082135s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 7.941474s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 9.052317s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 10.004125s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 10.871788s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 12.000251s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 12.91709s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 13.939s    | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 14.941261s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 16.031191s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 16.834676s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 17.724766s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 18.82167s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 19.878636s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 20.88244s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 21.8288s   | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 22.653542s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 23.587377s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 24.528324s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 25.639358s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 26.659231s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 27.596029s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 28.406835s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 29.33981s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 30.37907s  | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 31.283454s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 32.365059s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 33.180016s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 34.041546s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 35.085369s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 35.891895s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 37.0588s   | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
| explicit_content | <nil>     | <nil>                       | 38.211111s | <nil>      | LIKELIHOOD_VERY_UNLIKELY |
+------------------+-----------+-----------------------------+------------+------------+--------------------------+
```

There's currently a bug where the explicit content doesn't always come through.
