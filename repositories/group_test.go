// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"net/http"

	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestFetchGroup(c *check.C) {
	groupData := `{
	"Name": "my-group",
	"Teams": ["team 1"],
	"Pools": ["pool 1", "pool 2"]
}`
	Client.Transport = &cmdtest.Transport{Message: groupData, Status: http.StatusOK}
	group, err := FetchGroup("my-group")
	c.Assert(err, check.IsNil)
	c.Assert(group, check.NotNil)
	c.Assert(group.Name, check.Equals, "my-group")
	c.Assert(group.Teams, check.DeepEquals, []string{"team 1"})
	c.Assert(group.Pools, check.DeepEquals, []string{"pool 1", "pool 2"})
}

func (s *S) TestFetchGroupNotFound(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusNotFound}
	group, err := FetchGroup("my-group")
	c.Assert(err, check.IsNil)
	c.Assert(group, check.IsNil)
}

func (s *S) TestFetchGroupError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	group, err := FetchGroup("my-group")
	c.Assert(err, check.NotNil)
	c.Assert(group, check.IsNil)
}

func (s *S) TestFetchGroupInvalidResponse(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "not json", Status: http.StatusOK}
	group, err := FetchGroup("my-group")
	c.Assert(err, check.NotNil)
	c.Assert(group, check.IsNil)
}

func (s *S) TestFetchGroups(c *check.C) {
	groupsData := `[
	{
		"Name": "my-group",
		"Teams": ["team 1"],
		"Pools": ["pool 1", "pool 2"]
	},
	{
		"Name": "other group",
		"Teams": [],
		"Pools": ["pool 2"]
	}
]`
	Client.Transport = &cmdtest.Transport{Message: groupsData, Status: http.StatusOK}
	groups, err := FetchGroups()
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.HasLen, 2)
	c.Assert(groups[0].Name, check.Equals, "my-group")
	c.Assert(groups[0].Teams, check.DeepEquals, []string{"team 1"})
	c.Assert(groups[0].Pools, check.DeepEquals, []string{"pool 1", "pool 2"})
	c.Assert(groups[1].Name, check.Equals, "other group")
	c.Assert(groups[1].Teams, check.DeepEquals, []string{})
	c.Assert(groups[1].Pools, check.DeepEquals, []string{"pool 2"})
}

func (s *S) TestFetchGroupsError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	groups, err := FetchGroups()
	c.Assert(err, check.NotNil)
	c.Assert(groups, check.IsNil)
}

func (s *S) TestFetchGroupsInvalidResponse(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "not json", Status: http.StatusOK}
	groups, err := FetchGroups()
	c.Assert(err, check.NotNil)
	c.Assert(groups, check.IsNil)
}
