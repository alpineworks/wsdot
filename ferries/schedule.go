package ferries

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

type inSchedule struct {
	ScheduleID     int64             `json:"ScheduleID"`
	ScheduleName   string            `json:"ScheduleName"`
	ScheduleSeason int               `json:"ScheduleSeason"`
	SchedulePDFUrl string            `json:"SchedulePDFUrl"`
	ScheduleStart  *string           `json:"ScheduleStart"`
	ScheduleEnd    *string           `json:"ScheduleEnd"`
	AllRoutes      []int64           `json:"AllRoutes"`
	TerminalCombos []inTerminalCombo `json:"TerminalCombos"`
}

type inTerminalCombo struct {
	DepartingTerminalID   int64    `json:"DepartingTerminalID"`
	DepartingTerminalName string   `json:"DepartingTerminalName"`
	ArrivingTerminalID    int64    `json:"ArrivingTerminalID"`
	ArrivingTerminalName  string   `json:"ArrivingTerminalName"`
	SailingNotes          string   `json:"SailingNotes"`
	Annotations           []string `json:"Annotations"`
	Times                 []inTime `json:"Times"`
	AnnotationsIVR        []string `json:"AnnotationsIVR"`
}

type inTime struct {
	DepartingTime            *string `json:"DepartingTime"`
	ArrivingTime             *string `json:"ArrivingTime"`
	LoadingRule              int     `json:"LoadingRule"`
	VesselID                 int64   `json:"VesselID"`
	VesselName               string  `json:"VesselName"`
	VesselHandicapAccessible bool    `json:"VesselHandicapAccessible"`
	VesselPositionNum        int     `json:"VesselPositionNum"`
	Routes                   []int64 `json:"Routes"`
	AnnotationIndexes        []int   `json:"AnnotationIndexes"`
}

type Schedule struct {
	ScheduleID     int64           `json:"ScheduleID"`
	ScheduleName   string          `json:"ScheduleName"`
	ScheduleSeason int             `json:"ScheduleSeason"`
	SchedulePDFUrl string          `json:"SchedulePDFUrl"`
	ScheduleStart  *time.Time      `json:"ScheduleStart"`
	ScheduleEnd    *time.Time      `json:"ScheduleEnd"`
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
	DepartingTime            *time.Time `json:"DepartingTime"`
	ArrivingTime             *time.Time `json:"ArrivingTime"`
	LoadingRule              int        `json:"LoadingRule"`
	VesselID                 int64      `json:"VesselID"`
	VesselName               string     `json:"VesselName"`
	VesselHandicapAccessible bool       `json:"VesselHandicapAccessible"`
	VesselPositionNum        int        `json:"VesselPositionNum"`
	Routes                   []int64    `json:"Routes"`
	AnnotationIndexes        []int      `json:"AnnotationIndexes"`
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

	var inSchedule inSchedule
	if err := json.NewDecoder(resp.Body).Decode(&inSchedule); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	schedule := inScheduleToSchedule(inSchedule)

	return &schedule, nil
}

func inScheduleToSchedule(inSchedule inSchedule) Schedule {
	var (
		scheduleStart *time.Time = nil
		scheduleEnd   *time.Time = nil
		err           error
	)
	if inSchedule.ScheduleStart != nil {
		scheduleStart, err = wsdotTimeStringToTime(*inSchedule.ScheduleStart)
		if err != nil {
			slog.Warn("error parsing departing time", "error", err)
			scheduleStart = nil
		}
	}
	if inSchedule.ScheduleEnd != nil {
		scheduleEnd, err = wsdotTimeStringToTime(*inSchedule.ScheduleEnd)
		if err != nil {
			slog.Warn("error parsing arriving time", "error", err)
			scheduleEnd = nil
		}
	}

	schedule := Schedule{
		ScheduleID:     inSchedule.ScheduleID,
		ScheduleName:   inSchedule.ScheduleName,
		ScheduleSeason: inSchedule.ScheduleSeason,
		SchedulePDFUrl: inSchedule.SchedulePDFUrl,
		ScheduleStart:  scheduleStart,
		ScheduleEnd:    scheduleEnd,
		AllRoutes:      inSchedule.AllRoutes,
	}

	for _, combo := range inSchedule.TerminalCombos {
		schedule.TerminalCombos = append(schedule.TerminalCombos, inTerminalComboToTerminalCombo(combo))
	}

	return schedule
}

