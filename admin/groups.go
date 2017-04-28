// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-usage/repositories"
)

type GroupContext struct {
	Teams []repositories.Team
	Pools []repositories.Pool
}

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := repositories.FetchGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updated, _ := strconv.ParseBool(r.FormValue("updated"))
	context := struct {
		Groups  []repositories.Group
		Updated bool
	}{
		groups,
		updated,
	}
	err = render(w, "templates/groups/index.html", context)
	if err != nil {
		log.Println(err)
	}
}

func groupNewHandler(w http.ResponseWriter, r *http.Request) {
	groupContext, err := fetchGroupContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = render(w, "templates/groups/new.html", groupContext)
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
	groupContext, err := fetchGroupContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	context := struct {
		Group repositories.Group
		Teams []repositories.Team
		Pools []repositories.Pool
	}{
		*group,
		groupContext.Teams,
		groupContext.Pools,
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
	http.Redirect(w, r, "/admin/teamgroups?updated=true", 302)
}

func groupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	group := repositories.Group{Name: groupName}
	err := repositories.DeleteGroup(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func fetchGroupContext() (GroupContext, error) {
	context := GroupContext{}
	errs := make(chan error, 2)
	go func() {
		teams, err := repositories.FetchTeams()
		context.Teams = teams
		errs <- err
	}()
	go func() {
		pools, err := repositories.FetchPools()
		context.Pools = pools
		errs <- err
	}()
	for i := 0; i < cap(errs); i++ {
		err := <-errs
		if err != nil {
			return context, err
		}
	}
	return context, nil
}
