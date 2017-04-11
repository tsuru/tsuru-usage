// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/tsuru/tsuru-usage/prom"
)

type TeamPoolUsage struct {
	Team  string
	Month string
	Usage []PoolUsage
}

type PoolUsage struct {
	Pool  string
	Month string
	Usage float64
}

type monthUsage struct {
	month time.Month
	value model.Vector
}

func getPoolUsage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pool := vars["name"]
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return err
	}
	usage, err := poolYearUsage(pool, year)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(usage)
}

func getPoolUsageForGroup(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return err
	}
	poolSelector, err := poolSelectorForGroup(name)
	if err != nil {
		return err
	}
	usage, err := teamPoolsYearUsage(name, year, poolSelector)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(usage)
}

func teamPoolsYearUsage(team string, year int, selector string) ([]TeamPoolUsage, error) {
	result := runForYear(poolMonthlyUsage(selector, year))
	usage := make([]TeamPoolUsage, 12)
	for k, v := range result {
		var poolUsage []PoolUsage
		for _, val := range v {
			pool := string(val.Metric["pool"])
			value := float64(val.Value)
			poolUsage = append(poolUsage, PoolUsage{Pool: pool, Month: k.String(), Usage: value})
		}
		usage[k-1] = TeamPoolUsage{Team: team, Month: k.String(), Usage: poolUsage}
	}
	return usage, nil
}

func poolYearUsage(pool string, year int) ([]PoolUsage, error) {
	result := runForYear(poolMonthlyUsage(pool, year))
	usage := make([]PoolUsage, 12)
	for k, v := range result {
		var val float64
		if len(v) > 0 {
			val = float64(v[0].Value)
		}
		usage[k-1] = PoolUsage{Pool: pool, Month: k.String(), Usage: val}
	}
	return usage, nil
}

func poolMonthlyUsage(pool string, year int) func(month time.Month) (model.Vector, error) {
	sel := fmt.Sprintf("tsuru_usage_nodes{pool=~%q}", pool)
	return func(month time.Month) (model.Vector, error) {
		t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		return prom.GetAvgOverPeriod(sel, "30d", t, "pool")
	}
}

func poolSelectorForGroup(groupName string) (string, error) {
	group, err := FindTeamGroup(groupName)
	if err != nil {
		return "", err
	}
	return strings.Join(group.Pools, "|"), nil
}
