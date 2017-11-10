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
videos and outputs a CSV of annotations.


