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

func (s *S) TestAppUsage(c *check.C) {
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
				"Plan": "default plan",
				"Usage": 5,
				"Cost": {
					"MeasureUnit": "GB",
					"UnitCost": 2,
					"TotalCost": 10
				}
			}
		]
	},
	{
		"Month": "February",
		"Usage": [
			{
				"Plan": "planb",
				"Usage": 2,
				"Cost": {
					"MeasureUnit": "GB",
					"UnitCost": 3,
					"TotalCost": 6
				}
			}
		]
	}
]`
	repositories.Client.Transport = &cmdtest.Transport{Message: groupData, Status: http.StatusOK}
	Client.Transport = &cmdtest.Transport{Message: usageData, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/apps/mygroup/2017?group=true", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "default plan"), check.Equals, true)
	c.Assert(strings.Contains(body, "5.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "2.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "10.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "planb"), check.Equals, true)
	c.Assert(strings.Contains(body, "2.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "3.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "6.00 GB"), check.Equals, true)
	c.Assert(strings.Contains(body, "16.00 GB"), check.Equals, true)
}

func (s *S) TestAppUsageAPIError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/apps/mygroup/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestAppUsageInvalidJSON(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/apps/mygroup/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
