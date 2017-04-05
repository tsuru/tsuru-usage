package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/tsuru/tsuru-usage/prom"
)

const (
	Services = ResourceType("services")
	Apps     = ResourceType("apps")
)

type ResourceType string

func getTeamUsage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	team := vars["team"]
	resource := vars["resource"]
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return err
	}
	err = r.ParseForm()
	if err != nil {
		return err
	}
	group, _ := strconv.ParseBool(r.Form.Get("group"))
	teamSel := team
	if group {
		var err error
		teamSel, err = selectorForGroup(team)
		if err != nil {
			return err
		}
	}
	var usage interface{}
	switch ResourceType(resource) {
	case Services:
		usage, err = teamServicesYearUsage(team, year, teamSel)
	case Apps:
		usage, err = teamAppsYearUsage(team, year, teamSel)
	default:
		w.WriteHeader(http.StatusBadRequest)
		err = fmt.Errorf("invalid resource type: %q", resource)
	}
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(usage)
}

type TeamAppUsage struct {
	Team  string
	Month string
	Usage []AppUsage
}

type AppUsage struct {
	Plan  string
	Usage float64
}

func teamAppsYearUsage(team string, year int, teamSelector string) ([]TeamAppUsage, error) {
	result := runForYear(monthlyUsage("tsuru_usage_units", team, year, "plan"))
	usage := make([]TeamAppUsage, 12)
	for k, v := range result {
		var appUsage []AppUsage
		for _, val := range v {
			plan := val.Metric["plan"]
			appUsage = append(appUsage, AppUsage{Plan: string(plan), Usage: float64(val.Value)})
		}
		usage[k-1] = TeamAppUsage{Team: team, Month: k.String(), Usage: appUsage}
	}
	return usage, nil
}

func monthlyUsage(metric, team string, year int, by ...string) func(month time.Month) (model.Vector, error) {
	sel := fmt.Sprintf("%s{team=~%q}", metric, team)
	return func(month time.Month) (model.Vector, error) {
		t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		return prom.GetAvgOverPeriod(sel, "30d", t, by...)
	}
}

type TeamServiceUsage struct {
	Team  string
	Month string
	Usage []ServiceUsage
}

type ServiceUsage struct {
	Service string
	Plan    string
	Usage   float64
}

func teamServicesYearUsage(team string, year int, teamSelector string) ([]TeamServiceUsage, error) {
	result := runForYear(monthlyUsage("tsuru_usage_services", team, year, "service", "plan"))
	usage := make([]TeamServiceUsage, 12)
	for k, v := range result {
		var servUsage []ServiceUsage
		for _, val := range v {
			plan := val.Metric["plan"]
			service := val.Metric["service"]
			servUsage = append(servUsage, ServiceUsage{Plan: string(plan), Service: string(service), Usage: float64(val.Value)})
		}
		usage[k-1] = TeamServiceUsage{Team: team, Month: k.String(), Usage: servUsage}
	}
	return usage, nil
}
