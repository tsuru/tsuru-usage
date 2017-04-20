// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin

import (
	"log"
	"net/http"

	"github.com/tsuru/tsuru-usage/repositories"
)

func groupListHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := repositories.FetchGroups()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		Groups []repositories.Group
	}{
		groups,
	}
	err = render(w, "templates/index.html", context)
	if err != nil {
		log.Println(err)
	}
}
