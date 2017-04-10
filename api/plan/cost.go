// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plan

import (
	"github.com/tsuru/tsuru-usage/db"
	"gopkg.in/mgo.v2/bson"
)

const (
	AppPlan     = PlanType("app")
	ServicePlan = PlanType("service")
)

type PlanType string

type PlanCost struct {
	Service     string
	Plan        string
	Type        PlanType
	Cost        float64
	MeasureUnit string
}

func ListCosts() ([]PlanCost, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	var plans []PlanCost
	err = conn.PlanCosts().Find(nil).All(&plans)
	if err != nil {
		return nil, err
	}
	return plans, err
}

func Save(plan PlanCost) (bool, error) {
	conn, err := db.Conn()
	if err != nil {
		return false, err
	}
	info, err := conn.PlanCosts().Upsert(bson.M{"plan": plan.Plan, "service": plan.Service}, plan)
	if err != nil {
		return false, err
	}
	return info.Matched == 0, nil
}
