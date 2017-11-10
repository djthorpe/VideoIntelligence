package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/VideoIntelligence/service"
	"github.com/djthorpe/VideoIntelligence/util"
)

var (
	FlagServiceAccount = flag.String("sa", ".yt-video-intelligence.json", "Service Account JSON")
	FlagDebug          = flag.Bool("debug", false, "Debug")
)

func filenameToAbsolute(filename string) (string, error) {
	path, exists := util.ResolvePath(filename, util.UserDir())
	if exists {
		return path, nil
	} else {
		return "", fmt.Errorf("Missing file: %s", filename)
	}
}

func runMain(api *service.Service, uris []string) error {
	if len(uris) == 0 {
		return errors.New("Missing uri arguments")
	}

	// Return Annotate result for each URI
	for _, uri := range uris {
		if operation, err := api.Annotate(uri, service.ANNOTATION_EXPLICIT_CONTENT); err != nil {
			return err
		} else {
			for {
				if response, err := api.Status(operation); err != nil {
					return err
				} else {
					fmt.Println("Response=", response)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}

	// success
	return nil
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Obtain the filename (if relative path, then make it absolute relative to home folder)
	if serviceAccountPath, err := filenameToAbsolute(*FlagServiceAccount); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	} else if service, err := service.NewServiceFromServiceAccountJSON(serviceAccountPath, *FlagDebug); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	} else if err := runMain(service, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	}
}
