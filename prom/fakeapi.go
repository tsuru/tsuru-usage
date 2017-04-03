package prom

import "time"

type FakePrometheusAPI struct {
	results map[string]float64
}

func (p *FakePrometheusAPI) Add(selector, duration string, t time.Time, v float64) {
	if p.results == nil {
		p.results = make(map[string]float64)
	}
	p.results[selector+"/"+duration+"/"+t.String()] = v
}

func (p *FakePrometheusAPI) getAvgOverPeriod(selector, duration string, t time.Time) (float64, error) {
	if p.results == nil {
		p.results = make(map[string]float64)
	}
	return p.results[selector+"/"+duration+"/"+t.String()], nil
}
