// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type TabData struct {
	ActiveTab    string
	TeamOrGroup  string
	GroupingType string
	Year         string
}

type UsageValue float64

type handler func(http.ResponseWriter, *http.Request) error

func (u UsageValue) String() string {
	if u == 0.0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", u)
}

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var Client = &http.Client{}

// Router return a http.Handler with all web routes
func Router(m *mux.Router) {
	m.HandleFunc("/", indexHandler).Methods(http.MethodGet)
	m.HandleFunc("/pools", poolListHandler).Methods(http.MethodGet)
	m.HandleFunc("/teams", teamListHandler).Methods(http.MethodGet)
	m.HandleFunc("/teamgroups", groupListHandler).Methods(http.MethodGet)
	m.HandleFunc("/pools/{name}/{year}", poolUsageHandler).Methods(http.MethodGet)
	m.HandleFunc("/teamgroups/{group}/pools/{year}", groupPoolUsageHandler).Methods(http.MethodGet)
	m.HandleFunc("/apps/{teamOrGroup}/{year}", appUsageHandler).Methods(http.MethodGet)
	m.HandleFunc("/services/{teamOrGroup}/{year}", serviceUsageHandler).Methods(http.MethodGet)
}
