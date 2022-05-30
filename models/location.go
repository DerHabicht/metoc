package models

import (
	"regexp"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	"github.com/tzneal/coordconv"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
	Precision int
}

func (l *Location) Mgrs() string {
	c := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
	m, _ := coordconv.DefaultMGRSConverter.ConvertFromGeodetic(c, l.Precision)

	return m
}

func MgrsPrecision(mgrs string) (int, error) {
	r := regexp.MustCompile(`\d{1,2}[[:alpha:]]{3}(\d+)`)
	match := r.FindStringSubmatch(mgrs)

	if len(match) < 2 {
		return -1, errors.Errorf("attempt to determine precision of invalid MGRS string: %s", mgrs)
	}

	return len(match[1]) / 2, nil
}

func ParseLocationFromMgrs(name, mgrs string) (Location, error) {
	l, err := coordconv.DefaultMGRSConverter.ConvertToGeodetic(mgrs)
	if err != nil {
		return Location{}, err
	}

	p, err := MgrsPrecision(mgrs)

	return Location{
		Name:      name,
		Latitude:  l.Lat.Degrees(),
		Longitude: l.Lng.Degrees(),
		Precision: p,
	}, nil
}
