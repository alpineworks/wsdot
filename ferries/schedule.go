package ferries

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getRouteSchedulesAsJsonURL         = "https://www.wsdot.wa.gov/Ferries/API/Schedule/rest/schedroutes"
	getScheduleTodayByRouteIDAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Schedule/rest/scheduletoday/%d/%t"
)

type RouteSchedule struct {
	ScheduleID         int                     `json:"ScheduleID"`
	SchedRouteID       int                     `json:"SchedRouteID"`
	ContingencyOnly    bool                    `json:"ContingencyOnly"`
	RouteID            int                     `json:"RouteID"`
	RouteAbbrev        string                  `json:"RouteAbbrev"`
	Description        string                  `json:"Description"`
	SeasonalRouteNotes string                  `json:"SeasonalRouteNotes"`
	RegionID           int                     `json:"RegionID"`
	ServiceDisruptions []ServiceDisruption     `json:"ServiceDisruptions"`
	ContingencyAdj     []ContingencyAdjustment `json:"ContingencyAdj"`
}

type ServiceDisruption struct {
	BulletinID            int    `json:"BulletinID"`
	BulletinFlag          bool   `json:"BulletinFlag"`
	PublishDate           string `json:"PublishDate"`
	DisruptionDescription string `json:"DisruptionDescription"`
}

type ContingencyAdjustment struct {
	DateFrom               string  `json:"DateFrom"`
	DateThru               string  `json:"DateThru"`
	EventID                *int    `json:"EventID"`
	EventDescription       *string `json:"EventDescription"`
	AdjType                int     `json:"AdjType"`
	ReplacedBySchedRouteID *int    `json:"ReplacedBySchedRouteID"`
}

func (f *FerriesClient) GetRouteSchedules() ([]RouteSchedule, error) {
	req, err := http.NewRequest(http.MethodGet, getRouteSchedulesAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, f.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var schedules []RouteSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return schedules, nil
}

type Schedule struct {
	ScheduleID     int64           `json:"ScheduleID"`
	ScheduleName   string          `json:"ScheduleName"`
	ScheduleSeason int             `json:"ScheduleSeason"`
	SchedulePDFUrl string          `json:"SchedulePDFUrl"`
	ScheduleStart  string          `json:"ScheduleStart"`
	ScheduleEnd    string          `json:"ScheduleEnd"`
	AllRoutes      []int64         `json:"AllRoutes"`
	TerminalCombos []TerminalCombo `json:"TerminalCombos"`
}

type TerminalCombo struct {
	DepartingTerminalID   int64    `json:"DepartingTerminalID"`
	DepartingTerminalName string   `json:"DepartingTerminalName"`
	ArrivingTerminalID    int64    `json:"ArrivingTerminalID"`
	ArrivingTerminalName  string   `json:"ArrivingTerminalName"`
	SailingNotes          string   `json:"SailingNotes"`
	Annotations           []string `json:"Annotations"`
	Times                 []Time   `json:"Times"`
	AnnotationsIVR        []string `json:"AnnotationsIVR"`
}

type Time struct {
	DepartingTime            string  `json:"DepartingTime"`
	ArrivingTime             string  `json:"ArrivingTime"`
	LoadingRule              int     `json:"LoadingRule"`
	VesselID                 int64   `json:"VesselID"`
	VesselName               string  `json:"VesselName"`
	VesselHandicapAccessible bool    `json:"VesselHandicapAccessible"`
	VesselPositionNum        int     `json:"VesselPositionNum"`
	Routes                   []int64 `json:"Routes"`
	AnnotationIndexes        []int   `json:"AnnotationIndexes"`
}

func (f *FerriesClient) GetSchedulesTodayByRouteID(routeID int, onlyRemainingTimes bool) (*Schedule, error) {
	url := fmt.Sprintf(getScheduleTodayByRouteIDAsJsonURL, routeID, onlyRemainingTimes)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, f.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var schedules Schedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &schedules, nil
}
