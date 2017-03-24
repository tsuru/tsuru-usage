// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var unitsDesc = prometheus.NewDesc("tsuru_usage_units", "The current number of started/errored units", []string{"app", "pool"}, nil)

type TsuruCollector struct {
	client *tsuruClient
}

func (c *TsuruCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- unitsDesc
}

func (c *TsuruCollector) Collect(ch chan<- prometheus.Metric) {
	unitsCounts, err := c.client.fetchUnitsCount()
	if err != nil {
		log.Printf("failed to fetch units metrics: %s", err)
	}
	for _, u := range unitsCounts {
		ch <- prometheus.MustNewConstMetric(unitsDesc, prometheus.GaugeValue, float64(u.count), u.app, u.pool)
	}
}
