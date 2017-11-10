package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	"github.com/djthorpe/VideoIntelligence/service"
)

var (
	FlagServiceAccount = flag.String("sa", "", "Service Account JSON")
	FlagDebug          = flag.Bool("debug", false, "Debug")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create the service
	service, err := service.NewServiceFromServiceAccountJSON(*FlagServiceAccount, *FlagDebug)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	}
	fmt.Println(service)

	/*	call := service.Videos.Annotate(&api.GoogleCloudVideointelligenceV1beta1AnnotateVideoRequest{
			Features: []string{"LABEL_DETECTION", "SHOT_CHANGE_DETECTION", "SAFE_SEARCH_DETECTION"},
			InputUri: "gs://cloud-ml-sandbox/video/chicago.mp4",
		})
		if response, err := call.Do(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		} else {
			fmt.Println(response)
		}
	*/
}
