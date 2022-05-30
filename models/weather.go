package models

type PrecipitationType int

const (
	FreezingRain PrecipitationType = 0
	Hail         PrecipitationType = 1
	Snow         PrecipitationType = 2
	Rain         PrecipitationType = 3
)

type CloudCover string

const (
	Overcast  CloudCover = "OVC" // 8/8 cloud cover
	Broken    CloudCover = "BKN" // 5/8-7/8 cloud cover
	Scattered CloudCover = "SCT" // 3/8-5/8 cloud cover
	Few       CloudCover = "FEW" // 1/8-3/8 cloud cover
	Clear     CloudCover = "CLR" // <1/8 cloud cover
)

type TempRiskCategory string

const (
	ExtremeRisk  TempRiskCategory = "EXTREME"
	HighRisk     TempRiskCategory = "HIGH"
	ModerateRisk TempRiskCategory = "MODERATE"
	LowRisk      TempRiskCategory = "LOW"
	NoRisk       TempRiskCategory = "N/A"
)

func CalculateHeatCategory(heatIndex float64) TempRiskCategory {
	if heatIndex > 46.1 {
		return ExtremeRisk
	} else if heatIndex > 39.4 {
		return HighRisk
	} else if heatIndex > 32.8 {
		return ModerateRisk
	} else if heatIndex > 29.4 {
		return LowRisk
	}

	return NoRisk
}

func CalculateColdCategory(windChill float64) TempRiskCategory {
	if windChill < -29.4 {
		return ExtremeRisk
	} else if windChill < -18.3 {
		return HighRisk
	} else if windChill < -6.7 {
		return ModerateRisk
	} else if windChill < 4.4 {
		return LowRisk
	}

	return NoRisk
}

type PrecipitationIntensity string

const (
	HeavyPrecip    PrecipitationIntensity = "HEAVY"
	ModeratePrecip PrecipitationIntensity = "MODERATE"
	LightPrecip    PrecipitationIntensity = "LIGHT"
	NoPrecip       PrecipitationIntensity = "NONE"
)

func CalculatePrecipitationIntensity(precip float64) PrecipitationIntensity {
	if precip > 7.6 {
		return HeavyPrecip
	} else if precip > 2.5 {
		return ModeratePrecip
	} else if precip > 0 {
		return LightPrecip
	}

	return NoPrecip
}

type HourlyForecast struct {
	DateTime                 Dtg
	WindDirection            float64
	WindSpeed                float64
	WindGust                 float64
	Visibility               float64
	PrecipitationAmount      float64
	PrecipitationType        []PrecipitationType
	PrecipitationProbability float64
	CloudCover               CloudCover
	Temperature              float64
	FeelsLike                float64
	Dewpoint                 float64
	Pressure                 float64
}

type DailyForecast struct {
	Conditions   string
	Description  string
	HighTemp     float64
	LowTemp      float64
	FeelsLikeMax float64
	FeelsLikeMin float64
	Hours        map[int]HourlyForecast
}
