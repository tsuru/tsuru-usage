// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/tsuru"
)

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleError(h(w, r), w, r)
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		log.Printf("[api] %s: %s", r.URL.Path, err.Error())
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Router(m *mux.Router, tsuruAPI tsuru.TsuruAPI) {
	m.Handle("/pools/{name}/{year}", handler(getPoolUsage))
	m.Handle("/{resource}/{team}/{year}", handler(getTeamUsage))
	m.Handle("/teamgroups/{name}/pools/{year}", handler(getPoolUsageForGroup))
	m.Handle("/teamgroups/{name}", handler(updateTeamGroup)).Methods(http.MethodPut)
	m.Handle("/teamgroups/{name}", handler(viewTeamGroup)).Methods(http.MethodGet)
	m.Handle("/teamgroups/{name}", handler(removeTeamGroup)).Methods(http.MethodDelete)
	m.Handle("/teamgroups", handler(listTeamGroups))
	m.Handle("/plans/cost", handler(updatePlanCost)).Methods(http.MethodPut)
	m.Handle("/plans/cost", handler(listPlanCosts))
	m.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		handleError(listTeams(w, r, tsuruAPI), w, r)
	})
	m.HandleFunc("/pools", func(w http.ResponseWriter, r *http.Request) {
		handleError(listPools(w, r, tsuruAPI), w, r)
	})
}
