package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	getFaresCacheFlushDateURL = "http://www.wsdot.wa.gov/ferries/api/fares/rest/cacheflushdate"
)

// GetCacheFlushDate retrieves the fares cache flush date
func (v *VesselsClient) GetCacheFlushDate() (*time.Time, error) {
	req, err := http.NewRequest(http.MethodGet, getFaresCacheFlushDateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := v.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dateStr string
	if err := json.NewDecoder(resp.Body).Decode(&dateStr); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Parse the date string - assuming standard ISO format or custom format
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try alternative format if RFC3339 fails
		parsedTime, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
	}

	return &parsedTime, nil
}