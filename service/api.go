package service

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
	api "google.golang.org/api/videointelligence/v1beta1"
)

var (
	ErrorInvalidServiceAccount = errors.New("Invalid Service Account")
)

///////////////////////////////////////////////////////////////////////////////
// STRUCTS

type Service struct {
	service *api.Service
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns service
func NewServiceFromServiceAccountJSON(filename string, debug bool) (*Service, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	saConfig, err := google.JWTConfigFromJSON(bytes, api.CloudPlatformScope)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	client := saConfig.Client(getContext(debug))
	if api, err := api.New(client); err != nil {
		return nil, err
	} else {
		return &Service{api}, nil
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
