package models

type PlanLocation struct {
	Name  string   `yaml:"name"`
	Mgrs  string   `yaml:"mgrs"`
	Dates []string `yaml:"dates"`
}

type Plan struct {
	Tzoffset  int            `yaml:"tzoffset"`
	Locations []PlanLocation `yaml:"locations"`
}
