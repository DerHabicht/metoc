package orchestrations

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/derhabicht/metoc/models"
)

type ReportGenerator struct {
	plan           models.Plan
	reports        []models.MetocReport
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

func (r ReportGenerator) Generate() (string, error) {
	return "", nil
}
