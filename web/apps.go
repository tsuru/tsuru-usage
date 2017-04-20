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
	"github.com/tsuru/tsuru-usage/repositories"
)

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

func appUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOrGroup := vars["teamOrGroup"]
	year := vars["year"]
	isGroup, _ := strconv.ParseBool(r.FormValue("group"))
	groupingType := "team"
	backURL := "/web/teams"
	var group *repositories.Group
	var err error
	if isGroup {
		groupingType = "group"
		group, err = repositories.FetchGroup(teamOrGroup)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		backURL = "/web/teamgroups"
	}
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/apps/%s/%s?group=%t", host, teamOrGroup, year, isGroup)
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
	tabData := TabData{
		ActiveTab:    "apps",
		TeamOrGroup:  teamOrGroup,
		GroupingType: groupingType,
		Year:         year,
	}
	context := struct {
		TeamOrGroup  string
		GroupingType string
		Year         string
		Usage        []AppUsage
		Total        TotalAppCost
		TabData      TabData
		Group        *repositories.Group
		BackURL      string
	}{
		teamOrGroup,
		groupingType,
		year,
		usage,
		totalAppCost(usage),
		tabData,
		group,
		backURL,
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
