package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	v1 "github.com/djthorpe/VideoIntelligence/videointelligence/v1"
	v1beta2 "github.com/djthorpe/VideoIntelligence/videointelligence/v1beta2"
	oauth2 "golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
)

///////////////////////////////////////////////////////////////////////////////
// STRUCTS

type Service struct {
	videos *v1beta2.Service
	ops    *v1.Service
	status map[string]*Status
}

type Status struct {
	Uri      string
	Type     []AnnotationType
	Done     bool
	progress map[AnnotationType]*v1beta2.GoogleCloudVideointelligenceV1VideoAnnotationProgress
}

// AnnotationType are the types of annotations
type AnnotationType uint

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ANNOTATION_NONE                            = 0
	ANNOTATION_LABEL            AnnotationType = 1 << iota
	ANNOTATION_SHOT_CHANGE      AnnotationType = 1 << iota
	ANNOTATION_EXPLICIT_CONTENT AnnotationType = 1 << iota
)

var (
	ErrInvalidServiceAccount = errors.New("Invalid Service Account")
	ErrNotFound              = errors.New("Not Found")
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewServiceFromServiceAccountJSON returns service object and error given
// the filename to the Service Account JSON file which can be downloaded from the
// Google Developer Console
func NewServiceFromServiceAccountJSON(filename string, debug bool) (*Service, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ErrInvalidServiceAccount
	}
	saConfig, err := google.JWTConfigFromJSON(bytes, v1beta2.CloudPlatformScope)
	if err != nil {
		return nil, ErrInvalidServiceAccount
	}
	client := saConfig.Client(getContext(debug))
	if videos, err := v1beta2.New(client); err != nil {
		return nil, err
	} else if ops, err := v1.New(client); err != nil {
		return nil, err
	} else {
		return &Service{videos, ops, make(map[string]*Status)}, nil
	}
}

// Annotate will kick of the annotation process, and provide a unique ID on return
// for the annotation process. You can then use "OperationResponse" to return the
// result of the annotation when done.
func (this *Service) Annotate(uri string, flags AnnotationType) (string, error) {
	call := this.videos.Videos.Annotate(&v1beta2.GoogleCloudVideointelligenceV1beta2AnnotateVideoRequest{
		Features: annotateFlagArray(flags),
		InputUri: uri,
	})
	if response, err := call.Do(); err != nil {
		return "", err
	} else {
		// Append the operation name into the list of current operations
		this.status[response.Name] = &Status{
			Uri:      uri,
			Type:     annotateTypeArray(flags),
			progress: make(map[AnnotationType]*v1beta2.GoogleCloudVideointelligenceV1VideoAnnotationProgress, 3),
		}
		return response.Name, nil
	}
}

func (this *Service) Status(name string) (*Status, error) {
	status, exists := this.status[name]
	if exists == false {
		return nil, ErrNotFound
	}
	call := this.ops.Operations.Get(name)
	if response, err := call.Do(); err != nil {
		return nil, err
	} else {
		var progress v1beta2.GoogleCloudVideointelligenceV1AnnotateVideoProgress
		if err := json.Unmarshal(response.Metadata, &progress); err != nil {
			return nil, err
		}
		// decode the status codes
		for i, statusDetail := range progress.AnnotationProgress {
			annotationType := status.Type[i]
			status.progress[annotationType] = statusDetail
		}
		// set the done flag
		status.Done = response.Done
		return status, nil
	}
}

func (this *Service) ExplicitAnnotations(name string) {
	// TODO
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Returns context
func getContext(debug bool) context.Context {
	ctx := context.Background()
	if debug {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &LogTransport{http.DefaultTransport},
		})
	}
	return ctx
}

// Returns array of annotation flags as a string
func annotateFlagArray(flags AnnotationType) []string {
	flagArray := make([]string, 0, 3)
	if flags&ANNOTATION_LABEL != ANNOTATION_NONE {
		flagArray = append(flagArray, "LABEL_DETECTION")
	}
	if flags&ANNOTATION_SHOT_CHANGE != ANNOTATION_NONE {
		flagArray = append(flagArray, "SHOT_CHANGE_DETECTION")
	}
	if flags&ANNOTATION_EXPLICIT_CONTENT != ANNOTATION_NONE {
		flagArray = append(flagArray, "EXPLICIT_CONTENT_DETECTION")
	}
	return flagArray
}

// Returns array of annotation flags as a string
func annotateTypeArray(flags AnnotationType) []AnnotationType {
	typeArray := make([]AnnotationType, 0, 3)
	if flags&ANNOTATION_LABEL != ANNOTATION_NONE {
		typeArray = append(typeArray, ANNOTATION_LABEL)
	}
	if flags&ANNOTATION_SHOT_CHANGE != ANNOTATION_NONE {
		typeArray = append(typeArray, ANNOTATION_SHOT_CHANGE)
	}
	if flags&ANNOTATION_EXPLICIT_CONTENT != ANNOTATION_NONE {
		typeArray = append(typeArray, ANNOTATION_EXPLICIT_CONTENT)
	}
	return typeArray
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t AnnotationType) String() string {
	switch t {
	case ANNOTATION_LABEL:
		return "ANNOTATION_LABEL"
	case ANNOTATION_SHOT_CHANGE:
		return "ANNOTATION_SHOT_CHANGE"
	case ANNOTATION_EXPLICIT_CONTENT:
		return "ANNOTATION_EXPLICIT_CONTENT"
	default:
		return "ANNOTATION_NONE"
	}
}
