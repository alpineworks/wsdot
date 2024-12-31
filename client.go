package wsdot

import (
	"errors"
	"net/http"
)

type WSDOTClient struct {
	Client *http.Client
	ApiKey string
}

type WSDOTClientOption func(*WSDOTClient)

type WSDOTClientError error

var (
	ErrInvalidAPIKey WSDOTClientError = errors.New("invalid api key")
	ErrNoClient      WSDOTClientError = errors.New("no client")
)

const (
	ParamAccessCode = "AccessCode"
)

func NewWSDOTClient(options ...WSDOTClientOption) (*WSDOTClient, error) {
	client := &http.Client{}
	wsdotClient := &WSDOTClient{
		Client: client,
	}

	for _, option := range options {
		option(wsdotClient)
	}

	if wsdotClient.ApiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	return wsdotClient, nil
}

func WithHTTPClient(client *http.Client) WSDOTClientOption {
	return func(w *WSDOTClient) {
		w.Client = client
	}
}

func WithAPIKey(apiKey string) WSDOTClientOption {
	return func(w *WSDOTClient) {
		w.ApiKey = apiKey
	}
}
