// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/tsuru/tsuru-usage/db"
	"github.com/tsuru/tsuru-usage/prom"
	check "gopkg.in/check.v1"
)

func server(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()
	Router(m)
	m.ServeHTTP(w, r)
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
	server(recorder, request)
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
		{Team: "mygroup", Month: "January", Usage: []PoolUsage{
			{Pool: "mypool", Month: "January", Usage: 10},
			{Pool: "mypool2", Month: "January", Usage: 5},
		}},
		{Team: "mygroup", Month: "February", Usage: nil},
		{Team: "mygroup", Month: "March", Usage: []PoolUsage{
			{Pool: "mypool", Month: "March", Usage: 10},
			{Pool: "mypool2", Month: "March", Usage: 5},
		}},
		{Team: "mygroup", Month: "April", Usage: nil},
		{Team: "mygroup", Month: "May", Usage: nil},
		{Team: "mygroup", Month: "June", Usage: nil},
		{Team: "mygroup", Month: "July", Usage: nil},
		{Team: "mygroup", Month: "August", Usage: nil},
		{Team: "mygroup", Month: "September", Usage: nil},
		{Team: "mygroup", Month: "October", Usage: nil},
		{Team: "mygroup", Month: "November", Usage: nil},
		{Team: "mygroup", Month: "December", Usage: []PoolUsage{
			{Pool: "mypool2", Month: "December", Usage: 5},
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
	server(recorder, request)
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
