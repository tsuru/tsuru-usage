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
	status := r.FormValue("status")
	if status != "updated" {
		status = ""
	}
	context := struct {
		Groups []repositories.Group
		Status string
	}{
		groups,
		status,
	}
	err = render(w, "templates/groups/index.html", context)
	if err != nil {
		log.Println(err)
	}
}

func groupNewHandler(w http.ResponseWriter, r *http.Request) {
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
		Teams []repositories.Team
		Pools []repositories.Pool
	}{
		teams,
		pools,
	}
	err = render(w, "templates/groups/new.html", context)
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

func groupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	group := repositories.Group{
		Name:  vars["name"],
		Teams: r.Form["teams"],
		Pools: r.Form["pools"],
	}
	err = repositories.UpdateGroup(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/teamgroups?status=updated", 302)
}
