package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/derhabicht/metoc/models"
)

const SunriseSunsetUrl = "https://api.sunrise-sunset.org/json"

type SunriseSunsetResults struct {
	Sunrise                   string `json:"sunrise"`
	Sunset                    string `json:"sunset"`
	SolarNoon                 string `json:"solar_noon"`
	DayLength                 int    `json:"day_length"`
	CivilTwilightBegin        string `json:"civil_twilight_begin"`
	CivilTwilightEnd          string `json:"civil_twilight_end"`
	NauticalTwilightBegin     string `json:"nautical_twilight_begin"`
	NauticalTwilightEnd       string `json:"nautical_twilight_end"`
	AstronomicalTwilightBegin string `json:"astronomical_twilight_begin"`
	AstronomicalTwilightEnd   string `json:"astronomical_twilight_end"`
}

type SunriseSunsetResponse struct {
	Results SunriseSunsetResults `json:"results"`
	Status  string               `json:"status"`
}

func FetchSunriseSunsetData(location models.Location, date models.Dtg) (SunriseSunsetResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", SunriseSunsetUrl, nil)
	if err != nil {
		return SunriseSunsetResponse{}, errors.WithStack(err)
	}

	q := req.URL.Query()
	q.Add("formatted", "0")
	q.Add("lat", fmt.Sprintf("%f", location.Latitude))
	q.Add("lng", fmt.Sprintf("%f", location.Longitude))
	q.Add("date", date.IsoDate())
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return SunriseSunsetResponse{}, errors.WithStack(err)
	}

	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SunriseSunsetResponse{}, errors.WithStack(err)
	}

	var r SunriseSunsetResponse
	err = json.Unmarshal(raw, &r)
	if err != nil {
		return SunriseSunsetResponse{}, errors.WithStack(err)
	}

	return r, nil
}
