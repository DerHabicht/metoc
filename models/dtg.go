package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ParseAcp121Full   = "021504 Jan 2006"
	FormatAcp121Full  = "021504MST Jan 2006"
	ParseAcp121Short  = "021504"
	FormatAcp121Short = "021504MST"
	Acp121Date        = "_2 Jan 2006"
	ParseAcp121Time   = "1504"
	FormatAcp121Time  = "1504MST"
)

var AcpTimeZonesByOffset = map[int]string{
	1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "K", 11: "L", 12: "M",
	-1: "N", -2: "O", -3: "P", -4: "Q", -5: "R", -6: "S", -7: "T", -8: "U", -9: "V", -10: "W", -11: "X", -12: "Y",
	0: "Z",
}

var AcpTimeZonesByLetter = map[string]int{
	"A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6, "G": 7, "H": 8, "I": 9, "K": 10, "L": 11, "M": 12,
	"N": -1, "O": -2, "P": -3, "Q": -4, "R": -5, "S": -6, "T": -7, "U": -8, "V": -9, "W": -10, "X": -11, "Y": -12,
	"Z": 0,
}

func AcpLocation(offset int) *time.Location {
	return time.FixedZone(AcpTimeZonesByOffset[offset], offset*3600)
}

func ParseAcpTimeZone(tz string) (*time.Location, error) {
	if len(tz) != 1 {
		return nil, errors.Errorf("%s is not a valid ACP 121 timezone designator", tz)
	}
	offset, ok := AcpTimeZonesByLetter[strings.ToUpper(tz)]
	if !ok {
		return nil, errors.Errorf("%s is not a valid ACP 121 timezone designator", tz)
	}

	return AcpLocation(offset), nil
}

type Dtg struct {
	datetime time.Time
}

func (d Dtg) Full() string {
	return d.datetime.Format(FormatAcp121Full)
}

func (d Dtg) Short() string {
	return d.datetime.Format(FormatAcp121Short)
}

func (d Dtg) Date() string {
	return d.datetime.Format(Acp121Date)
}

func (d Dtg) Time() string {
	return d.datetime.Format(FormatAcp121Time)
}

func (d Dtg) IsoDate() string {
	return d.datetime.Format("2006-01-02")
}

func (d Dtg) Hour() int {
	return d.datetime.Hour()
}

func (d Dtg) Unix() int64 {
	return d.datetime.Unix()
}

func DtgNow(tzoffset int) Dtg {
	l := AcpLocation(tzoffset)
	d := time.Now().In(l)

	return Dtg{d}
}

func ParseDtg(dtg string) (Dtg, error) {
	r := regexp.MustCompile(`(\d{4,6})([[:alpha:]]).*`)
	match := r.FindStringSubmatch(dtg)

	if len(match) == 4 {
		l, err := ParseAcpTimeZone(match[2])
		if err != nil {
			return Dtg{}, errors.WithStack(err)
		}

		d := match[1] + match[3]

		// Try full
		dt, err := time.ParseInLocation(ParseAcp121Full, d, l)
		if err == nil {
			return Dtg{dt}, nil
		}

		// Try short
		dt, err = time.ParseInLocation(ParseAcp121Short, d, l)
		if err == nil {
			return Dtg{dt}, nil
		}

		// Try time
		// Try short
		dt, err = time.ParseInLocation(ParseAcp121Short, d, l)
		if err != nil {
			return Dtg{}, errors.WithMessagef(err, "failed to parse %s as an ACP DTG", dtg)
		}

		return Dtg{dt}, nil
	}

	dt, err := time.Parse(Acp121Date, dtg)
	if err != nil {
		return Dtg{}, errors.WithMessagef(err, "failed to parse %s as an ACP date", dtg)
	}

	return Dtg{dt}, nil
}

func ParseIsoDateTimeToDtg(iso string, tzoffset int) (Dtg, error) {
	location := AcpLocation(tzoffset)

	dt, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		dt, err = time.Parse("2006-01-02T15:04:05", iso)
		if err != nil {
			return Dtg{}, errors.WithMessagef(err, "failed to parse %s as an ISO date-time", iso)
		}
	}

	localized := dt.In(location)

	return Dtg{localized}, nil
}

func ParseIsoDateToDtg(iso string) (Dtg, error) {
	dt, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return Dtg{}, errors.WithMessagef(err, "failed to parse %s as an ISO date", iso)
	}

	return Dtg{dt}, nil
}

func ParseTimeToDtg(t string, tzoffset int) (Dtg, error) {
	location := AcpLocation(tzoffset)

	var dt time.Time

	dt, err := time.ParseInLocation("15:04:05", t, location)
	if err != nil {
		dt, err = time.ParseInLocation("15:04", t, location)
		if err != nil {
			return Dtg{}, errors.WithMessagef(err, "failed to parse %s as a time", t)
		}
	}

	return Dtg{dt}, nil
}
