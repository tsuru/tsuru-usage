// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var _ TsuruAPI = &TsuruClient{}

var Client TsuruAPI

type TsuruAPI interface {
	ListApps() ([]App, error)
	ListServiceInstances(service string) ([]ServiceInstance, error)
	ListNodes() ([]Node, error)
	ListPools() ([]Pool, error)
	ListTeams() ([]Team, error)
}

type App struct {
	Name      string
	Plan      Plan
	Units     []Unit
	Pool      string
	TeamOwner string
}

type Unit struct {
	Status string
}

type Plan struct {
	Name string
}

type ServiceInstance struct {
	ServiceName string
	Name        string
	PlanName    string
	TeamOwner   string
	Info        map[string]string
}

type nodeResult struct {
	Nodes []Node
}

type Node struct {
	Metadata NodeMetadata
}

type NodeMetadata struct {
	Pool string
}

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type TsuruClient struct {
	addr       string
	token      string
	httpClient RequestDoer
}

type Pool struct {
	Name string
}

type Team struct {
	Name string
}

func NewClient(addr, token string) *TsuruClient {
	return &TsuruClient{addr: addr, token: token, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *TsuruClient) ListApps() ([]App, error) {
	var apps []App
	err := c.fetchList("apps", &apps)
	return apps, err
}

func (c *TsuruClient) ListServiceInstances(service string) ([]ServiceInstance, error) {
	var result []ServiceInstance
	err := c.fetchList("services/"+service, &result)
	return result, err
}

func (c *TsuruClient) ListNodes() ([]Node, error) {
	var result nodeResult
	err := c.fetchList("node", &result)
	if err != nil {
		return nil, err
	}
	return result.Nodes, err
}

func (c *TsuruClient) ListPools() ([]Pool, error) {
	var result []Pool
	err := c.fetchList("pools", &result)
	return result, err
}

func (c *TsuruClient) ListTeams() ([]Team, error) {
	var result []Team
	err := c.fetchList("teams", &result)
	return result, err
}

func (c *TsuruClient) fetchList(path string, v interface{}) error {
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
