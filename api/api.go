// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Router(m *mux.Router) {
	m.Handle("/pools/{name}/{year}", handler(getPoolUsage))
	m.Handle("/{resource}/{team}/{year}", handler(getTeamUsage))
	m.Handle("/teamgroups/{name}", handler(updateTeamGroup)).Methods(http.MethodPut)
	m.Handle("/teamgroups", handler(listTeamGroups))
}
