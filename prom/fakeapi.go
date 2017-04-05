// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prom

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/common/model"
)

type FakePrometheusAPI struct {
	results map[string]model.Vector
}

func key(selector, duration string, t time.Time, by ...string) string {
	keyFmt := "%s/%s/%s/%s"
	return fmt.Sprintf(keyFmt, selector, duration, t.String(), strings.Join(by, ","))
}

func (p *FakePrometheusAPI) Add(selector, duration string, t time.Time, v model.Vector, by ...string) {
	if p.results == nil {
		p.results = make(map[string]model.Vector)
	}
	p.results[key(selector, duration, t, by...)] = v
}

func (p *FakePrometheusAPI) getAvgOverPeriod(selector, duration string, t time.Time, by ...string) (model.Vector, error) {
	if p.results == nil {
		p.results = make(map[string]model.Vector)
	}
	return p.results[key(selector, duration, t, by...)], nil
}
