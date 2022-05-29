package views

import (
	"fmt"
	"github.com/derhabicht/metoc/models"
	"strings"
)

const Template = `\documentclass[10pt]{article}

\usepackage{booktabs}
\usepackage[margin=0.5in]{geometry}

\renewcommand{\familydefault}{\ttdefault}
\setlength{\parindent}{0pt}

\pagestyle{empty}

\begin{document}
\section*{Planning METOC Report}
Generated: %{GENERATED_DTG}

%{LOCATION_DATA}

\end{document}
`

const LocationDataTemplate = `
\subsection*{%{LOCATION_NAME}}

\subsubsection*{Sun and Moon Data}
\begin{tabular}{lrrrrrrrrcrr}
\toprule
	Date &  AT  &  NT &  CT &  SR &  SS &  CT &  NT &  AT &  LP &  MR &  MS \\
\midrule
	%{ASTRO_DATA_ROWS}
\bottomrule
\end{tabular}
`

const AstroDataTemplate = `
	%{DATE} &
	%{AT_BEGIN} &
	%{NT_BEGIN} &
	%{CT_BEGIN} &
	%{SR} &
	%{SS} &
	%{CT_END} &
	%{NT_END} &
	%{AT_END} &
	%{LP} &
	%{MR} &
	%{MS} \\
`

type LatexView struct {
	generated models.Dtg
	locations map[string]string
	astroDataRows map[string]string
	wxDataRows map[string]string
}

func NewLatexView(generated models.Dtg) LatexView {
	return LatexView{
		generated: generated,
		locations: make(map[string]string),
		astroDataRows: make(map[string]string),
		wxDataRows: make(map[string]string),
	}
}

func (lv *LatexView) AddLocation(key, name, mgrs string) {
	lv.locations[key] = fmt.Sprintf("%s (%s)", name, mgrs)
	lv.astroDataRows[key] = ""
	lv.wxDataRows[key] = ""
}

func (lv *LatexView) AddAstroData(locationKey string, date models.Dtg, data models.AstroData) {
	row := AstroDataTemplate

	row = strings.Replace(row, "%{DATE}", date.Date(), -1)
	row = strings.Replace(row, "%{AT_BEGIN}", data.AstronomicalTwilight.Begin.Time(), -1)
	row = strings.Replace(row, "%{NT_BEGIN}", data.NauticalTwilight.Begin.Time(), -1)
	row = strings.Replace(row, "%{CT_BEGIN}", data.CivilTwilight.Begin.Time(), -1)
	row = strings.Replace(row, "%{SR}", data.Sunrise.Time(), -1)
	row = strings.Replace(row, "%{SS}", data.Sunset.Time(), -1)
	row = strings.Replace(row, "%{CT_END}", data.CivilTwilight.End.Time(), -1)
	row = strings.Replace(row, "%{NT_END}", data.NauticalTwilight.End.Time(), -1)
	row = strings.Replace(row, "%{AT_END}", data.AstronomicalTwilight.End.Time(), -1)
	row = strings.Replace(row, "%{LP}", string(data.Phase), -1)
	row = strings.Replace(row, "%{MR}", data.MoonRise.Short(), -1)
	row = strings.Replace(row, "%{MS}", data.MoonSet.Short(), -1)

	lv.astroDataRows[locationKey] += row
}

func (lv LatexView) Build() string {
	loc := ""

	for k, v := range lv.locations {
		locStr := LocationDataTemplate
		locStr = strings.Replace(locStr, "%{LOCATION_NAME}", v, -1)
		locStr = strings.Replace(locStr, "%{ASTRO_DATA_ROWS}", lv.astroDataRows[k], -1)

		loc += locStr
	}

	out := Template
	out = strings.Replace(out, "%{GENERATED_DTG}", lv.generated.Full(), -1)
	return strings.Replace(out, "%{LOCATION_DATA}", loc, -1)
}