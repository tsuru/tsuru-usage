// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestGroupPoolUsage(c *check.C) {
	data := `[
	{
		"Month": "January",
		"Usage": [
			{
				"Pool": "pool1",
				"Usage": 5
			},
			{
				"Pool": "pool2",
				"Usage": 7
			}
		]
	},
	{
		"Month": "February",
		"Usage": [
			{
				"Pool": "pool2",
				"Usage": 2
			}
		]
	}
]`
	Client.Transport = &cmdtest.Transport{Message: data, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool1"), check.Equals, true)
	c.Assert(strings.Contains(body, "5.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool2"), check.Equals, true)
	c.Assert(strings.Contains(body, "7.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "12.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool2"), check.Equals, true)
	c.Assert(strings.Contains(body, "2.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "Total"), check.Equals, true)
	c.Assert(strings.Contains(body, "14.00"), check.Equals, true)
}

func (s *S) TestGroupPoolUsageAPIError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestGroupPoolUsageInvalidJSON(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
