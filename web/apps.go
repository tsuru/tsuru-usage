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
	Team  string
	Month string
	Usage []struct {
		Plan  string
		Usage int
	}
}

func appListHandler(w http.ResponseWriter, r *http.Request) {
	apps := []string{"app 1", "app 2", "app 3", "app 4"}
	render(w, "web/templates/apps/index.html", apps)
}

func appUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["name"]
	year := vars["year"]
	url := fmt.Sprintf("/api/apps/%s/%s", app, year)
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
		AppName string
		Year    string
		Usage   []AppUsage
	}{
		app,
		year,
		usage,
	}
	err = render(w, "web/templates/apps/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}
