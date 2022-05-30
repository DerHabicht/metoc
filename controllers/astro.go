package controllers

import (
	"github.com/pkg/errors"

	"github.com/derhabicht/metoc/clients"
	"github.com/derhabicht/metoc/models"
)

func FetchSunData(location models.Location, date models.Dtg, tzoffset int) (models.SunData, error) {
	resp, err := clients.FetchSunriseSunsetData(location, date)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	sunrise, err := models.ParseIsoDateTimeToDtg(resp.Results.Sunrise, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	sunset, err := models.ParseIsoDateTimeToDtg(resp.Results.Sunset, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	beginAt, err := models.ParseIsoDateTimeToDtg(resp.Results.AstronomicalTwilightBegin, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	beginNt, err := models.ParseIsoDateTimeToDtg(resp.Results.NauticalTwilightBegin, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	beginCt, err := models.ParseIsoDateTimeToDtg(resp.Results.CivilTwilightBegin, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	endAt, err := models.ParseIsoDateTimeToDtg(resp.Results.AstronomicalTwilightEnd, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	endNt, err := models.ParseIsoDateTimeToDtg(resp.Results.NauticalTwilightEnd, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	endCt, err := models.ParseIsoDateTimeToDtg(resp.Results.CivilTwilightEnd, tzoffset)
	if err != nil {
		return models.SunData{}, errors.WithStack(err)
	}

	return models.SunData{
		Sunrise: sunrise,
		Sunset:  sunset,
		AstronomicalTwilight: models.Twilight{
			Begin: beginAt,
			End:   endAt,
		},
		NauticalTwilight: models.Twilight{
			Begin: beginNt,
			End:   endNt,
		},
		CivilTwilight: models.Twilight{
			Begin: beginCt,
			End:   endCt,
		},
	}, nil
}

func FetchMoonData(location models.Location, date models.Dtg, tzoffset int) (models.MoonData, error) {
	resp, err := clients.FetchUsnoData(location, date, tzoffset)
	if err != nil {
		return models.MoonData{}, errors.WithStack(err)
	}

	var rise models.Dtg
	var set models.Dtg

	moondata := resp.Properties.Data.MoonData
	for _, d := range moondata {
		if d.Phenomenon == "Rise" {
			rise, err = models.ParseTimeToDtg(d.Time, tzoffset)
			if err != nil {
				return models.MoonData{}, errors.WithStack(err)
			}
		} else if d.Phenomenon == "Set" {
			set, err = models.ParseTimeToDtg(d.Time, tzoffset)
			if err != nil {
				return models.MoonData{}, errors.WithStack(err)
			}
		}
	}

	var phase models.LunarPhase
	switch resp.Properties.Data.CurrentPhase {
	case "New Moon":
		phase = models.New
	case "Waxing Crescent":
		phase = models.WaxingCrescent
	case "First Quarter":
		phase = models.FirstQuarter
	case "Waxing Gibbous":
		phase = models.WaxingGibbous
	case "Full Moon":
		phase = models.Full
	case "Waning Gibbous":
		phase = models.WaningGibbous
	case "Third Quarter":
		phase = models.LastQuarter
	case "Waning Crescent":
		phase = models.WaningCrescent
	default:
		return models.MoonData{}, errors.Errorf(
			"failed to parse %s as a moon phase",
			resp.Properties.Data.CurrentPhase,
		)
	}

	return models.MoonData{
		Phase:    phase,
		MoonRise: rise,
		MoonSet:  set,
	}, nil
}

func FetchDailyAstroData(location models.Location, dtg models.Dtg, tzoffset int) (models.AstroData, error) {
	sunData, err := FetchSunData(location, dtg, tzoffset)
	if err != nil {
		return models.AstroData{}, errors.WithStack(err)
	}

	moonData, err := FetchMoonData(location, dtg, tzoffset)
	if err != nil {
		return models.AstroData{}, errors.WithStack(err)
	}

	return models.AstroData{SunData: sunData, MoonData: moonData}, nil
}

func FetchAstroDataForReport(report *models.MetocReport) error {
	for _, date := range report.Dates {
		ad, err := FetchDailyAstroData(report.Location, date, report.TzOffset)
		if err != nil {
			return errors.WithStack(err)
		}

		report.AstroData[date] = ad
	}

	return nil
}
