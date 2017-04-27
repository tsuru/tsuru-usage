// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

func Test(t *testing.T) { check.TestingT(t) }

func runServer() *mux.Router {
	r := mux.NewRouter()
	Router(r.PathPrefix("/admin").Subrouter())
	return r
}

func makeMultiConditionalTransport(messages []string) *cmdtest.MultiConditionalTransport {
	trueFunc := func(*http.Request) bool { return true }
	cts := make([]cmdtest.ConditionalTransport, len(messages))
	for i, message := range messages {
		cts[i] = cmdtest.ConditionalTransport{
			Transport: cmdtest.Transport{Message: message, Status: http.StatusOK},
			CondFunc:  trueFunc,
		}
	}
	return &cmdtest.MultiConditionalTransport{ConditionalTransports: cts}
}
