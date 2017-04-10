// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AppUsage struct {
	Month string
	Usage []struct {
		Plan  string
		Usage float64
	}
}

func appTeamListHandler(w http.ResponseWriter, r *http.Request) {
	teams := []string{"team1", "team2", "team3", "team4"}
	render(w, "web/templates/apps/index.html", teams)
}

func appUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]
	year := vars["year"]
	url := fmt.Sprintf("/api/apps/%s/%s", team, year)
	response, err := http.Get(url)
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
		Total float64
	}{
		team,
		year,
		usage,
		totalAppUsage(usage),
	}
	err = render(w, "web/templates/apps/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalAppUsage(usage []AppUsage) float64 {
	var result float64
	for _, month := range usage {
		for _, item := range month.Usage {
			result += item.Usage
		}
	}
	return result
}
