// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tsuru/tsuru-usage/tsuru"
)

var (
	unitsDesc    = prometheus.NewDesc("tsuru_usage_units", "The current number of started/errored units", []string{"app", "pool", "plan", "team"}, nil)
	nodesDesc    = prometheus.NewDesc("tsuru_usage_nodes", "The current number of nodes", []string{"pool"}, nil)
	servicesDesc = prometheus.NewDesc("tsuru_usage_services", "The current number of service instances", []string{"service", "instance", "team", "plan"}, nil)
	collectErr   = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "tsuru_usage_collector_errors", Help: "The error count while fetching metrics"}, []string{"op"})
	buckets      = append(prometheus.DefBuckets, []float64{15, 30, 45, 50}...)
	collectHist  = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "tsuru_usage_collector_duration_seconds", Help: "The duration of collector runs", Buckets: buckets})
)

type TsuruCollector struct {
	client   CollectableClient
	services []string
}

type CollectableClient interface {
	fetchUnitsCount() ([]unitCount, error)
	fetchNodesCount() (map[string]int, error)
	fetchServicesInstances(service string) ([]serviceInstance, error)
}

func init() {
	prometheus.MustRegister(collectErr)
	prometheus.MustRegister(collectHist)
}

func Register(tsuruAPI tsuru.TsuruAPI, services []string) {
	prometheus.MustRegister(&TsuruCollector{client: &tsuruCollectorClient{api: tsuruAPI}, services: services})
}

func (c *TsuruCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- unitsDesc
	ch <- nodesDesc
}

func (c *TsuruCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	defer func() {
		collectHist.Observe(time.Since(now).Seconds())
	}()
	wg := sync.WaitGroup{}
	collects := []func(chan<- prometheus.Metric){c.collectUnits, c.collectNodes, c.collectInstances}
	wg.Add(len(collects))
	for _, collect := range collects {
		go func(f func(chan<- prometheus.Metric)) {
			f(ch)
			wg.Done()
		}(collect)
	}
	wg.Wait()
}

func (c *TsuruCollector) collectUnits(ch chan<- prometheus.Metric) {
	unitsCounts, err := c.client.fetchUnitsCount()
	if err != nil {
		log.Printf("failed to fetch units metrics: %s", err)
		collectErr.WithLabelValues("units").Inc()
	}
	for _, u := range unitsCounts {
		ch <- prometheus.MustNewConstMetric(unitsDesc, prometheus.GaugeValue, float64(u.count), u.app, u.pool, u.plan, u.team)
	}
}

func (c *TsuruCollector) collectNodes(ch chan<- prometheus.Metric) {
	nodesCounts, err := c.client.fetchNodesCount()
	if err != nil {
		log.Printf("failed to fetch nodes metrics: %s", err)
		collectErr.WithLabelValues("nodes").Inc()
	}
	for p, c := range nodesCounts {
		ch <- prometheus.MustNewConstMetric(nodesDesc, prometheus.GaugeValue, float64(c), p)
	}
}

func (c *TsuruCollector) collectInstances(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.services))
	for _, s := range c.services {
		go func(s string) {
			instances, err := c.client.fetchServicesInstances(s)
			if err != nil {
				log.Printf("failed to fetch services metrics: %s", err)
				collectErr.WithLabelValues("services").Inc()
			}
			for _, i := range instances {
				ch <- prometheus.MustNewConstMetric(servicesDesc, prometheus.GaugeValue, float64(i.count), i.ServiceName, i.Name, i.TeamOwner, i.PlanName)
			}
			wg.Done()
		}(s)
	}
	wg.Wait()
}
