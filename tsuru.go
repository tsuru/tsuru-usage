// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type tsuruClient struct {
	addr       string
	token      string
	httpClient RequestDoer
}

type unitCount struct {
	app   string
	pool  string
	plan  string
	team  string
	count int
}

type nodeResult struct {
	Nodes []node
}

type node struct {
	Metadata nodeMetadata
}
type nodeMetadata struct {
	Pool string
}

type app struct {
	Name      string
	Plan      plan
	Units     []unit
	Pool      string
	TeamOwner string
}

type unit struct {
	Status string
}

type plan struct {
	Name string
}

func newClient(addr, token string) *tsuruClient {
	return &tsuruClient{addr: addr, token: token, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *tsuruClient) fetchUnitsCount() ([]unitCount, error) {
	var apps []app
	err := c.fetchList("apps", &apps)
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

func (c *tsuruClient) fetchNodesCount() (map[string]int, error) {
	var result nodeResult
	err := c.fetchList("node", &result)
	if err != nil {
		return nil, err
	}
	count := make(map[string]int)
	for _, n := range result.Nodes {
		count[n.Metadata.Pool]++
	}
	return count, nil
}

func (c *tsuruClient) fetchServicesInstances(services []string) ([]serviceInstance, error) {
	var result []serviceInstance
	for i := range services {
		var partial []serviceInstance
		err := c.fetchList("services/"+services[i], &partial)
		if err != nil {
			return nil, err
		}
		result = append(result, partial...)
	}
	return result, nil
}

type serviceInstance struct {
	ServiceName string
	Name        string
	PlanName    string
	TeamOwner   string
	Info        map[string]string
}

func (c *tsuruClient) fetchList(path string, v interface{}) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.addr, path), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "bearer "+c.token)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("returned non OK status code: %s", response.Status)
	}
	return json.NewDecoder(response.Body).Decode(v)
}
