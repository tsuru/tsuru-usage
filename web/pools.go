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

type PoolUsage struct {
	Pool  string
	Month string
	Usage float64
}

func poolListHandler(w http.ResponseWriter, r *http.Request) {
	pools := []string{"staging", "prod", "workshop"}
	render(w, "templates/pools/index.html", pools)
}

func poolUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pool := vars["name"]
	year := vars["year"]
	host := os.Getenv("HOST")
	url := fmt.Sprintf("%s/api/pools/%s/%s", host, pool, year)
	response, err := Client.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []PoolUsage
	err = json.NewDecoder(response.Body).Decode(&usage)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		PoolName string
		Year     string
		Usage    []PoolUsage
		Total    float64
	}{
		pool,
		year,
		usage,
		totalPoolUsage(usage),
	}
	err = render(w, "templates/pools/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalPoolUsage(usage []PoolUsage) float64 {
	var result float64
	for _, item := range usage {
		result += item.Usage
	}
	return result
}
