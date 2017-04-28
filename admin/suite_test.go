// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"errors"
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

func makeMultiTransport(urls []string, messages []string) *multiTransport {
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
	return &multiTransport{ConditionalTransports: cts}
}

type multiTransport struct {
	ConditionalTransports []cmdtest.ConditionalTransport
}

func (m *multiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, ct := range m.ConditionalTransports {
		if ct.CondFunc(req) {
			return ct.RoundTrip(req)
		}
	}
	return &http.Response{Body: nil, StatusCode: 500}, errors.New("condition failed")
}
