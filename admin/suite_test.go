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

func makeMultiConditionalTransport(urls []string, messages []string) *cmdtest.MultiConditionalTransport {
	cts := make([]cmdtest.ConditionalTransport, len(messages))
	for i, message := range messages {
		url := urls[i]
		cts[i] = cmdtest.ConditionalTransport{
			Transport: cmdtest.Transport{Message: message, Status: http.StatusOK},
			CondFunc: func(r *http.Request) bool {
				return r.URL.Path == url
			},
		}
	}
	return &cmdtest.MultiConditionalTransport{ConditionalTransports: cts}
}
