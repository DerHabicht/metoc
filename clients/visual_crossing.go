package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/derhabicht/metoc/models"
)

const VisualCrossingUrl = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

type VisualCrossingHourlyData struct {
	Datetime                 string   `json:"datetime"`
	DatetimeEpoch            int      `json:"DatetimeEpoch"`
	Temperature              float64  `json:"temp"`
	FeelsLike                float64  `json:"feelslike"`
	Humidity                 float64  `json:"humidity"`
	Dewpoint                 float64  `json:"dewpoint"`
	Precipitation            float64  `json:"precip"`
	PrecipitationProbability float64  `json:"precipprob"`
	Snow                     float64  `json:"snow"`
	SnowDepth                float64  `json:"snowdepth"`
	PrecipitationType        []string `json:"preciptype"`
	WindGust                 float64  `json:"windgust"`
	WindSpeed                float64  `json:"windspeed"`
	WindDirection            float64  `json:"winddir"`
	Pressure                 float64  `json:"pressure"`
	Visibility               float64  `json:"visibility"`
	CloudCover               float64  `json:"cloudcover"`
	SolarRadiation           float64  `json:"solarradiation"`
	SolarEnergy              float64  `json:"solarenergy"`
	UvIndex                  float64  `json:"uvindex"`
	SevereRisk               float64  `json:"severerisk"`
	Conditions               string   `json:"conditions"`
	Icon                     string   `json:"icon"`
	Source                   string   `json:"source"`
}

type VisualCrossingDailyData struct {
	Datetime          string                     `json:"datetime"`
	DatetimeEpoch     int                        `json:"datetimeEpoch"`
	TempMax           float64                    `json:"tempmax"`
	TempMin           float64                    `json:"tempmin"`
	Temp              float64                    `json:"temp"`
	FeelsLikeMax      float64                    `json:"feelslikemax"`
	FeelsLikeMin      float64                    `json:"feelslikemin"`
	FeelsLike         float64                    `json:"feelslike"`
	DewPoint          float64                    `json:"dew"`
	Humidity          float64                    `json:"humidity"`
	Precipitation     float64                    `json:"precip"`
	PrecipProbability float64                    `json:"precipprob"`
	PrecipCover       float64                    `json:"precipcover"`
	PrecipType        []string                   `json:"preciptype"`
	Snow              float64                    `json:"snow"`
	SnowDepth         float64                    `json:"snowdepth"`
	WindGust          float64                    `json:"windgust"`
	WindSpeed         float64                    `json:"windspeed"`
	WindDirection     float64                    `json:"winddir"`
	Pressure          float64                    `json:"pressure"`
	CloudCover        float64                    `json:"cloudcover"`
	Visibility        float64                    `json:"visibility"`
	SolarRadiation    float64                    `json:"solarradiation"`
	SolarEnergy       float64                    `json:"solarenergy"`
	UvIndex           float64                    `json:"uvindex"`
	SevereRisk        float64                    `json:"severerisk"`
	Sunrise           string                     `json:"sunrise"`
	SunriseEpoch      int                        `json:"sunriseEpoch"`
	MoonPhase         float64                    `json:"moonphase"`
	Conditions        string                     `json:"conditions"`
	Description       string                     `json:"description"`
	Icon              string                     `json:"icon"`
	Source            string                     `json:"source"`
	Hours             []VisualCrossingHourlyData `json:"hours"`
}

type VisualCrossingResponse struct {
	QueryCost       int                       `json:"queryCost"`
	Latitude        float64                   `json:"latitude"`
	Longitude       float64                   `json:"Longitude"`
	ResolvedAddress string                    `json:"resolvedAddress"`
	Address         string                    `json:"address"`
	Timezone        string                    `json:"timezone"`
	TzOffset        float64                   `json:"tzoffset"`
	Days            []VisualCrossingDailyData `json:"days"`
}

func FetchVisualCrossingData(apiKey string, location models.Location, startDate, endDate models.Dtg) (VisualCrossingResponse, error) {
	client := &http.Client{}

	path := fmt.Sprintf(
		"%s/%s/%s/%s",
		VisualCrossingUrl,
		url.PathEscape(fmt.Sprintf("%f,%f", location.Latitude, location.Longitude)),
		url.PathEscape(startDate.IsoDate()),
		url.PathEscape(endDate.IsoDate()),
	)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return VisualCrossingResponse{}, errors.WithStack(err)
	}

	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("unitGroup", "metric")
	q.Add("include", "days,hours")
	q.Add("contentType", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return VisualCrossingResponse{}, errors.WithStack(err)
	}

	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return VisualCrossingResponse{}, errors.WithStack(err)
	}

	var vcr VisualCrossingResponse
	err = json.Unmarshal(raw, &vcr)
	if err != nil {
		return VisualCrossingResponse{}, errors.WithStack(err)
	}

	return vcr, nil
}
