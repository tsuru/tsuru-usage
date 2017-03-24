package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type FakeDoer struct {
	response http.Response
}

func (d *FakeDoer) Do(request *http.Request) (*http.Response, error) {
	return &d.response, nil
}

func TestFetchNodesCount(t *testing.T) {
	body := `{
	"nodes": [
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "dev", "meta2": "bar"}},
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "dev", "meta2": "bar"}},
		{"Address": "http://localhost1:8080", "Status": "disabled", "Metadata": {"pool": "prod", "meta2": "bar"}}
	]
}`
	f := &FakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := tsuruClient{httpClient: f}
	counts, err := client.fetchNodesCount()
	if err != nil {
		t.Errorf("Expected err to be nil. Got %s", err)
	}
	expectedCounts := map[string]int{
		"dev":  2,
		"prod": 1,
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("Expected %#+v. Got %#+v", expectedCounts, counts)
	}
}

func TestFetchUnitsCount(t *testing.T) {
	body := `[
{"ip":"10.10.10.11","name":"app1","pool": "pool1", "units":[{"ID":"sapp1/0","Status":"started"}]},
{"ip":"10.10.10.11","name":"app3","pool": "pool2", "units":[{"ID":"sapp1/0","Status":"stopped"}]},
{"ip":"10.10.10.11","name":"app4","pool": "pool2"},
{"ip":"10.10.10.10","name":"app2", "pool":"pool1", "units":[{"ID":"app2/0","Status":"started"},{"ID":"app2/0","Status":"error"}]}]`
	f := &FakeDoer{
		response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
	}
	client := tsuruClient{httpClient: f}
	counts, err := client.fetchUnitsCount()
	if err != nil {
		t.Errorf("Expected err to be nil. Got %s", err)
	}
	expectedCounts := []unitCount{
		{app: "app1", pool: "pool1", count: 1},
		{app: "app3", pool: "pool2", count: 0},
		{app: "app4", pool: "pool2", count: 0},
		{app: "app2", pool: "pool1", count: 2},
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("Expected %#+v. Got %#+v", expectedCounts, counts)
	}
}
