// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"testing"

	"github.com/gorilla/mux"

	"gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

func Test(t *testing.T) { check.TestingT(t) }

func runServer() *mux.Router {
	r := mux.NewRouter()
	Router(r.PathPrefix("/web").Subrouter())
	return r
}
