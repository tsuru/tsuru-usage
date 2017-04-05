package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"strings"

	"net/url"

	"github.com/tsuru/tsuru-usage/db"
	check "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *S) SetUpTest(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Database.DropDatabase()
	c.Assert(err, check.IsNil)
}

func (s *S) TestUpdateTeamGroup(c *check.C) {
	recorder := httptest.NewRecorder()
	params := url.Values{}
	params.Set("teams", "myteam")
	reqBody := strings.NewReader(params.Encode())
	request, err := http.NewRequest(http.MethodPut, "/teamgroups/mygroup", reqBody)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	var groups []TeamGroup
	err = conn.TeamGroups().Find(nil).All(&groups)
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.DeepEquals, []TeamGroup{{Name: "mygroup", Teams: []string{"myteam"}}})
	recorder = httptest.NewRecorder()
	params["teams"] = append(params["teams"], "mynewteam")
	reqBody = strings.NewReader(params.Encode())
	request, err = http.NewRequest(http.MethodPut, "/teamgroups/mygroup", reqBody)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	err = conn.TeamGroups().Find(nil).All(&groups)
	c.Assert(err, check.IsNil)
	c.Assert(groups, check.DeepEquals, []TeamGroup{{Name: "mygroup", Teams: []string{"myteam", "mynewteam"}}})
}

func (s *S) TestListTeamGroups(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	conn.TeamGroups().Insert(
		bson.M{"name": "group1", "teams": []string{"team1", "team2"}},
		bson.M{"name": "group2", "teams": []string{"team3"}},
	)
	expected := []TeamGroup{
		{Name: "group1", Teams: []string{"team1", "team2"}},
		{Name: "group2", Teams: []string{"team3"}},
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/teamgroups", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	var body []TeamGroup
	err = json.Unmarshal(recorder.Body.Bytes(), &body)
	c.Assert(err, check.IsNil)
	c.Assert(body, check.DeepEquals, expected)
	c.Assert(recorder.HeaderMap.Get("Content-type"), check.DeepEquals, "application/json")
}
