package ferries

import (
	"alpineworks.io/wsdot"
)

type FerriesClient struct {
	wsdot *wsdot.WSDOTClient
}

func NewFerriesClient(wsdotClient *wsdot.WSDOTClient) (*FerriesClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &FerriesClient{
		wsdot: wsdotClient,
	}, nil
}
