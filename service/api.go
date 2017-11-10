package service

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	v1 "github.com/djthorpe/VideoIntelligence/videointelligence/v1"
	v1beta2 "github.com/djthorpe/VideoIntelligence/videointelligence/v1beta2"
	oauth2 "golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
)

var (
	ErrorInvalidServiceAccount = errors.New("Invalid Service Account")
)

///////////////////////////////////////////////////////////////////////////////
// STRUCTS

type Service struct {
	videos *v1beta2.Service
	ops    *v1.Service
	names  []string
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns service
func NewServiceFromServiceAccountJSON(filename string, debug bool) (*Service, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	saConfig, err := google.JWTConfigFromJSON(bytes, v1beta2.CloudPlatformScope)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	client := saConfig.Client(getContext(debug))
	if videos, err := v1beta2.New(client); err != nil {
		return nil, err
	} else if ops, err := v1.New(client); err != nil {
		return nil, err
	} else {
		return &Service{videos, ops, make([]string, 0)}, nil
	}
}

func (this *Service) Annotate(uri string) (string, error) {
	call := this.videos.Videos.Annotate(&v1beta2.GoogleCloudVideointelligenceV1beta2AnnotateVideoRequest{
		Features: []string{"LABEL_DETECTION", "SHOT_CHANGE_DETECTION", "EXPLICIT_CONTENT_DETECTION"},
		InputUri: uri,
	})
	if response, err := call.Do(); err != nil {
		return "", err
	} else {
		// Append the operation name into the list of current operations
		this.names = append(this.names, response.Name)
		return response.Name, nil
	}
}

func (this *Service) OperationStatus(name string) (*v1.GoogleLongrunningOperation, error) {
	call := this.ops.Operations.Get(name)
	if response, err := call.Do(); err != nil {
		return nil, err
	} else {
		return response, nil
	}
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
