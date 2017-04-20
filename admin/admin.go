// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

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
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/teamgroups", 301)
	})
	m.HandleFunc("/teamgroups", groupListHandler)
	m.HandleFunc("/teamgroups/{name}", groupEditHandler).Methods(http.MethodGet)
	m.HandleFunc("/teamgroups/{name}", groupUpdateHandler).Methods(http.MethodPut)
}
