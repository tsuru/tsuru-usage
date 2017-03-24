// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type tsuruClient struct {
	addr       string
	token      string
	httpClient *http.Client
}

type unitCount struct {
	app   string
	pool  string
	team  string
	count int
}

type app struct {
	Name  string
	Units []unit
	Pool  string
}

type unit struct {
	Status string
}

func newClient(addr, token string) *tsuruClient {
	return &tsuruClient{addr: addr, token: token, httpClient: &http.Client{}}
}

func (c *tsuruClient) fetchUnitsCount() ([]unitCount, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/apps", c.addr), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "bearer "+c.token)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	defer response.Body.Close()
	var apps []app
	err = json.NewDecoder(response.Body).Decode(&apps)
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
		counts = append(counts, unitCount{app: a.Name, pool: a.Pool, count: count})
	}
	return counts, nil
}
