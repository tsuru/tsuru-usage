// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"net/http"

	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestFetchTeams(c *check.C) {
	teamsData := `[
	{
		"Name": "team 1"
	},
	{
		"Name": "team 2"
	},
	{
		"Name": "team 3"
	}
]`
	Client.Transport = &cmdtest.Transport{Message: teamsData, Status: http.StatusOK}
	teams, err := FetchTeams()
	c.Assert(err, check.IsNil)
	c.Assert(teams, check.HasLen, 3)
	c.Assert(teams[0].Name, check.Equals, "team 1")
	c.Assert(teams[1].Name, check.Equals, "team 2")
	c.Assert(teams[2].Name, check.Equals, "team 3")
}

func (s *S) TestFetchTeamsError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	teams, err := FetchTeams()
	c.Assert(err, check.NotNil)
	c.Assert(teams, check.IsNil)
}

func (s *S) TestFetchTeamsInvalidResponse(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "not json", Status: http.StatusOK}
	teams, err := FetchTeams()
	c.Assert(err, check.NotNil)
	c.Assert(teams, check.IsNil)
}
