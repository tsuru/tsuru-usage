// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"strings"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

var Client PrometheusAPI

type PrometheusAPI interface {
	getAvgOverPeriod(selector, duration string, t time.Time, by ...string) (model.Vector, error)
}

type prometheusAPI struct {
	queryAPI prometheus.QueryAPI
}

func init() {
	client, err := prometheus.New(prometheus.Config{Address: os.Getenv("PROMETHEUS_HOST")})
	if err != nil {
		log.Fatalf("unable to initialize prometheus client: %s", err)
	}
	Client = &prometheusAPI{
		queryAPI: prometheus.NewQueryAPI(client),
	}
}

func GetAvgOverPeriod(selector, duration string, t time.Time, by ...string) (model.Vector, error) {
	return Client.getAvgOverPeriod(selector, duration, t, by...)
}

func (p *prometheusAPI) getAvgOverPeriod(selector, duration string, t time.Time, by ...string) (model.Vector, error) {
	query := fmt.Sprintf("avg(avg_over_time(%s[%s]))", selector, duration)
	if len(by) > 0 {
		query += " by (" + strings.Join(by, ",") + ")"
	}
	result, err := p.queryAPI.Query(context.Background(), query, t)
	if err != nil {
		return nil, err
	}
	vec, ok := result.(model.Vector)
	if !ok {
		return nil, errors.New("failed to parse result from query")
	}
	if len(vec) == 0 || vec == nil {
		return model.Vector{&model.Sample{Value: model.SampleValue(0)}}, nil
	}
	return vec, nil
}
