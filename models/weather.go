package models

type HourlyForecast struct {
}

type DailyForecast struct {
	hours map[int]HourlyForecast
}
