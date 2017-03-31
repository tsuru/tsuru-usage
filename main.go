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
	"github.com/tsuru/tsuru-usage/api"
	"github.com/tsuru/tsuru-usage/exporter"
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
	exporter.Register(tsuruEndpoint, tsuruToken, services)
	runServer(port)
}

func runServer(port string) {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("HTTP server listening at :%s...\n", port)
	log.Fatal(s.ListenAndServe())
}

func router() http.Handler {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	api.Router(apiRouter)
	n := negroni.Classic()
	r.Handle("/metrics", promhttp.Handler())
	n.UseHandler(r)
	return n
}
