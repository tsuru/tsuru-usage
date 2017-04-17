// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	check "gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

func Test(t *testing.T) { check.TestingT(t) }

type fakeDoer struct {
	response http.Response
}

func (d *fakeDoer) Do(request *http.Request) (*http.Response, error) {
	return &d.response, nil
}

func (s *S) TestListNodes(c *check.C) {
	body := `{
	"nodes": [
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "dev", "meta2": "bar"}},
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "dev", "meta2": "bar"}},
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "prod", "meta2": "bar"}}
	]
}`
	f := &fakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := TsuruClient{httpClient: f}
	nodes, err := client.ListNodes()
	c.Assert(err, check.IsNil)
	expectedNodes := []Node{
		{Metadata: NodeMetadata{Pool: "dev"}},
		{Metadata: NodeMetadata{Pool: "dev"}},
		{Metadata: NodeMetadata{Pool: "prod"}},
	}
	c.Assert(nodes, check.DeepEquals, expectedNodes)
}

func (s *S) TestListApps(c *check.C) {
	body := `[
{"ip":"10.10.10.11","name":"app1","pool": "pool1", "teamowner":"admin", "units":[{"ID":"sapp1/0","Status":"started"}]},
{"ip":"10.10.10.11","name":"app3","pool": "pool2", "units":[{"ID":"sapp1/0","Status":"stopped"}]},
{"ip":"10.10.10.11","name":"app4","pool": "pool2"},
{"ip":"10.10.10.10","name":"app2", "pool":"pool1", "units":[{"ID":"app2/0","Status":"started"},{"ID":"app2/0","Status":"error"}]}]`
	f := &fakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := TsuruClient{httpClient: f}
	apps, err := client.ListApps()
	c.Assert(err, check.IsNil)
	expectedApps := []App{
		{Name: "app1", Pool: "pool1", TeamOwner: "admin", Units: []Unit{{Status: "started"}}},
		{Name: "app3", Pool: "pool2", Units: []Unit{{Status: "stopped"}}},
		{Name: "app4", Pool: "pool2"},
		{Name: "app2", Pool: "pool1", Units: []Unit{{Status: "started"}, {Status: "error"}}},
	}
	c.Assert(apps, check.DeepEquals, expectedApps)
}

func (s *S) TestListServicesInstances(c *check.C) {
	body := `[
	{"Apps":[],"Id":0,"Info":{"Address":"127.0.0.1","Instances":"2"},"Name":"instance-rpaas","PlanName":"plan1","ServiceName":"rpaas","TeamOwner":"myteam","Teams":["myteam"]},
	{"Apps":[],"Id":0,"Info":{"Address":"127.0.0.1"},"Name":"instance-rpaas","PlanName":"plan1","ServiceName":"rpaas","TeamOwner":"myteam","Teams":["myteam"]}
]`
	f := &fakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := TsuruClient{httpClient: f}
	instances, err := client.ListServiceInstances("rpaas")
	c.Assert(err, check.IsNil)
	expectedInstances := []ServiceInstance{
		{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1", "Instances": "2"}},
		{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1"}},
	}
	c.Assert(instances, check.DeepEquals, expectedInstances)
}

func (s *S) TestListPools(c *check.C) {
	body := `[
	{"Name":"pool1","Teams":[],"Public":true,"Default":false,"Provisioner":""},
	{"Name":"pool2","Teams":[],"Public":true,"Default":false,"Provisioner":""},
	{"Name":"pool3","Teams":[],"Public":true,"Default":false,"Provisioner":""}
]`
	f := &fakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := TsuruClient{httpClient: f}
	pools, err := client.ListPools()
	c.Assert(err, check.IsNil)
	expected := []Pool{
		{Name: "pool1"},
		{Name: "pool2"},
		{Name: "pool3"},
	}
	c.Assert(pools, check.DeepEquals, expected)
}

func (s *S) TestListTeams(c *check.C) {
	body := `[
	{"name":"team1","permissions":["app","team","service","service-instance"]},
	{"name":"team2","permissions":["app","team","service","service-instance"]},
	{"name":"team3","permissions":["app","team","service","service-instance"]}
]`
	f := &fakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := TsuruClient{httpClient: f}
	teams, err := client.ListTeams()
	c.Assert(err, check.IsNil)
	expected := []Team{
		{Name: "team1"},
		{Name: "team2"},
		{Name: "team3"},
	}
	c.Assert(teams, check.DeepEquals, expected)
}
