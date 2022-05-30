package controllers

import (
	"fmt"
	"github.com/derhabicht/metoc/clients"
	"github.com/derhabicht/metoc/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ParsePrecipType(raw []string) []models.PrecipitationType {
	var types []models.PrecipitationType

	for _, t := range raw {
		switch t {
		case "rain":
			types = append(types, models.Rain)
		case "snow":
			types = append(types, models.Snow)
		case "freezingrain":
			types = append(types, models.FreezingRain)
		case "ice":
			types = append(types, models.Hail)
		}
	}

	return types
}

func ParseCloudCover(raw float64) models.CloudCover {
	if raw < 12.5 {
		return models.Clear
	} else if raw < 37.5 {
		return models.Scattered
	} else if raw < 62.5 {
		return models.Few
	} else if raw < 100.0 {
		return models.Broken
	}

	return models.Overcast
}

func ParseHourlyForecast(dtg models.Dtg, raw clients.VisualCrossingHourlyData) models.HourlyForecast {
	return models.HourlyForecast{
		DateTime:                 dtg,
		WindDirection:            raw.WindDirection,
		WindSpeed:                raw.WindSpeed,
		WindGust:                 raw.WindGust,
		Visibility:               raw.Visibility,
		PrecipitationAmount:      raw.Precipitation,
		PrecipitationType:        ParsePrecipType(raw.PrecipitationType),
		PrecipitationProbability: raw.PrecipitationProbability,
		CloudCover:               ParseCloudCover(raw.CloudCover),
		Temperature:              raw.Temperature,
		FeelsLike:                raw.FeelsLike,
		Dewpoint:                 raw.Dewpoint,
		Pressure:                 raw.Pressure,
	}
}

func ParseDailyForecast(raw clients.VisualCrossingDailyData, day models.Dtg, tzoffset int) (models.DailyForecast, error) {
	df := models.DailyForecast{
		Conditions:   raw.Conditions,
		Description:  raw.Description,
		HighTemp:     raw.TempMax,
		LowTemp:      raw.TempMin,
		FeelsLikeMax: raw.FeelsLikeMax,
		FeelsLikeMin: raw.FeelsLikeMin,
		Hours:        make(map[int]models.HourlyForecast),
	}

	for _, hour := range raw.Hours {
		dt := fmt.Sprintf("%sT%s", day.IsoDate(), hour.Datetime)
		dtg, err := models.ParseIsoDateTimeToDtg(dt, tzoffset)
		if err != nil {
			return models.DailyForecast{}, errors.WithStack(err)
		}

		hf := ParseHourlyForecast(dtg, hour)
		df.Hours[dtg.Hour()] = hf
	}

	return df, nil
}

func FetchWeatherDataForReport(report *models.MetocReport) error {
	data, err := clients.FetchVisualCrossingData(
		viper.GetString("visual-crossing-api-key"),
		report.Location,
		report.Dates[0],
		report.Dates[len(report.Dates)-1],
	)

	if err != nil {
		return errors.WithStack(err)
	}

	for _, day := range data.Days {
		dtg, err := models.ParseIsoDateToDtg(day.Datetime)

		df, err := ParseDailyForecast(day, dtg, int(data.TzOffset))
		if err != nil {
			return errors.WithStack(err)
		}

		report.Forecast[dtg] = df
	}

	return nil
}
