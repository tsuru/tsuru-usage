package api

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"encoding/json"

	"github.com/ajg/form"
	"github.com/tsuru/tsuru-usage/api/plan"
	check "gopkg.in/check.v1"
)

func (s *S) TestUpdatePlanCost(c *check.C) {
	recorder := httptest.NewRecorder()
	p := plan.PlanCost{
		Type:        plan.AppPlan,
		Plan:        "small",
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	reqBody, err := form.EncodeToString(p)
	c.Assert(err, check.IsNil)
	request, err := http.NewRequest(http.MethodPut, "/plans/cost", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	plans, err := plan.ListCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []plan.PlanCost{p})
	recorder = httptest.NewRecorder()
	p.Cost = 2
	reqBody, err = form.EncodeToString(p)
	c.Assert(err, check.IsNil)
	request, err = http.NewRequest(http.MethodPut, "/plans/cost", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	plans, err = plan.ListCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []plan.PlanCost{p})
}

func (s *S) TestListPlanCosts(c *check.C) {
	p1 := plan.PlanCost{
		Type:        plan.AppPlan,
		Plan:        "small",
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	p2 := plan.PlanCost{
		Type:        plan.ServicePlan,
		Service:     "rpaas",
		Plan:        "small",
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	_, err := plan.Save(p1)
	c.Assert(err, check.IsNil)
	_, err = plan.Save(p2)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/plans/cost", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	var plans []plan.PlanCost
	err = json.NewDecoder(recorder.Body).Decode(&plans)
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []plan.PlanCost{p1, p2})
}
