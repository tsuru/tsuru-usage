// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/db"
	"gopkg.in/mgo.v2/bson"
)

type TeamGroup struct {
	Name  string
	Teams []string
}

func FindTeamGroup(name string) (*TeamGroup, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	var group TeamGroup
	err = conn.TeamGroups().Find(bson.M{"name": name}).One(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func updateTeamGroup(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	err = r.ParseForm()
	if err != nil {
		return err
	}
	teams := r.Form["teams"]
	info, err := conn.TeamGroups().Upsert(bson.M{"name": name}, bson.M{"name": name, "teams": teams})
	if err != nil {
		return err
	}
	if info.Matched == 0 {
		w.WriteHeader(http.StatusCreated)
	}
	return nil
}

func listTeamGroups(w http.ResponseWriter, r *http.Request) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	var teamGroups []TeamGroup
	err = conn.TeamGroups().Find(nil).All(&teamGroups)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(teamGroups)
}
