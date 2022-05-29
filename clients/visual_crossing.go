package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/derhabicht/metoc/models"
)

const VisualCrossingUrl = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

func FetchVisualCrossingData(apiKey string, location models.Location, startDate, endDate models.Dtg) ([]byte, error) {
	client := &http.Client{}

	path := fmt.Sprintf(
		"%s/%s/%s/%s",
		VisualCrossingUrl,
		url.PathEscape(fmt.Sprintf("%f,%f", location.Latitude, location.Longitude)),
		url.PathEscape(startDate.IsoDate()),
		url.PathEscape(endDate.Date()),
	)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("unitGroup", "metric")
	q.Add("include", url.QueryEscape("days,hours"))
	q.Add("contentType", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
