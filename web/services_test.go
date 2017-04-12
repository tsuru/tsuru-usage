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

func (s *S) TestServiceUsage(c *check.C) {
	data := `[
	{
		"Month": "January",
		"Usage": [
			{
				"Service": "service 1",
				"Plan": "plan 1",
				"Usage": 11,
				"Cost": {
					"MeasureUnit": "GB",
					"UnitCost": 2,
					"TotalCost": 22
				}
			}
		]
	},
	{
		"Month": "February",
		"Usage": [
			{
				"Service": "service 2",
				"Plan": "plan 2",
				"Usage": 4,
				"Cost": {
					"MeasureUnit": "GB",
					"UnitCost": 3,
					"TotalCost": 12
				}
			}
		]
	}
]`
	Client.Transport = &cmdtest.Transport{Message: data, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/services/mygroup/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "plan 1"), check.Equals, true)
	c.Assert(strings.Contains(body, "11"), check.Equals, true)
	c.Assert(strings.Contains(body, "2 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "22 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "plan 2"), check.Equals, true)
	c.Assert(strings.Contains(body, "4"), check.Equals, true)
	c.Assert(strings.Contains(body, "3 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "12 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "Total"), check.Equals, true)
	c.Assert(strings.Contains(body, "15"), check.Equals, true)
	c.Assert(strings.Contains(body, "5 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "34 GB"), check.Equals, true)
}

func (s *S) TestServiceUsageAPIError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/services/mygroup/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestServiceUsageInvalidJSON(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/web/services/mygroup/2017/", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
