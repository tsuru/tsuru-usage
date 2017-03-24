// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	addr := flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	tsuruEndpoint := flag.String("tsuru-address", "", "The tsuru API address to fetch resources from.")
	tsuruToken := flag.String("tsuru-token", "", "Tsuru API user token.")
	flag.Parse()

	if *tsuruEndpoint == "" {
		log.Fatal("Must set tsuru endpoint with \"--tsuru-address\" flag.")
	}
	if *tsuruToken == "" {
		log.Fatal("Must set tsuru token with \"--tsuru-token\" flag.")
	}

	http.Handle("/metrics", promhttp.Handler())

	tsuruClient := newClient(*tsuruEndpoint, *tsuruToken)
	prometheus.MustRegister(&TsuruCollector{client: tsuruClient})

	log.Printf("HTTP server listening at %s...\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
