package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"alpineworks.io/wsdot"
)

const (
	getCacheFlushDateAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/cacheflushdate"
)

type CacheFlushDate struct {
	CacheFlushDate string `json:"CacheFlushDate"`
}

func (v *VesselsClient) GetCacheFlushDate() (*time.Time, error) {
	req, err := http.NewRequest(http.MethodGet, getCacheFlushDateAsJsonURL, nil)
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

	var cacheFlushDate CacheFlushDate
	if err := json.NewDecoder(resp.Body).Decode(&cacheFlushDate); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Parse the WSDOT time string format
	parsedTime, err := wsdotTimeStringToTime(cacheFlushDate.CacheFlushDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing cache flush date: %v", err)
	}

	return parsedTime, nil
}