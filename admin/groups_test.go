// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/tsuru/tsuru-usage/repositories"
	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestGroupList(c *check.C) {
	groupsData := `[
	{
		"Name": "my-group",
		"Teams": ["team 1"],
		"Pools": ["pool 1", "pool 2"]
	},
	{
		"Name": "other group",
		"Teams": ["team 2", "team 3"],
		"Pools": ["pool 3"]
	}
]`
	repositories.Client.Transport = &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: groupsData, Status: http.StatusOK},
		CondFunc: func(r *http.Request) bool {
			return r.URL.Path == "/api/teamgroups"
		},
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "my-group"), check.Equals, true)
	c.Assert(strings.Contains(body, "team 1"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool 1, pool 2"), check.Equals, true)
	c.Assert(strings.Contains(body, "other group"), check.Equals, true)
	c.Assert(strings.Contains(body, "team 2, team 3"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool 3"), check.Equals, true)
}

func (s *S) TestGroupListWithError(c *check.C) {
	repositories.Client.Transport = &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Status: http.StatusInternalServerError},
		CondFunc: func(r *http.Request) bool {
			return r.URL.Path == "/api/teamgroups"
		},
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestGroupNew(c *check.C) {
	teamsData := `[
	{
		"Name": "team a"
	},
	{
		"Name": "team b"
	}
]`
	poolsData := `[
	{
		"Name": "pool a"
	},
	{
		"Name": "pool b"
	}
]`
	repositories.Client.Transport = makeAnyConditionalTransport([]string{"/api/teams", "/api/pools"}, []string{teamsData, poolsData})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups/new", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "<option>team a</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>team b</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>pool a</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>pool b</option>"), check.Equals, true)
}

func (s *S) TestGroupEdit(c *check.C) {
	groupData := `{
	"Name": "mygroup",
	"Teams": [],
	"Pools": []
}`
	teamsData := `[
	{
		"Name": "team a"
	},
	{
		"Name": "team b"
	}
]`
	poolsData := `[
	{
		"Name": "pool a"
	},
	{
		"Name": "pool b"
	}
]`
	repositories.Client.Transport = makeAnyConditionalTransport([]string{"/api/teamgroups/mygroup", "/api/teams", "/api/pools"}, []string{groupData, teamsData, poolsData})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups/mygroup", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "<option>team a</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>team b</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>pool a</option>"), check.Equals, true)
	c.Assert(strings.Contains(body, "<option>pool b</option>"), check.Equals, true)
}

func (s *S) TestGroupEditGroupNotFound(c *check.C) {
	repositories.Client.Transport = &cmdtest.AnyConditionalTransport{ConditionalTransports: []cmdtest.ConditionalTransport{
		cmdtest.ConditionalTransport{
			Transport: cmdtest.Transport{Status: http.StatusNotFound},
			CondFunc: func(r *http.Request) bool {
				return r.URL.Path == "/api/teamgroups/mygroup"
			},
		},
		cmdtest.ConditionalTransport{
			Transport: cmdtest.Transport{Message: "[]", Status: http.StatusOK},
			CondFunc: func(r *http.Request) bool {
				return true
			},
		},
	}}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups/mygroup", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestGroupEditRequestError(c *check.C) {
	repositories.Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/admin/teamgroups/mygroup", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestGroupUpdate(c *check.C) {
	repositories.Client.Transport = &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Status: http.StatusOK},
		CondFunc: func(r *http.Request) bool {
			return r.URL.Path == "/api/teamgroups/mygroup" && r.Method == http.MethodPut
		},
	}
	recorder := httptest.NewRecorder()
	v := url.Values{"teams": []string{"team 1"}, "pools": []string{"pool 1", "pool 2"}}
	request, err := http.NewRequest(http.MethodPost, "/admin/teamgroups/mygroup", strings.NewReader(v.Encode()))
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
}

func (s *S) TestGroupUpdateError(c *check.C) {
	repositories.Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	v := url.Values{"teams": []string{"team 1"}, "pools": []string{"pool 1", "pool 2"}}
	request, err := http.NewRequest(http.MethodPost, "/admin/teamgroups/mygroup", strings.NewReader(v.Encode()))
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestGroupDelete(c *check.C) {
	repositories.Client.Transport = &cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Status: http.StatusOK},
		CondFunc: func(r *http.Request) bool {
			return r.URL.Path == "/api/teamgroups/mygroup" && r.Method == http.MethodDelete
		},
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, "/admin/teamgroups/mygroup", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}

func (s *S) TestGroupDeleteError(c *check.C) {
	repositories.Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, "/admin/teamgroups/mygroup", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
