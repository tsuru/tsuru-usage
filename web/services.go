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

type ServiceCost struct {
	MeasureUnit string
	UnitCost    UsageValue
	TotalCost   UsageValue
}

type TotalServiceCost struct {
	ServiceCost
	Usage        UsageValue
	CostPerMonth map[string]UsageValue
}

type ServiceUsage struct {
	Month string
	Usage []struct {
		Service string
		Plan    string
		Usage   UsageValue
		Cost    ServiceCost
	}
}

func (s ServiceCost) UnitCostValue() string {
	str := s.UnitCost.String()
	if str == "0" {
		return str
	}
	return fmt.Sprintf("%s %s", str, s.MeasureUnit)
}

func (s ServiceCost) TotalCostValue() string {
	str := s.TotalCost.String()
	if str == "0" {
		return str
	}
	return fmt.Sprintf("%s %s", str, s.MeasureUnit)
}

func (t TotalServiceCost) MonthValue(month string) string {
	str := t.CostPerMonth[month].String()
	if str == "0" {
		return str
	}
	return fmt.Sprintf("%s %s", str, t.MeasureUnit)
}

func serviceUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOrGroup := vars["teamOrGroup"]
	year := vars["year"]
	group, _ := strconv.ParseBool(r.FormValue("group"))
	groupingType := "team"
	if group {
		groupingType = "group"
	}
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/services/%s/%s?group=%t", host, teamOrGroup, year, group)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []ServiceUsage
	err = json.NewDecoder(response.Body).Decode(&usage)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tabData := TabData{
		ActiveTab:    "services",
		TeamOrGroup:  teamOrGroup,
		GroupingType: groupingType,
		Year:         year,
	}
	context := struct {
		TeamOrGroup  string
		GroupingType string
		Year         string
		Usage        []ServiceUsage
		Total        TotalServiceCost
		TabData      TabData
	}{
		teamOrGroup,
		groupingType,
		year,
		usage,
		totalServiceCost(usage),
		tabData,
	}
	err = render(w, "templates/services/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalServiceCost(usage []ServiceUsage) TotalServiceCost {
	total := TotalServiceCost{CostPerMonth: make(map[string]UsageValue)}
	for _, month := range usage {
		for _, item := range month.Usage {
			if item.Cost.MeasureUnit != "" && total.MeasureUnit == "" {
				total.MeasureUnit = item.Cost.MeasureUnit
			}
			total.TotalCost += item.Cost.TotalCost
			total.Usage += item.Usage
			total.CostPerMonth[month.Month] += UsageValue(item.Cost.TotalCost)
		}
	}
	return total
}
