// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prom

import (
	"testing"
	"time"

	check "gopkg.in/check.v1"

	"golang.org/x/net/context"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type S struct{}

var _ = check.Suite(&S{})

func Test(t *testing.T) { check.TestingT(t) }

type fakeQueryAPI struct {
	q string
	t time.Time
	m model.Value
	e error
}

func (q *fakeQueryAPI) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	q.q = query
	q.t = ts
	return q.m, q.e
}

func (q *fakeQueryAPI) QueryRange(ctx context.Context, query string, r prometheus.Range) (model.Value, error) {
	return nil, nil
}

func (s *S) TestAvgOverPeriod(c *check.C) {
	api := &fakeQueryAPI{
		m: model.Vector{
			{Value: model.SampleValue(10)},
		},
	}
	Client = &prometheusAPI{
		queryAPI: api,
	}
	d := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	f, err := GetAvgOverPeriod("metric{label=\"a\"}", "10d", d)
	c.Assert(err, check.IsNil)
	c.Assert(float64(f[0].Value), check.Equals, float64(10))
	c.Assert(api.q, check.DeepEquals, "avg(avg_over_time(metric{label=\"a\"}[10d]))")
	c.Assert(api.t, check.DeepEquals, d)
}

func (s *S) TestAvgOverPeriodGroupBy(c *check.C) {
	api := &fakeQueryAPI{
		m: model.Vector{
			{Value: model.SampleValue(10)},
		},
	}
	Client = &prometheusAPI{
		queryAPI: api,
	}
	d := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := GetAvgOverPeriod("metric{label=\"a\"}", "10d", d, "pool", "team")
	c.Assert(err, check.IsNil)
	c.Assert(api.q, check.DeepEquals, "avg(avg_over_time(metric{label=\"a\"}[10d])) by (pool,team)")
}
