package vessels

import (
	"alpineworks.io/wsdot"
)

type VesselsClient struct {
	wsdot *wsdot.WSDOTClient
}

func NewVesselsClient(wsdotClient *wsdot.WSDOTClient) (*VesselsClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &VesselsClient{
		wsdot: wsdotClient,
	}, nil
}
