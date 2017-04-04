// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/tsuru/tsuru-usage/prom"
)

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

func runForYear(f func(month time.Month) (model.Vector, error)) map[time.Month]model.Vector {
	results := make(chan monthUsage, 12)
	usage := make(map[time.Month]model.Vector)
	wg := sync.WaitGroup{}
	wg.Add(12)
	for m := 1; m <= 12; m++ {
		go func(m int) {
			result, err := f(time.Month(m))
			if err != nil {
				log.Printf("failed to get month %s usage: %s", time.Month(m).String(), err)
			}
			results <- monthUsage{month: time.Month(m), value: result}
			wg.Done()
		}(m)
	}
	wg.Wait()
	close(results)
	for u := range results {
		usage[u.month] = u.value
	}
	return usage
}

func poolMonthlyUsage(pool string, year int) func(month time.Month) (model.Vector, error) {
	sel := fmt.Sprintf("tsuru_usage_nodes{pool=%q}", pool)
	return func(month time.Month) (model.Vector, error) {
		t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		return prom.GetAvgOverPeriod(sel, "30d", t)
	}
}
