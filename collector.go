// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	unitsDesc  = prometheus.NewDesc("tsuru_usage_units", "The current number of started/errored units", []string{"app", "pool"}, nil)
	nodesDesc  = prometheus.NewDesc("tsuru_usage_nodes", "The current number of nodes", []string{"pool"}, nil)
	collectErr = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "tsuru_usage_collector_errors", Help: "The error count while fetching metrics"}, []string{"op"})
)

func init() {
	prometheus.MustRegister(collectErr)
}

type TsuruCollector struct {
	client *tsuruClient
}

func (c *TsuruCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- unitsDesc
	ch <- nodesDesc
}

func (c *TsuruCollector) Collect(ch chan<- prometheus.Metric) {
	unitsCounts, err := c.client.fetchUnitsCount()
	if err != nil {
		log.Printf("failed to fetch units metrics: %s", err)
		collectErr.WithLabelValues("units").Inc()
	}
	for _, u := range unitsCounts {
		ch <- prometheus.MustNewConstMetric(unitsDesc, prometheus.GaugeValue, float64(u.count), u.app, u.pool)
	}
	nodesCounts, err := c.client.fetchNodesCount()
	if err != nil {
		log.Printf("failed to fetch nodes metrics: %s", err)
		collectErr.WithLabelValues("nodes").Inc()
	}
	for p, c := range nodesCounts {
		ch <- prometheus.MustNewConstMetric(nodesDesc, prometheus.GaugeValue, float64(c), p)
	}
}
