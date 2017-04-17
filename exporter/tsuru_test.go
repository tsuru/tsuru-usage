// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"reflect"
	"testing"

	"github.com/tsuru/tsuru-usage/tsuru"
)

func TestFetchNodesCount(t *testing.T) {
	f := &tsuru.FakeTsuruAPI{
		Nodes: []tsuru.Node{
			{Metadata: tsuru.NodeMetadata{Pool: "dev"}},
			{Metadata: tsuru.NodeMetadata{Pool: "dev"}},
			{Metadata: tsuru.NodeMetadata{Pool: "prod"}},
		},
	}
	client := tsuruCollectorClient{api: f}
	counts, err := client.fetchNodesCount()
	if err != nil {
		t.Errorf("Expected err to be nil. Got %s", err)
	}
	expectedCounts := map[string]int{
		"dev":  2,
		"prod": 1,
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("Expected %#+v. Got %#+v", expectedCounts, counts)
	}
}

func TestFetchUnitsCount(t *testing.T) {
	f := &tsuru.FakeTsuruAPI{
		Apps: []tsuru.App{
			{Name: "app1", Pool: "pool1", TeamOwner: "admin", Units: []tsuru.Unit{{Status: "started"}}},
			{Name: "app3", Pool: "pool2", Units: []tsuru.Unit{{Status: "stopped"}}},
			{Name: "app4", Pool: "pool2"},
			{Name: "app2", Pool: "pool1", Units: []tsuru.Unit{{Status: "started"}, {Status: "error"}}},
		},
	}
	client := tsuruCollectorClient{api: f}
	counts, err := client.fetchUnitsCount()
	if err != nil {
		t.Errorf("Expected err to be nil. Got %s", err)
	}
	expectedCounts := []unitCount{
		{app: "app1", pool: "pool1", count: 1, team: "admin"},
		{app: "app3", pool: "pool2", count: 0},
		{app: "app4", pool: "pool2", count: 0},
		{app: "app2", pool: "pool1", count: 2},
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("Expected %#+v. Got %#+v", expectedCounts, counts)
	}
}

func TestFetchServicesInstances(t *testing.T) {
	f := &tsuru.FakeTsuruAPI{
		Instances: map[string][]tsuru.ServiceInstance{
			"rpaas": []tsuru.ServiceInstance{
				{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1", "Instances": "2"}},
				{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1"}},
			},
		},
	}
	client := tsuruCollectorClient{api: f}
	instances, err := client.fetchServicesInstances("rpaas")
	if err != nil {
		t.Errorf("Expected err to be nil. Got %s", err)
	}
	expectedInstances := []serviceInstance{
		{ServiceInstance: tsuru.ServiceInstance{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1", "Instances": "2"}}, count: 2},
		{ServiceInstance: tsuru.ServiceInstance{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1"}}, count: 1},
	}
	if !reflect.DeepEqual(instances, expectedInstances) {
		t.Errorf("Expected %#+v. Got %#+v", expectedInstances, instances)
	}
}
