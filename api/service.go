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

func getTeamServicesUsage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	team := vars["team"]
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return err
	}
	usage, err := teamServicesYearUsage(team, year)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(usage)
}

func teamServicesYearUsage(team string, year int) ([]TeamServiceUsage, error) {
	result := runForYear(teamServicesMonthlyUsage(team, year))
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

func teamServicesMonthlyUsage(team string, year int) func(month time.Month) (model.Vector, error) {
	sel := fmt.Sprintf("tsuru_usage_services{team=%q}", team)
	return func(month time.Month) (model.Vector, error) {
		t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		return prom.GetAvgOverPeriod(sel, "30d", t, "service", "plan")
	}
}
