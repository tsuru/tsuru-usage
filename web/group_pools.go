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

type GroupPoolUsage struct {
	Month string
	Usage []struct {
		Pool  string
		Usage UsageValue
	}
}

type TotalGroupPoolUsage struct {
	Total         UsageValue
	TotalPerMonth map[string]UsageValue
}

func groupPoolUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	year := vars["year"]
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teamgroups/%s/pools/%s", host, group, year)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []GroupPoolUsage
	err = json.NewDecoder(response.Body).Decode(&usage)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tabData := TabData{
		ActiveTab:    "pools",
		TeamOrGroup:  group,
		GroupingType: "group",
		Year:         year,
	}
	context := struct {
		Group      string
		Year       string
		Usage      []GroupPoolUsage
		TotalUsage TotalGroupPoolUsage
		TabData    TabData
	}{
		group,
		year,
		usage,
		totalGroupPoolUsage(usage),
		tabData,
	}
	err = render(w, "templates/group_pools/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalGroupPoolUsage(usage []GroupPoolUsage) TotalGroupPoolUsage {
	total := TotalGroupPoolUsage{TotalPerMonth: make(map[string]UsageValue)}
	for _, month := range usage {
		for _, item := range month.Usage {
			total.Total += item.Usage
			total.TotalPerMonth[month.Month] += item.Usage
		}
	}
	return total
}
