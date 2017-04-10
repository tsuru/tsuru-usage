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

type ServiceUsage struct {
	Month string
	Usage []struct {
		Service string
		Plan    string
		Usage   float64
	}
}

func serviceTeamListHandler(w http.ResponseWriter, r *http.Request) {
	teams := []string{"team 1", "team 2", "team 3", "team 4"}
	render(w, "web/templates/services/index.html", teams)
}

func serviceUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]
	year := vars["year"]
	url := fmt.Sprintf("/api/services/%s/%s", team, year)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []ServiceUsage
	json.NewDecoder(response.Body).Decode(&usage)
	context := struct {
		Team  string
		Year  string
		Usage []ServiceUsage
		Total float64
	}{
		team,
		year,
		usage,
		totalServiceUsage(usage),
	}
	err = render(w, "web/templates/services/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalServiceUsage(usage []ServiceUsage) float64 {
	var result float64
	for _, month := range usage {
		for _, item := range month.Usage {
			result += item.Usage
		}
	}
	return result
}
