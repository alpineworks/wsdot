package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"alpineworks.io/wsdot"
)

const (
	getTerminalWaitTimesURL = "http://www.wsdot.wa.gov/ferries/api/terminals/rest/terminalwaittimes"
)

// TerminalWaitTime represents wait time information for a terminal
type TerminalWaitTime struct {
	TerminalID      int    `json:"TerminalID"`
	TerminalName    string `json:"TerminalName"`
	TerminalAbbrev  string `json:"TerminalAbbrev"`
	WaitTimeNotes   string `json:"WaitTimeNotes"`
	WaitTime        int    `json:"WaitTime"`
	LastUpdated     string `json:"LastUpdated"`
	DisplayOrder    int    `json:"DisplayOrder"`
}

// GetTerminalWaitTimes retrieves wait times for all terminals
func (v *VesselsClient) GetTerminalWaitTimes() ([]TerminalWaitTime, error) {
	req, err := http.NewRequest(http.MethodGet, getTerminalWaitTimesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, v.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var waitTimes []TerminalWaitTime
	if err := json.NewDecoder(resp.Body).Decode(&waitTimes); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return waitTimes, nil
}

// GetTerminalWaitTimeByID retrieves wait time for a specific terminal
func (v *VesselsClient) GetTerminalWaitTimeByID(terminalID int) (*TerminalWaitTime, error) {
	url := fmt.Sprintf("%s/%s", getTerminalWaitTimesURL, strconv.Itoa(terminalID))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, v.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var waitTime TerminalWaitTime
	if err := json.NewDecoder(resp.Body).Decode(&waitTime); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &waitTime, nil
}