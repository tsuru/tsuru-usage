// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"net/http"

	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestFetchPools(c *check.C) {
	poolsData := `[
	{
		"Name": "pool 1"
	},
	{
		"Name": "pool 2"
	}
]`
	Client.Transport = &cmdtest.Transport{Message: poolsData, Status: http.StatusOK}
	pools, err := FetchPools()
	c.Assert(err, check.IsNil)
	c.Assert(pools, check.HasLen, 2)
	c.Assert(pools[0].Name, check.Equals, "pool 1")
	c.Assert(pools[1].Name, check.Equals, "pool 2")
}

func (s *S) TestFetchPoolsError(c *check.C) {
	Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	pools, err := FetchPools()
	c.Assert(err, check.NotNil)
	c.Assert(pools, check.IsNil)
}

func (s *S) TestFetchPoolsInvalidResponse(c *check.C) {
	Client.Transport = &cmdtest.Transport{Message: "not json", Status: http.StatusOK}
	pools, err := FetchPools()
	c.Assert(err, check.NotNil)
	c.Assert(pools, check.IsNil)
}
