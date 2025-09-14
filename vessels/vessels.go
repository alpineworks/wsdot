package vessels

import (
	"alpineworks.io/wsdot"
)

// VesselsClient provides access to all WSDOT Vessels, Terminals, and Fares API endpoints
type VesselsClient struct {
	wsdot *wsdot.WSDOTClient
}

// NewVesselsClient creates a new unified Vessels client
func NewVesselsClient(wsdotClient *wsdot.WSDOTClient) (*VesselsClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &VesselsClient{
		wsdot: wsdotClient,
	}, nil
}

// GetWSDOTClient returns the underlying WSDOT client
func (v *VesselsClient) GetWSDOTClient() *wsdot.WSDOTClient {
	return v.wsdot
}