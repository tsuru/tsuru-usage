package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/prom"
	check "gopkg.in/check.v1"
)

var _ = check.Suite(&S{})

type S struct{}

func Test(t *testing.T) { check.TestingT(t) }

func server(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()
	Router(m)
	m.ServeHTTP(w, r)
}

func (s *S) TestGetPoolUsage(c *check.C) {
	nextDay := func(month time.Month) time.Time {
		return time.Date(2017, month+1, 1, 0, 0, 0, 0, time.UTC)
	}
	expected := map[string]float64{
		"April":     0,
		"August":    0,
		"December":  2,
		"February":  0,
		"January":   5,
		"July":      0,
		"June":      0,
		"March":     10,
		"May":       0,
		"November":  0,
		"October":   0,
		"September": 0,
	}
	fakeAPI := &prom.FakePrometheusAPI{}
	fakeAPI.Add("tsuru_usage_nodes{pool=\"mypool\"}", "30d", nextDay(time.March), 10)
	fakeAPI.Add("tsuru_usage_nodes{pool=\"mypool\"}", "30d", nextDay(time.January), 5)
	fakeAPI.Add("tsuru_usage_nodes{pool=\"mypool\"}", "30d", nextDay(time.December), 2)
	prom.Client = fakeAPI
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/pool/mypool/2017", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body map[string]float64
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}
