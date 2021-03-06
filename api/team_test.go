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
	"github.com/tsuru/tsuru-usage/api/plan"
	"github.com/tsuru/tsuru-usage/db"
	"github.com/tsuru/tsuru-usage/prom"
	"github.com/tsuru/tsuru-usage/tsuru"
	check "gopkg.in/check.v1"
)

func (s *S) TestTeamList(c *check.C) {
	expected := []tsuru.Team{
		{Name: "team1"},
		{Name: "team2"},
		{Name: "team2"},
	}
	s.tsuruAPI.Teams = expected
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/teams", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []tsuru.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestGetTeamUsageGroup(c *check.C) {
	_, err := plan.Save(plan.PlanCost{Plan: "large", Type: plan.AppPlan, Cost: 3, MeasureUnit: "GB"})
	c.Assert(err, check.IsNil)
	expected := []TeamAppUsage{
		{Team: "mygroup", Month: "January", Usage: []AppUsage{{Plan: "default", Usage: 10}}},
		{Team: "mygroup", Month: "February", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "March", Usage: []AppUsage{
			{Plan: "default", Usage: 10},
			{Plan: "large", Usage: 5, Cost: UsageCost{MeasureUnit: "GB", UnitCost: 3, TotalCost: 15}}},
		},
		{Team: "mygroup", Month: "April", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "May", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "June", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "July", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "August", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "September", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "October", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "November", Usage: []AppUsage(nil)},
		{Team: "mygroup", Month: "December", Usage: []AppUsage(nil)},
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	vector := model.Vector{
		&model.Sample{Metric: model.Metric{"plan": "default"}, Value: model.SampleValue(10)},
		&model.Sample{Metric: model.Metric{"plan": "large"}, Value: model.SampleValue(5)},
	}
	fakeAPI.Add("tsuru_usage_units{team=~\"team1|team2\"}", "30d", s.nextDay(time.March), vector, "plan")
	fakeAPI.Add("tsuru_usage_units{team=~\"team1|team2\"}", "30d", s.nextDay(time.January), model.Vector{vector[0]}, "plan")
	prom.Client = fakeAPI
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Insert(TeamGroup{Name: "mygroup", Teams: []string{"team1", "team2"}})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/apps/mygroup/2017?group=true", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamAppUsage
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestGetTeamAppsUsage(c *check.C) {
	expected := []TeamAppUsage{
		{Team: "myteam", Month: "January", Usage: []AppUsage{{Plan: "default", Usage: 10}}},
		{Team: "myteam", Month: "February", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "March", Usage: []AppUsage{{Plan: "default", Usage: 10}, {Plan: "large", Usage: 5}}},
		{Team: "myteam", Month: "April", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "May", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "June", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "July", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "August", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "September", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "October", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "November", Usage: []AppUsage(nil)},
		{Team: "myteam", Month: "December", Usage: []AppUsage(nil)},
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	vector := model.Vector{
		&model.Sample{Metric: model.Metric{"plan": "default"}, Value: model.SampleValue(10)},
		&model.Sample{Metric: model.Metric{"plan": "large"}, Value: model.SampleValue(5)},
	}
	fakeAPI.Add("tsuru_usage_units{team=~\"myteam\"}", "30d", s.nextDay(time.March), vector, "plan")
	fakeAPI.Add("tsuru_usage_units{team=~\"myteam\"}", "30d", s.nextDay(time.January), model.Vector{vector[0]}, "plan")
	prom.Client = fakeAPI
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/apps/myteam/2017", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamAppUsage
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestGetTeamServicesUsage(c *check.C) {
	_, err := plan.Save(plan.PlanCost{Plan: "default", Service: "serv1", Type: plan.ServicePlan, Cost: 3, MeasureUnit: "GB"})
	c.Assert(err, check.IsNil)
	_, err = plan.Save(plan.PlanCost{Plan: "default", Service: "serv2", Type: plan.ServicePlan, Cost: 1, MeasureUnit: "GB"})
	c.Assert(err, check.IsNil)
	expected := []TeamServiceUsage{
		{Team: "myteam", Month: "January", Usage: []ServiceUsage{
			{Plan: "default", Service: "serv1", Usage: 10, Cost: UsageCost{MeasureUnit: "GB", UnitCost: 3, TotalCost: 30}}},
		},
		{Team: "myteam", Month: "February", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "March", Usage: []ServiceUsage{
			{Plan: "default", Service: "serv1", Usage: 10, Cost: UsageCost{MeasureUnit: "GB", UnitCost: 3, TotalCost: 30}},
			{Plan: "default", Service: "serv2", Usage: 5, Cost: UsageCost{MeasureUnit: "GB", UnitCost: 1, TotalCost: 5}}},
		},
		{Team: "myteam", Month: "April", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "May", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "June", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "July", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "August", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "September", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "October", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "November", Usage: []ServiceUsage(nil)},
		{Team: "myteam", Month: "December", Usage: []ServiceUsage(nil)},
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	vector := model.Vector{
		&model.Sample{Metric: model.Metric{"plan": "default", "service": "serv1"}, Value: model.SampleValue(10)},
		&model.Sample{Metric: model.Metric{"plan": "default", "service": "serv2"}, Value: model.SampleValue(5)},
	}
	fakeAPI.Add("tsuru_usage_services{team=~\"myteam\"}", "30d", s.nextDay(time.March), vector, "service", "plan")
	fakeAPI.Add("tsuru_usage_services{team=~\"myteam\"}", "30d", s.nextDay(time.January), model.Vector{vector[0]}, "service", "plan")
	prom.Client = fakeAPI
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/services/myteam/2017", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamServiceUsage
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}
