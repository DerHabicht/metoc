package models

type MetocReport struct {
	dates     []Dtg
	Generated Dtg
	Location  Location
	AstroData map[Dtg]AstroData
	Forecast  map[Dtg]DailyForecast
}

func NewMetocReport(locationName, locationMgrs string) (MetocReport, error) {
	l, err := ParseLocationFromMgrs(locationName, locationMgrs)
	if err != nil {
		return MetocReport{}, err
	}

	return MetocReport{Location: l}, nil
}

func (m *MetocReport) AddDateToMetocReport(date Dtg) error {
	return nil
}
