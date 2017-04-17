// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"net/http"
	"strconv"

	"github.com/tsuru/tsuru-usage/tsuru"
)

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type tsuruCollectorClient struct {
	api tsuru.TsuruAPI
}

type unitCount struct {
	app   string
	pool  string
	plan  string
	team  string
	count int
}

type serviceInstance struct {
	tsuru.ServiceInstance
	count int
}

func (c *tsuruCollectorClient) fetchUnitsCount() ([]unitCount, error) {
	apps, err := c.api.ListApps()
	if err != nil {
		return nil, err
	}
	var counts []unitCount
	for _, a := range apps {
		var count int
		for _, u := range a.Units {
			if u.Status == "started" || u.Status == "error" {
				count++
			}
		}
		counts = append(counts, unitCount{app: a.Name, pool: a.Pool, count: count, plan: a.Plan.Name, team: a.TeamOwner})
	}
	return counts, nil
}

func (c *tsuruCollectorClient) fetchNodesCount() (map[string]int, error) {
	nodes, err := c.api.ListNodes()
	if err != nil {
		return nil, err
	}
	count := make(map[string]int)
	for _, n := range nodes {
		count[n.Metadata.Pool]++
	}
	return count, nil
}

func (c *tsuruCollectorClient) fetchServicesInstances(service string) ([]serviceInstance, error) {
	instances, err := c.api.ListServiceInstances(service)
	if err != nil {
		return nil, err
	}
	serviceInstances := make([]serviceInstance, len(instances))
	for i := range instances {
		count := 1
		if str := instances[i].Info["Instances"]; str != "" {
			if v, err := strconv.Atoi(str); err == nil {
				count = v
			}
		}
		serviceInstances[i].ServiceInstance = instances[i]
		serviceInstances[i].count = count
	}
	return serviceInstances, nil
}
