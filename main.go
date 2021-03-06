// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tsuru/tsuru-usage/admin"
	"github.com/tsuru/tsuru-usage/api"
	"github.com/tsuru/tsuru-usage/exporter"
	"github.com/tsuru/tsuru-usage/tsuru"
	"github.com/tsuru/tsuru-usage/web"
	"github.com/urfave/negroni"
)

func main() {
	port := os.Getenv("PORT")
	tsuruEndpoint := os.Getenv("TSURU_HOST")
	tsuruToken := os.Getenv("USAGE_USER_TOKEN")
	tsuruServicesStr := os.Getenv("USAGE_SERVICES")
	if port == "" {
		port = "8080"
	}
	if tsuruEndpoint == "" {
		log.Fatal("Must set tsuru endpoint with TSURU_HOST env")
	}
	if tsuruToken == "" {
		log.Fatal("Must set tsuru token with USAGE_USER_TOKEN env")
	}
	var services []string
	if tsuruServicesStr != "" {
		services = strings.Split(tsuruServicesStr, ",")
	}
	client := tsuru.NewClient(tsuruEndpoint, tsuruToken)
	exporter.Register(client, services)
	runServer(port, client)
}

func runServer(port string, tsuruAPI tsuru.TsuruAPI) {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router(tsuruAPI),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("HTTP server listening at :%s...\n", port)
	log.Fatal(s.ListenAndServe())
}

func router(tsuruAPI tsuru.TsuruAPI) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusPermanentRedirect)
	})
	apiRouter := r.PathPrefix("/api").Subrouter()
	api.Router(apiRouter, tsuruAPI)
	webRouter := r.PathPrefix("/web").Subrouter()
	web.Router(webRouter)
	adminRouter := r.PathPrefix("/admin").Subrouter()
	admin.Router(adminRouter)
	n := negroni.Classic()
	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	n.UseHandler(r)
	return n
}
