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
	"sort"

	"github.com/tsuru/tsuru-usage/repositories"
)

type Pool struct {
	Name string
}

type Team struct {
	Name string
}

func poolListHandler(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/pools", host)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var pools []Pool
	err = json.NewDecoder(response.Body).Decode(&pools)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Name < pools[j].Name
	})
	render(w, "templates/list/pools.html", pools)
}

func teamListHandler(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teams", host)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var teams []Team
	err = json.NewDecoder(response.Body).Decode(&teams)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	render(w, "templates/list/teams.html", teams)
}

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := repositories.FetchGroups()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render(w, "templates/list/groups.html", groups)
}
