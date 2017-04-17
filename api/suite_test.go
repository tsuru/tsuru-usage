// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/db"

	"github.com/tsuru/tsuru-usage/tsuru"
	check "gopkg.in/check.v1"
)

var _ = check.Suite(&S{})

type S struct {
	tsuruAPI *tsuru.FakeTsuruAPI
}

func Test(t *testing.T) { check.TestingT(t) }

func (s *S) SetUpTest(c *check.C) {
	os.Setenv("MONGODB_DATABASE_NAME", "tsuru_usage_api_tests")
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Database.DropDatabase()
	c.Assert(err, check.IsNil)
	s.tsuruAPI = &tsuru.FakeTsuruAPI{}
}

func (s *S) TearDownSuite(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Database.DropDatabase()
	c.Assert(err, check.IsNil)
}

func (s *S) server(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()
	Router(m, s.tsuruAPI)
	m.ServeHTTP(w, r)
}
