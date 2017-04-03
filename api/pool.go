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
	"github.com/tsuru/tsuru-usage/prom"
)

type monthUsage struct {
	month time.Month
	usage float64
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

func poolYearUsage(pool string, year int) (map[string]float64, error) {
	usages := make(chan monthUsage, 12)
	usage := make(map[string]float64)
	wg := sync.WaitGroup{}
	for m := 1; m <= 12; m++ {
		wg.Add(1)
		go func(m int) {
			defer wg.Done()
			usage, err := poolMonthUsage(pool, year, time.Month(m))
			if err != nil {
				log.Printf("failed to get month %s usage: %s", time.Month(m).String(), err)
			}
			usages <- monthUsage{month: time.Month(m), usage: usage}
		}(m)
	}
	wg.Wait()
	close(usages)
	for u := range usages {
		usage[u.month.String()] = u.usage
	}
	return usage, nil
}

func poolMonthUsage(pool string, year int, month time.Month) (float64, error) {
	t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	sel := fmt.Sprintf("tsuru_usage_nodes{pool=%q}", pool)
	return prom.GetAvgOverPeriod(sel, "30d", t)
}
