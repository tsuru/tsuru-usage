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
)

type Team struct {
	Name string
}

type Group struct {
	Name string
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
	render(w, "templates/list/teams.html", teams)
}

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teamgroups", host)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var groups []Group
	err = json.NewDecoder(response.Body).Decode(&groups)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render(w, "templates/list/groups.html", groups)
}
