package api

import (
	"time"

	"github.com/prometheus/common/model"
	check "gopkg.in/check.v1"
)

func (s *S) TestRunForYear(c *check.C) {
	f := func(month time.Month) (model.Vector, error) {
		return model.Vector{&model.Sample{Value: model.SampleValue(month)}}, nil
	}
	results := runForYear(f)
	for m, v := range results {
		c.Assert(float64(v[0].Value), check.DeepEquals, float64(m))
	}
	c.Assert(len(results), check.Equals, 12)
}
