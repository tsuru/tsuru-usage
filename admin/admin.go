// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"crypto/subtle"
	"net/http"
	"os"
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
	m.HandleFunc("/teamgroups", basicAuth(groupListHandler))
	m.HandleFunc("/teamgroups/new", basicAuth(groupNewHandler)).Methods(http.MethodGet)
	m.HandleFunc("/teamgroups/{name}", basicAuth(groupEditHandler)).Methods(http.MethodGet)
	m.HandleFunc("/teamgroups/{name}", basicAuth(groupUpdateHandler)).Methods(http.MethodPost)
	m.HandleFunc("/teamgroups/{name}", basicAuth(groupDeleteHandler)).Methods(http.MethodDelete)
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if len(username) == 0 || len(password) == 0 {
		return handler
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="password"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}
		handler(w, r)
	}
}
