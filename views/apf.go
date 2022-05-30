package views

import (
	"fmt"
	"sort"

	"github.com/derhabicht/metoc/models"
)

func EncodeRemarks(hf models.HourlyForecast) string {
	rmk := "RMK "

	heatCat := models.CalculateHeatCategory(hf.FeelsLike)
	if heatCat != models.NoRisk {
		rmk += fmt.Sprintf("HEATCAT %s ", heatCat)
	}

	coldCat := models.CalculateColdCategory(hf.FeelsLike)
	if coldCat != models.NoRisk {
		rmk += fmt.Sprintf("COLDCAT %s ", coldCat)
	}

	if hf.PrecipitationProbability >= 10 {
		rmk += fmt.Sprintf("PRECIP PROB %.0f%%", hf.PrecipitationProbability)
	}

	if rmk != "RMK " {
		return rmk
	}

	return ""
}

func EncodePressure(pressure float64) string {
	inhg := pressure / 33.864

	return fmt.Sprintf("A%04.0f", inhg)
}

func EncodeTemperature(temp float64, dewpoint float64) string {
	tempStr := ""
	dewStr := ""

	if temp < 0 {
		temp *= -1
		tempStr = fmt.Sprintf("M%02.0f", temp)
	} else {
		tempStr = fmt.Sprintf("%02.0f", temp)
	}

	if dewpoint < 0 {
		dewpoint *= -1
		dewStr = fmt.Sprintf("M%02.0f", dewpoint)
	} else {
		dewStr = fmt.Sprintf("%02.0f", dewpoint)
	}

	return fmt.Sprintf("%s/%s", tempStr, dewStr)
}

func EncodeWeather(precipType []models.PrecipitationType, precipAmount float64) string {
	if len(precipType) < 1 {
		return ""
	}

	intensity := models.CalculatePrecipitationIntensity(precipAmount)

	prefix := ""
	switch intensity {
	case models.HeavyPrecip:
		prefix = "+"
	case models.LightPrecip:
		prefix = "-"
	case models.NoPrecip:
		return ""
	}

	sort.Slice(precipType, func(i, j int) bool { return precipType[i] < precipType[j] })

	wx := ""
	switch precipType[0] {
	case models.FreezingRain:
		wx = "FZRA"
	case models.Hail:
		wx = "GR"
	case models.Snow:
		wx = "SN"
	case models.Rain:
		wx = "RN"
	}

	return fmt.Sprintf("%s%s", prefix, wx)
}

func EncodeVisibility(visibility float64) string {
	return fmt.Sprintf("%.0fSM", visibility/1.609)
}

func EncodeWind(direction, speed, gust float64) string {
	speed_knots := speed / 1.852
	gust_knots := gust / 1.852

	if (gust_knots - speed_knots) < 5 {
		return fmt.Sprintf("%03.0f%02.f", direction, speed_knots)
	}

	return fmt.Sprintf("%03.0f%02.0fG%02.0f", direction, speed_knots, gust_knots)
}

func EncodeHourlyForecast(hf models.HourlyForecast) string {
	fcst := fmt.Sprintf("FM%s", hf.DateTime.Short())
	wind := EncodeWind(hf.WindDirection, hf.WindSpeed, hf.WindGust)
	visibility := EncodeVisibility(hf.Visibility)
	wx := EncodeWeather(hf.PrecipitationType, hf.PrecipitationAmount)
	sky := string(hf.CloudCover)
	temp := EncodeTemperature(hf.Temperature, hf.Dewpoint)
	pressure := EncodePressure(hf.Pressure)
	remarks := EncodeRemarks(hf)

	if wx != "" {
		return fmt.Sprintf("%s %s %s %s %s %s %s %s\n",
			fcst,
			wind,
			visibility,
			wx,
			sky,
			temp,
			pressure,
			remarks,
		)
	}

	return fmt.Sprintf("%s %s %s %s %s %s %s\n",
		fcst,
		wind,
		visibility,
		sky,
		temp,
		pressure,
		remarks,
	)
}

func EncodeDailyForecast(df models.DailyForecast, location models.Location, generated models.Dtg) string {
	fcst := fmt.Sprintf("APF %s %s\n", location.Mgrs(), generated.Short())

	for i := 0; i < 24; i++ {
		fcst += EncodeHourlyForecast(df.Hours[i])
	}

	return fcst
}
