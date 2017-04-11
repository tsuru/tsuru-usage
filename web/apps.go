// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var Client = &http.Client{}

type AppCost struct {
	MeasureUnit string
	UnitCost    float64
	TotalCost   float64
}

type TotalAppCost struct {
	AppCost
	Usage float64
}

type AppUsage struct {
	Month string
	Usage []struct {
		Plan  string
		Usage float64
		Cost  AppCost
	}
}

func appTeamListHandler(w http.ResponseWriter, r *http.Request) {
	teams := []string{"team1", "team2", "team3", "team4"}
	render(w, "templates/apps/index.html", teams)
}

func appUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]
	year := vars["year"]
	host := os.Getenv("HOST")
	url := fmt.Sprintf("%s/api/apps/%s/%s", host, team, year)
	response, err := Client.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []AppUsage
	json.NewDecoder(response.Body).Decode(&usage)
	context := struct {
		Team  string
		Year  string
		Usage []AppUsage
		Total TotalAppCost
	}{
		team,
		year,
		usage,
		totalAppCost(usage),
	}
	err = render(w, "templates/apps/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalAppCost(usage []AppUsage) TotalAppCost {
	total := TotalAppCost{}
	for _, month := range usage {
		for _, item := range month.Usage {
			if item.Cost.MeasureUnit != "" && total.MeasureUnit == "" {
				total.MeasureUnit = item.Cost.MeasureUnit
			}
			total.UnitCost += item.Cost.UnitCost
			total.TotalCost += item.Cost.TotalCost
			total.Usage += item.Usage
		}
	}
	return total
}
