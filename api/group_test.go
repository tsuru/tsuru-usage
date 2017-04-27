// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"net/url"

	"github.com/tsuru/tsuru-usage/db"
	check "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *S) TestUpdateTeamGroup(c *check.C) {
	recorder := httptest.NewRecorder()
	params := url.Values{}
	params.Set("teams", "myteam")
	params.Set("pools", "mypool")
	reqBody := strings.NewReader(params.Encode())
	request, err := http.NewRequest(http.MethodPut, "/teamgroups/mygroup", reqBody)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	var groups []TeamGroup
	err = conn.TeamGroups().Find(nil).All(&groups)
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.DeepEquals, []TeamGroup{{Name: "mygroup", Teams: []string{"myteam"}, Pools: []string{"mypool"}}})
	recorder = httptest.NewRecorder()
	params["teams"] = append(params["teams"], "mynewteam")
	params.Del("pools")
	reqBody = strings.NewReader(params.Encode())
	request, err = http.NewRequest(http.MethodPut, "/teamgroups/mygroup", reqBody)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	err = conn.TeamGroups().Find(nil).All(&groups)
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.DeepEquals, []TeamGroup{{Name: "mygroup", Teams: []string{"myteam", "mynewteam"}, Pools: []string{}}})
}

func (s *S) TestListTeamGroups(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	conn.TeamGroups().Insert(
		bson.M{"name": "group1", "teams": []string{"team1", "team2"}},
		bson.M{"name": "group2", "teams": []string{"team3"}, "pools": []string{"pool1"}},
	)
	expected := []TeamGroup{
		{Name: "group1", Teams: []string{"team1", "team2"}},
		{Name: "group2", Teams: []string{"team3"}, Pools: []string{"pool1"}},
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/teamgroups", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamGroup
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestViewTeamGroup(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	conn.TeamGroups().Insert(
		bson.M{"name": "group1", "teams": []string{"team1", "team2"}},
		bson.M{"name": "group2", "teams": []string{"team3"}, "pools": []string{"pool1"}},
	)
	expected := TeamGroup{Name: "group1", Teams: []string{"team1", "team2"}}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/teamgroups/group1", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body TeamGroup
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}

func (s *S) TestRemoveTeamGroup(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	conn.TeamGroups().Insert(
		bson.M{"name": "group1", "teams": []string{"team1", "team2"}},
		bson.M{"name": "group2", "teams": []string{"team3"}, "pools": []string{"pool1"}},
	)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, "/teamgroups/group1", nil)
	c.Assert(err, check.IsNil)
	s.server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var groups []TeamGroup
	err = conn.TeamGroups().Find(nil).All(&groups)
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.DeepEquals, []TeamGroup{{Name: "group2", Teams: []string{"team3"}, Pools: []string{"pool1"}}})
}
