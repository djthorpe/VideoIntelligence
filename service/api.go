package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	v1 "github.com/djthorpe/VideoIntelligence/videointelligence/v1"
	v1beta2 "github.com/djthorpe/VideoIntelligence/videointelligence/v1beta2"
	oauth2 "golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC STRUCTS

// Service defines the client for the Video Intellgence API
type Service struct {
	videos *v1beta2.Service
	ops    *v1.Service
	status map[string]*Status
}

// Status defines the current operation status
type Status struct {
	Name     string
	Uri      string
	Done     bool
	Type     []AnnotationType
	Progress map[AnnotationType]*Progress
	Updated  time.Time
}

// Progress defines progress on the annotation operations
type Progress struct {
	Done       bool
	Percent    int64
	StartTime  time.Time
	UpdateTime time.Time
}

// Response defines the response of the annotation operations
type Response struct{}

// AnnotationType are the types of annotations
type AnnotationType uint

///////////////////////////////////////////////////////////////////////////////
// PRIVATE STRUCTS

type my_time struct {
	time.Time
}

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ANNOTATION_NONE                            = 0
	ANNOTATION_LABEL            AnnotationType = 1 << iota
	ANNOTATION_SHOT_CHANGE      AnnotationType = 1 << iota
	ANNOTATION_EXPLICIT_CONTENT AnnotationType = 1 << iota
)

const (
	// Duration which to fetch remote status
	duration_CACHE_EXPIRY time.Duration = 1 * time.Minute
)

var (
	ErrInvalidServiceAccount = errors.New("Invalid Service Account")
	ErrNotFound              = errors.New("Not found")
	ErrInProgress            = errors.New("In progress")
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
			Name:     response.Name,
			Uri:      uri,
			Type:     annotateTypeArray(flags),
			Progress: make(map[AnnotationType]*Progress, 3),
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
			startTime, _ := time.Parse(time.RFC3339Nano, statusDetail.StartTime)
			updateTime, _ := time.Parse(time.RFC3339Nano, statusDetail.UpdateTime)
			done := (statusDetail.ProgressPercent == 100)
			status.Progress[annotationType] = &Progress{
				done,
				statusDetail.ProgressPercent,
				startTime,
				updateTime,
			}
		}
		// set the done flag and updated flag
		status.Done = response.Done
		status.Updated = time.Now()
		return status, nil
	}
}

func (this *Service) ExplicitAnnotations(name string) (*Response, error) {
	// Get progress on operation, if not completed then return error
	if progress, err := this.getCachedProgress(name, ANNOTATION_LABEL, duration_CACHE_EXPIRY); err != nil {
		return nil, err
	} else if progress.Done == false {
		return nil, ErrInProgress
	}
	// TODO
	return &Response{}, nil
}

func (this *Service) LabelAnnotations(name string) (*Response, error) {
	// Get progress on operation, if not completed then return error
	if progress, err := this.getCachedProgress(name, ANNOTATION_LABEL, duration_CACHE_EXPIRY); err != nil {
		return nil, err
	} else if progress.Done == false {
		return nil, ErrInProgress
	}
	// TODO
	return &Response{}, nil
}

func (this *Service) ShotChangeAnnotations(name string) (*Response, error) {
	// Get progress on operation, if not completed then return error
	if progress, err := this.getCachedProgress(name, ANNOTATION_LABEL, duration_CACHE_EXPIRY); err != nil {
		return nil, err
	} else if progress.Done == false {
		return nil, ErrInProgress
	}
	// TODO
	return &Response{}, nil
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

// getCachedStatus returns a status object, or a refresh the status object if
// hasn't been updated in a while
func (this *Service) getCachedStatus(name string, cacheExpiry time.Duration) (*Status, error) {
	var (
		status *Status
		exists bool
		fetch  bool
		err    error
	)

	// Set fetch flag which indicates we need ro re-fetch the status object
	if status, exists = this.status[name]; exists == false {
		fetch = true
	} else if time.Now().Sub(status.Updated) >= cacheExpiry {
		fetch = true
	}

	// Fetch the status object (side-effect is that it's set in 'this')
	if fetch {
		if status, err = this.Status(name); err != nil {
			return nil, err
		}
	}

	// Return the status object
	return status, nil
}

func (this *Service) getCachedProgress(name string, annotationType AnnotationType, cacheExpiry time.Duration) (*Progress, error) {
	if status, err := this.getCachedStatus(name, cacheExpiry); err != nil {
		return nil, err
	} else if progress, exists := status.Progress[annotationType]; exists == false {
		return nil, ErrNotFound
	} else {
		return progress, nil
	}
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

func (s Status) String() string {
	progress := make([]string, 0, 3)
	for _, annotationType := range []AnnotationType{ANNOTATION_LABEL, ANNOTATION_EXPLICIT_CONTENT, ANNOTATION_SHOT_CHANGE} {
		annotationProgress, exists := s.Progress[annotationType]
		if exists {
			progress = append(progress, fmt.Sprintf("%v=%v", annotationType, annotationProgress))
		}
	}
	return fmt.Sprintf("{ name=%v uri=%v updated=%v done=%v progress=[ %v ] }", s.Name, s.Uri, my_time{s.Updated}, s.Done, strings.Join(progress, ","))
}

func (p Progress) String() string {
	return fmt.Sprintf("{ done=%v percent=%v start=%v updated=%v }", p.Done, p.Percent, my_time{p.StartTime}, my_time{p.UpdateTime})
}

func (t my_time) String() string {
	difference := -t.Sub(time.Now())
	if difference.Seconds() < 60 {
		return fmt.Sprintf("%v secs ago", int(difference.Seconds()))
	}
	if difference.Minutes() < 90 {
		return fmt.Sprintf("%v mins ago", int(difference.Minutes()))
	}
	return fmt.Sprintf("%v hours ago", int(difference.Hours()))
}