func inTerminalComboToTerminalCombo(inCombo inTerminalCombo) TerminalCombo {
	combo := TerminalCombo{
		DepartingTerminalID:   inCombo.DepartingTerminalID,
		DepartingTerminalName: inCombo.DepartingTerminalName,
		ArrivingTerminalID:    inCombo.ArrivingTerminalID,
		ArrivingTerminalName:  inCombo.ArrivingTerminalName,
		SailingNotes:          inCombo.SailingNotes,
		Annotations:           inCombo.Annotations,
	}

	for _, time := range inCombo.Times {
		combo.Times = append(combo.Times, inTimeToTime(time))
	}

	return combo
}

func inTimeToTime(inTime inTime) Time {
	var (
		departingTime *time.Time = nil
		arrivingTime  *time.Time = nil
		err           error
	)
	if inTime.DepartingTime != nil {
		departingTime, err = wsdotTimeStringToTime(*inTime.DepartingTime)
		if err != nil {
			slog.Warn("error parsing departing time", "error", err)
			departingTime = nil
		}
	}
	if inTime.ArrivingTime != nil {
		arrivingTime, err = wsdotTimeStringToTime(*inTime.ArrivingTime)
		if err != nil {
			slog.Warn("error parsing arriving time", "error", err)
			arrivingTime = nil
		}
	}

	return Time{
		DepartingTime:            departingTime,
		ArrivingTime:             arrivingTime,
		LoadingRule:              inTime.LoadingRule,
		VesselID:                 inTime.VesselID,
		VesselName:               inTime.VesselName,
		VesselHandicapAccessible: inTime.VesselHandicapAccessible,
		VesselPositionNum:        inTime.VesselPositionNum,
		Routes:                   inTime.Routes,
		AnnotationIndexes:        inTime.AnnotationIndexes,
	}
}

func wsdotTimeStringToTime(wsdotTime string) (*time.Time, error) {
	// /Date(1742713200000-0700)/
	re := regexp.MustCompile(`^/Date\((\d+)([+-]\d{4})\)/$`)

	matches := re.FindStringSubmatch(wsdotTime)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid WSDOT time string format: %s", wsdotTime)
	}

	// Parse the timestamp
	milliseconds, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing milliseconds: %v", err)
	}

	// WARNING: The offset is not currently used in the conversion, as it's already accounted for

	// // Parse the offset
	// offset, err := parseOffsetMilli(matches[2])
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing offset: %v", err)
	// }

	// // Adjust the milliseconds with the offset
	// milliseconds += int64(offset)

	// Convert milliseconds to time.Time
	t := time.UnixMilli(milliseconds)

	return &t, nil
}

func parseOffsetMilli(offset string) (int, error) {
	if len(offset) != 5 {
		return 0, fmt.Errorf("failed to parse offset - length incorrect")
	}

	sign := string(offset[0])

	offsetHours := string(offset[1:3])
	offsetMinutes := string(offset[3:5])

	hours, err := strconv.ParseInt(offsetHours, 10, 64)
	if err != nil || (hours < 0 || hours > 23) {
		return 0, fmt.Errorf("error parsing hours: %v", err)
	}

	minutes, err := strconv.ParseInt(offsetMinutes, 10, 64)
	if err != nil || (minutes < 0 || minutes > 59) {
		return 0, fmt.Errorf("error parsing minutes: %v", err)
	}

	totalMillis := (hours*60 + minutes) * 60 * 1000

	switch sign {
	case "+":
		return int(totalMillis), nil
	case "-":
		return -int(totalMillis), nil
	default:
		return 0, fmt.Errorf("invalid sign: %s", string(sign))
	}
}
