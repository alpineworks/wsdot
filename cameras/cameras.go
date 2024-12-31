package cameras

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"alpineworks.io/wsdot"
)

const (
	getCamerasAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/HighwayCameras/HighwayCamerasREST.svc/GetCamerasAsJson"
	getCameraAsJsonURL  = "http://www.wsdot.wa.gov/Traffic/api/HighwayCameras/HighwayCamerasREST.svc/GetCameraAsJson"

	ParamCameraID = "CameraID"
)

type CamerasClient struct {
	wsdot *wsdot.WSDOTClient
}

func NewCamerasClient(wsdotClient *wsdot.WSDOTClient) (*CamerasClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &CamerasClient{
		wsdot: wsdotClient,
	}, nil
}

type CameraLocation struct {
	Description any     `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    int     `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

type Camera struct {
	CameraID         int            `json:"CameraID"`
	CameraLocation   CameraLocation `json:"CameraLocation"`
	CameraOwner      string         `json:"CameraOwner"`
	Description      any            `json:"Description"`
	DisplayLatitude  float64        `json:"DisplayLatitude"`
	DisplayLongitude float64        `json:"DisplayLongitude"`
	ImageHeight      int            `json:"ImageHeight"`
	ImageURL         string         `json:"ImageURL"`
	ImageWidth       int            `json:"ImageWidth"`
	IsActive         bool           `json:"IsActive"`
	OwnerURL         string         `json:"OwnerURL"`
	Region           string         `json:"Region"`
	SortOrder        int            `json:"SortOrder"`
	Title            string         `json:"Title"`
}

func (c *CamerasClient) GetCameras() ([]Camera, error) {
	req, err := http.NewRequest(http.MethodGet, getCamerasAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamAccessCode, c.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cameras []Camera
	if err := json.NewDecoder(resp.Body).Decode(&cameras); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return cameras, nil
}

func (c *CamerasClient) GetCamera(cameraID int) (*Camera, error) {
	req, err := http.NewRequest(http.MethodGet, getCameraAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamAccessCode, c.wsdot.ApiKey)
	q.Add(ParamCameraID, strconv.Itoa(cameraID))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var camera *Camera
	if err := json.NewDecoder(resp.Body).Decode(&camera); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return camera, nil
}
