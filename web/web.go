// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type handler func(http.ResponseWriter, *http.Request) error

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

// Router return a http.Handler with all web routes
func Router(m *mux.Router) {
	m.HandleFunc("/", indexHandler).Methods("GET")
	m.HandleFunc("/pools/", poolListHandler).Methods("GET")
	m.HandleFunc("/pools/{name}/{year}/", poolUsageHandler).Methods("GET")
	m.HandleFunc("/apps/", appTeamListHandler).Methods("GET")
	m.HandleFunc("/apps/{team}/{year}/", appUsageHandler).Methods("GET")
	m.HandleFunc("/services/", serviceTeamListHandler).Methods("GET")
	m.HandleFunc("/services/{team}/{year}/", serviceUsageHandler).Methods("GET")
}
