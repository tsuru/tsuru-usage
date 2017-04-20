// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"

	"github.com/tsuru/tsuru-usage/repositories"
)

func poolListHandler(w http.ResponseWriter, r *http.Request) {
	pools, err := repositories.FetchPools()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render(w, "templates/list/pools.html", pools)
}

func teamListHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := repositories.FetchTeams()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render(w, "templates/list/teams.html", teams)
}

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := repositories.FetchGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render(w, "templates/list/groups.html", groups)
}
