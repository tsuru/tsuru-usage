// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type fakeClient struct {
	units     []unitCount
	nodes     map[string]int
	instances []serviceInstance
}

func (c *fakeClient) fetchUnitsCount() ([]unitCount, error) {
	return c.units, nil
}

func (c *fakeClient) fetchNodesCount() (map[string]int, error) {
	return c.nodes, nil
}

func (c *fakeClient) fetchServicesInstances(service string) ([]serviceInstance, error) {
	return c.instances, nil
}

func TestCollectUnits(t *testing.T) {
	units := []unitCount{
		{app: "app1", pool: "pool1", count: 1, team: "admin"},
		{app: "app3", pool: "pool2", count: 0},
		{app: "app4", pool: "pool2", count: 0},
		{app: "app2", pool: "pool1", count: 2},
	}
	expectedLabels := []map[string]string{
		{"app": "app1", "pool": "pool1", "team": "admin", "plan": ""},
		{"app": "app3", "pool": "pool2", "team": "", "plan": ""},
		{"app": "app4", "pool": "pool2", "team": "", "plan": ""},
		{"app": "app2", "pool": "pool1", "team": "", "plan": ""},
	}
	client := &fakeClient{units: units}
	c := TsuruCollector{client: client}
	ch := make(chan prometheus.Metric, len(units))
	c.collectUnits(ch)
	for i := range units {
		checkMetric(<-ch, float64(units[i].count), expectedLabels[i], t)
	}
}

func TestCollectNodes(t *testing.T) {
	nodes := map[string]int{
		"pool1": 0,
		"pool2": 3,
		"pool3": 1,
	}
	client := &fakeClient{nodes: nodes}
	c := TsuruCollector{client: client}
	ch := make(chan prometheus.Metric, len(nodes))
	c.collectNodes(ch)
	for i := 0; i < len(nodes); i++ {
		node := <-ch
		metric := &dto.Metric{}
		node.Write(metric)
		value := metric.GetGauge().GetValue()
		label := metric.GetLabel()
		if label[0].GetName() != "pool" {
			t.Errorf("Expected pool label. Got %s", label[0].GetName())
		}
		if float64(nodes[label[0].GetValue()]) != value {
			t.Errorf("Expected value to be %d. Got %f.", nodes[label[0].GetValue()], value)
		}
	}
}

func TestCollectInstances(t *testing.T) {
	instances := []serviceInstance{
		{ServiceName: "rpaas", Name: "instance-rpaas", PlanName: "plan1", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1", "Instances": "2"}, count: 2},
		{ServiceName: "rpaas", Name: "instance-rpaas", TeamOwner: "myteam", Info: map[string]string{"Address": "127.0.0.1"}, count: 1},
	}
	expectedLabels := []map[string]string{
		{"service": "rpaas", "instance": "instance-rpaas", "plan": "plan1", "team": "myteam"},
		{"service": "rpaas", "instance": "instance-rpaas", "plan": "", "team": "myteam"},
	}
	client := &fakeClient{instances: instances}
	c := TsuruCollector{client: client, services: []string{"rpaas"}}
	ch := make(chan prometheus.Metric, len(instances))
	c.collectInstances(ch)
	for i := range instances {
		checkMetric(<-ch, float64(instances[i].count), expectedLabels[i], t)
	}
}

func checkMetric(m prometheus.Metric, expectedValue float64, expectedLabels map[string]string, t *testing.T) {
	metric := &dto.Metric{}
	m.Write(metric)
	value := metric.GetGauge().GetValue()
	if value != expectedValue {
		t.Errorf("Expected value to be %f. Got %f.", expectedValue, value)
	}
	label := metric.GetLabel()
	if len(label) != len(expectedLabels) {
		t.Errorf("Expected %d labels. Got %d.", len(expectedLabels), len(label))
	}
	for _, l := range label {
		expectedValue := expectedLabels[l.GetName()]
		valLabel := l.GetValue()
		if expectedValue != valLabel {
			t.Errorf("Expected label %s with value %s. Got %s.", l.GetName(), expectedValue, l.GetValue())
		}
	}
}
