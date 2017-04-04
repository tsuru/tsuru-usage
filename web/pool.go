// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import "net/http"

func poolListHandler(w http.ResponseWriter, r *http.Request) {
	pools := []string{"staging", "prod", "workshop"}
	render(w, "web/templates/pool/index.html", pools)
}
