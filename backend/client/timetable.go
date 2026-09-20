package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (c *LibrusClient) GetTimetable(
	ctx context.Context,
	weekStart time.Time,
) (*LibrusTimetableResponse, error) {
	u, err := url.Parse(c.APIBaseURL + "/Timetables")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("weekStart", weekStart.Format("2006-01-02"))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("timetable: librus API responded with %s", resp.Status)
	}

	var data LibrusTimetableResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
