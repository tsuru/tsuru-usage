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
	"strconv"

	"github.com/gorilla/mux"
)

var Client = &http.Client{}

type AppCost struct {
	MeasureUnit string
	UnitCost    UsageValue
	TotalCost   UsageValue
}

type TotalAppCost struct {
	AppCost
	Usage UsageValue
}

type AppUsage struct {
	Month string
	Usage []struct {
		Plan  string
		Usage UsageValue
		Cost  AppCost
	}
}

func (a AppCost) UnitCostValue() string {
	str := a.UnitCost.String()
	if str == "0" {
		return str
	}
	return fmt.Sprintf("%s %s", str, a.MeasureUnit)
}

func (a AppCost) TotalCostValue() string {
	str := a.TotalCost.String()
	if str == "0" {
		return str
	}
	return fmt.Sprintf("%s %s", str, a.MeasureUnit)
}

func appTeamListHandler(w http.ResponseWriter, r *http.Request) {
	teams := []string{"team1", "team2", "team3", "team4"}
	render(w, "templates/apps/index.html", teams)
}

func appUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOrGroup := vars["teamOrGroup"]
	year := vars["year"]
	group, _ := strconv.ParseBool(r.FormValue("group"))
	groupingType := "team"
	if group {
		groupingType = "group"
	}
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/apps/%s/%s?group=%t", host, teamOrGroup, year, group)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []AppUsage
	err = json.NewDecoder(response.Body).Decode(&usage)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		TeamOrGroup  string
		GroupingType string
		Year         string
		Usage        []AppUsage
		Total        TotalAppCost
	}{
		teamOrGroup,
		groupingType,
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
