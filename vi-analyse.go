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
	FlagServiceAccount  = flag.String("sa", ".yt-video-intelligence.json", "Service Account JSON")
	FlagDebug           = flag.Bool("debug", false, "Debug")
	FlagShotChange      = flag.Bool("shot", false, "Annotate for Shot Changes")
	FlagLabel           = flag.Bool("label", true, "Annotate for Labels")
	FlagExplicitContent = flag.Bool("explicit", false, "Annotate for Explicit Content")
)

func filenameToAbsolute(filename string) (string, error) {
	path, exists := util.ResolvePath(filename, util.UserDir())
	if exists {
		return path, nil
	} else {
		return "", fmt.Errorf("Missing file: %s", filename)
	}
}

func annotationFlags() service.AnnotationType {
	var flags service.AnnotationType
	if *FlagShotChange {
		flags |= service.ANNOTATION_SHOT_CHANGE
	}
	if *FlagLabel {
		flags |= service.ANNOTATION_LABEL
	}
	if *FlagExplicitContent {
		flags |= service.ANNOTATION_EXPLICIT_CONTENT
	}
	return flags
}

func outputResponseEntity(output *util.Output, t string, label *service.EntityAnnotation) {
	for i := range label.Segments {
		entity_description := label.Entity.Description
		for i, category := range label.Categories {
			if i == 0 {
				entity_description = entity_description + " > "
			} else {
				entity_description = entity_description + ", "
			}
			entity_description = entity_description + category.Description
		}
		if i == 0 {
			output.AppendMap(map[string]interface{}{
				"type":        t,
				"value":       label,
				"entity":      label.Entity.EntityId,
				"description": entity_description,
				"start":       label.Segments[0].StartOffset,
				"end":         label.Segments[0].EndOffset,
				"confidence":  label.Segments[0].Confidence,
			})
		} else {
			output.AppendMap(map[string]interface{}{
				"type":        "",
				"value":       "",
				"entity":      "",
				"description": "",
				"start":       label.Segments[i].StartOffset,
				"end":         label.Segments[i].EndOffset,
				"confidence":  label.Segments[i].Confidence,
			})
		}
	}
}

func outputResponse(status *service.Status, output *util.Output) {
	if len(status.Annotations.Shots) > 0 {
		for _, shot := range status.Annotations.Shots {
			output.AppendMap(map[string]interface{}{
				"type":  "shot",
				"start": shot.StartOffset,
				"end":   shot.EndOffset,
				"value": shot,
			})
		}
	}
	if len(status.Annotations.ShotLabels) > 0 {
		for _, label := range status.Annotations.ShotLabels {
			outputResponseEntity(output, "shot_label", label)
		}
	}
	if len(status.Annotations.SegmentLabels) > 0 {
		for _, label := range status.Annotations.SegmentLabels {
			outputResponseEntity(output, "segment_label", label)
		}
	}
	if len(status.Annotations.ExplicitContent) > 0 {
		for _, annotation := range status.Annotations.ExplicitContent {
			output.AppendMap(map[string]interface{}{
				"type":       "explicit_content",
				"start":      annotation.Offset,
				"confidence": annotation.Likelihood,
				"value":      annotation,
			})
		}
	}

	output.RenderASCII()
}

func runMain(api *service.Service, uris []string) error {
	if len(uris) == 0 {
		return errors.New("Missing uri arguments")
	}

	// Gather output
	output := util.NewOutput("type", "entity", "description", "start", "end", "confidence")
	// Add value column if debug
	if *FlagDebug {
		output.AddColumns("value")
	}

	// Return Annotate result for each URI
	for _, uri := range uris {
		if operation, err := api.Annotate(uri, annotationFlags()); err != nil {
			return err
		} else {
			for {
				status, err := api.Status(operation)
				if err != nil {
					return err
				}
				fmt.Printf("Percent Complete=%v%%\n", status.PercentComplete())
				if status.Done {
					outputResponse(status, output)
					break
				} else {
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
