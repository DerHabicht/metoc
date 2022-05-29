package models

import "time"

type MetocReport struct {
	Dates     []Dtg
	Generated Dtg
	Location  Location
	TzOffset  int
	AstroData map[Dtg]AstroData
	Forecast  map[Dtg]DailyForecast
}

func NewMetocReport(locationName, locationMgrs string, tzoffset int) (MetocReport, error) {
	l, err := ParseLocationFromMgrs(locationName, locationMgrs)
	if err != nil {
		return MetocReport{}, err
	}

	return MetocReport{
		Generated: Dtg{time.Now()},
		Location: l,
		TzOffset: tzoffset,
		AstroData: make(map[Dtg]AstroData),
		Forecast: make(map[Dtg]DailyForecast),
	}, nil
}

func (m *MetocReport) AddDateToMetocReport(date Dtg) {
	m.Dates = append(m.Dates, date)
	m.AstroData[date] = AstroData{}
	m.Forecast[date] = DailyForecast{}
}
