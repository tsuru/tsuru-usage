// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/repositories"
)

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := repositories.FetchGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		Groups []repositories.Group
	}{
		groups,
	}
	err = render(w, "templates/groups/index.html", context)
	if err != nil {
		log.Println(err)
	}
}

func groupEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	group, err := repositories.FetchGroup(groupName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	teams, err := repositories.FetchTeams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pools, err := repositories.FetchPools()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		Group repositories.Group
		Teams []repositories.Team
		Pools []repositories.Pool
	}{
		*group,
		teams,
		pools,
	}
	err = render(w, "templates/groups/edit.html", context)
	if err != nil {
		log.Println(err)
	}
}
