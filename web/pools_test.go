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

func (s *S) TestPoolUsage(c *check.C) {
	data := `[
	{
		"Month": "January",
		"Usage": 10
	},
	{
		"Month": "February",
		"Usage": 6
	},
	{
		"Month": "March",
		"Usage": 3
	}
]`
	Client.Transport = &cmdtest.Transport{Message: data, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/pools/mypool/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "10"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "6"), check.Equals, true)
	c.Assert(strings.Contains(body, "March"), check.Equals, true)
	c.Assert(strings.Contains(body, "3"), check.Equals, true)
	c.Assert(strings.Contains(body, "Total"), check.Equals, true)
	c.Assert(strings.Contains(body, "19"), check.Equals, true)
}

func (s *S) TestPoolUsageAPIError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/pools/mypool/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestPoolUsageInvalidJSON(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/pools/mypool/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
