package models

type LunarPhase string

const (
	New            LunarPhase = "NEW"
	WaxingCrescent LunarPhase = "WXC"
	FirstQuarter   LunarPhase = "1ST"
	WaxingGibbous  LunarPhase = "WXG"
	Full           LunarPhase = "FUL"
	WaningGibbous  LunarPhase = "WNG"
	LastQuarter    LunarPhase = "LST"
	WaningCrescent LunarPhase = "WNC"
)

type Twilight struct {
	Begin Dtg
	End   Dtg
}

type SunData struct {
	Sunrise              Dtg
	Sunset               Dtg
	AstronomicalTwilight Twilight
	NauticalTwilight     Twilight
	CivilTwilight        Twilight
}

type MoonData struct {
	Phase LunarPhase
	Rise  Dtg
	Set   Dtg
}

type AstroData struct {
	SunData
	MoonData
}
