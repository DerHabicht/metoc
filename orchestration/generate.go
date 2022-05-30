package orchestration

import (
	"github.com/derhabicht/metoc/controllers"
	"github.com/derhabicht/metoc/views"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/derhabicht/metoc/models"
)

type ReportGenerator struct {
	plan    models.Plan
	reports []models.MetocReport
}

func NewReportGenerator(planFile string) (ReportGenerator, error) {
	plan, err := ParsePlanFile(planFile)
	if err != nil {
		return ReportGenerator{}, errors.WithStack(err)
	}

	return ReportGenerator{plan, []models.MetocReport{}}, nil
}

func ParsePlanFile(planFile string) (models.Plan, error) {
	raw, err := os.ReadFile(planFile)
	if err != nil {
		return models.Plan{}, errors.WithStack(err)
	}

	plan := models.Plan{}
	err = yaml.Unmarshal(raw, &plan)
	if err != nil {
		return models.Plan{}, errors.WithStack(err)
	}

	return plan, nil
}

func loadPlanLocationDates(planLocation models.PlanLocation, report *models.MetocReport) error {
	for _, d := range planLocation.Dates {
		dtg, err := models.ParseIsoDateToDtg(d)
		if err != nil {
			return errors.WithStack(err)
		}

		report.AddDateToMetocReport(dtg)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (rg ReportGenerator) Generate(outfile string) error {
	for _, l := range rg.plan.Locations {
		r, err := models.NewMetocReport(l.Name, l.Mgrs, rg.plan.Tzoffset)
		if err != nil {
			return errors.WithStack(err)
		}

		err = loadPlanLocationDates(l, &r)
		if err != nil {
			return errors.WithStack(err)
		}

		rg.reports = append(rg.reports, r)
	}

	lv := views.NewLatexView(models.DtgNow(rg.plan.Tzoffset))

	for _, report := range rg.reports {
		lv.AddLocation(report.Location.Mgrs(), report.Location.Name, report.Location.Mgrs())

		err := controllers.FetchAstroDataForReport(&report)
		if err != nil {
			return errors.WithStack(err)
		}

		err = controllers.FetchWeatherDataForReport(&report)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, date := range report.Dates {
			lv.AddAstroData(report.Location.Mgrs(), date, report.AstroData[date])
			apf := views.EncodeDailyForecast(report.Forecast[date], report.Location, report.Generated)
			lv.AddWxData(report.Location.Mgrs(), date, report.Forecast[date], apf)
		}
	}

	err := os.WriteFile(outfile, []byte(lv.Build()), 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
