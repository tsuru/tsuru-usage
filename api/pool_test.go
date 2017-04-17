// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/prometheus/common/model"
	"github.com/tsuru/tsuru-usage/db"
	"github.com/tsuru/tsuru-usage/prom"
	"github.com/tsuru/tsuru-usage/tsuru"
	check "gopkg.in/check.v1"
)

func (s *S) TestPoolList(c *check.C) {
	expected := []tsuru.Pool{
		{Name: "pool1"},
		{Name: "pool2"},
		{Name: "pool2"},
	}
	s.tsuruAPI.Pools = expected
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/pools", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []tsuru.Pool
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestGetPoolUsage(c *check.C) {
	toVector := func(v float64) model.Vector {
		return model.Vector{&model.Sample{Value: model.SampleValue(v)}}
	}
	expected := []PoolUsage{
		{Pool: "mypool", Month: "January", Usage: 5},
		{Pool: "mypool", Month: "February", Usage: 0},
		{Pool: "mypool", Month: "March", Usage: 10},
		{Pool: "mypool", Month: "April", Usage: 0},
		{Pool: "mypool", Month: "May", Usage: 0},
		{Pool: "mypool", Month: "June", Usage: 0},
		{Pool: "mypool", Month: "July", Usage: 0},
		{Pool: "mypool", Month: "August", Usage: 0},
		{Pool: "mypool", Month: "September", Usage: 0},
		{Pool: "mypool", Month: "October", Usage: 0},
		{Pool: "mypool", Month: "November", Usage: 0},
		{Pool: "mypool", Month: "December", Usage: 2},
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool\"}", "30d", s.nextDay(time.March), toVector(10), "pool")
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool\"}", "30d", s.nextDay(time.January), toVector(5), "pool")
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool\"}", "30d", s.nextDay(time.December), toVector(2), "pool")
	prom.Client = fakeAPI
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/pools/mypool/2017", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []PoolUsage
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestGetPoolUsageForGroup(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Insert(TeamGroup{Name: "mygroup", Pools: []string{"mypool", "mypool2"}})
	c.Assert(err, check.IsNil)
	expected := []TeamPoolUsage{
		{Month: "January", Usage: []PoolUsage{
			{Pool: "mypool", Usage: 10},
			{Pool: "mypool2", Usage: 5},
		}},
		{Month: "February", Usage: nil},
		{Month: "March", Usage: []PoolUsage{
			{Pool: "mypool", Usage: 10},
			{Pool: "mypool2", Usage: 5},
		}},
		{Month: "April", Usage: nil},
		{Month: "May", Usage: nil},
		{Month: "June", Usage: nil},
		{Month: "July", Usage: nil},
		{Month: "August", Usage: nil},
		{Month: "September", Usage: nil},
		{Month: "October", Usage: nil},
		{Month: "November", Usage: nil},
		{Month: "December", Usage: []PoolUsage{
			{Pool: "mypool2", Usage: 5},
		}},
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	vector := model.Vector{
		&model.Sample{Metric: model.Metric{"pool": "mypool"}, Value: model.SampleValue(10)},
		&model.Sample{Metric: model.Metric{"pool": "mypool2"}, Value: model.SampleValue(5)},
	}
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool|mypool2\"}", "30d", s.nextDay(time.March), vector, "pool")
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool|mypool2\"}", "30d", s.nextDay(time.January), vector, "pool")
	fakeAPI.Add("tsuru_usage_nodes{pool=~\"mypool|mypool2\"}", "30d", s.nextDay(time.December), model.Vector{vector[1]}, "pool")
	prom.Client = fakeAPI
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamPoolUsage
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) nextDay(month time.Month) time.Time {
	return time.Date(2017, month+1, 1, 0, 0, 0, 0, time.UTC)
}
