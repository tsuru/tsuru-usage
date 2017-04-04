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

type TeamAppUsage struct {
	Team  string
	Month string
	Usage []AppUsage
}

type AppUsage struct {
	Plan  string
	Usage float64
}

func getTeamAppsUsage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	team := vars["team"]
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return err
	}
	usage, err := teamAppsYearUsage(team, year)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(usage)
}

func teamAppsYearUsage(team string, year int) ([]TeamAppUsage, error) {
	result := runForYear(teamAppsMonthlyUsage(team, year))
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

func teamAppsMonthlyUsage(team string, year int) func(month time.Month) (model.Vector, error) {
	sel := fmt.Sprintf("tsuru_usage_units{team=%q}", team)
	return func(month time.Month) (model.Vector, error) {
		t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		return prom.GetAvgOverPeriod(sel, "30d", t, "plan")
	}
}
