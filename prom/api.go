package prom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

var Client PrometheusAPI

type PrometheusAPI interface {
	getAvgOverPeriod(selector, duration string, t time.Time) (float64, error)
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

func GetAvgOverPeriod(selector, duration string, t time.Time) (float64, error) {
	return Client.getAvgOverPeriod(selector, duration, t)
}

func (p *prometheusAPI) getAvgOverPeriod(selector, duration string, t time.Time) (float64, error) {
	query := fmt.Sprintf("avg(avg_over_time(%s[%s]))", selector, duration)
	result, err := p.queryAPI.Query(context.Background(), query, t)
	if err != nil {
		return 0, err
	}
	vec, ok := result.(model.Vector)
	if !ok {
		return 0, errors.New("failed to parse result from query")
	}
	if len(vec) == 0 {
		return 0, nil
	}
	return float64(vec[0].Value), nil
}
