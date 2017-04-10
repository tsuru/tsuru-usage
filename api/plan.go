package api

import (
	"encoding/json"
	"net/http"

	"github.com/ajg/form"
	"github.com/tsuru/tsuru-usage/api/plan"
)

func updatePlanCost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	p := plan.PlanCost{}
	dec := form.NewDecoder(nil)
	dec.IgnoreCase(true)
	dec.IgnoreUnknownKeys(true)
	dec.DecodeValues(&p, r.Form)
	created, err := plan.Save(p)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}
	return nil
}

func listPlanCosts(w http.ResponseWriter, r *http.Request) error {
	plans, err := plan.ListCosts()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if len(plans) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	return json.NewEncoder(w).Encode(plans)
}
