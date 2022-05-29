package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/derhabicht/metoc/models"
)

const UsnoUrl = "https://aa.usno.navy.mil/api/rstt/oneday"

type UsnoGeometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type string `json:"type"`
}

type UsnoClosestPhase struct {
	Day int `json:"day"`
	Month int `json:"month"`
	Phase string `json:"phase"`
	Time string `json:"time"`
	Year int `json:"year"`
}

type UsnoPhenomenon struct {
	Phenomenon string `json:"phen"`
	Time string `json:"time"`
}

type UsnoData struct {
	ClosestPhase UsnoClosestPhase `json:"closestphase"`
	CurrentPhase string `json:"curphase"`
	Day int `json:"day"`
	DayOfWeek string `json:"day_of_week"`
	Fracillum string `json:"fracillum"`
	IsDst bool `json:"isdst"`
	Month int `json:"month"`
	MoonData []UsnoPhenomenon `json:"moondata"`
	SunData []UsnoPhenomenon `json:"sundata"`
	TimezoneOffset float64 `json:"tz"`
	Year int `json:"year"`
}

type UsnoProperties struct {
	Data UsnoData `json:"data"`
}

type UsnoResponse struct {
	ApiVersion string `json:"apiversion"`
	Geometry UsnoGeometry `json:"geometry"`
	Properties UsnoProperties `json:"properties"`
	Type string `json:"type"`
}

func FetchUsnoData(location models.Location, date models.Dtg, tzoffset int) (UsnoResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", UsnoUrl, nil)
	if err != nil {
		return UsnoResponse{}, errors.WithStack(err)
	}

	q := req.URL.Query()
	q.Add("id", "ag7if")
	q.Add("coords", fmt.Sprintf("%f,%f", location.Latitude, location.Longitude))
	q.Add("date", date.IsoDate())
	q.Add("tz", fmt.Sprintf("%d", tzoffset))
	q.Add("dst", "false")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return UsnoResponse{}, errors.WithStack(err)
	}

	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UsnoResponse{}, errors.WithStack(err)
	}

	var r UsnoResponse
	err = json.Unmarshal(raw, &r)
	if err != nil {
		return UsnoResponse{}, errors.WithStack(err)
	}

	return r, nil
}
