// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tsuru/tsuru-usage/repositories"
	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestServiceUsage(c *check.C) {
	groupData := `{
	"Name": "group 1",
	"Teams": ["team 1", "team 2"],
	"Pools": ["pool 1", "pool 2"]
}`
	usageData := `[
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
			},
			{
				"Service": "service 2",
				"Plan": "plan 3",
				"Usage": 6,
				"Cost": {
					"MeasureUnit": "GB",
					"UnitCost": 1,
					"TotalCost": 6
				}
			}
		]
	}
]`
	repositories.Client.Transport = &cmdtest.Transport{Message: groupData, Status: http.StatusOK}
	Client.Transport = &cmdtest.Transport{Message: usageData, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/services/mygroup/2017?group=true", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "plan 1"), check.Equals, true)
	c.Assert(strings.Contains(body, "11.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "2.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "22.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "plan 2"), check.Equals, true)
	c.Assert(strings.Contains(body, "4.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "3.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "12.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "plan 3"), check.Equals, true)
	c.Assert(strings.Contains(body, "6.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "1.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "6.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "18.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "Total"), check.Equals, true)
	c.Assert(strings.Contains(body, "21.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "40.00 GB"), check.Equals, true)
}

func (s *S) TestServiceUsageAPIError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/services/mygroup/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestServiceUsageInvalidJSON(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/services/mygroup/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
