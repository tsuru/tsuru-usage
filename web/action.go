// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, templatePath string, data interface{}) error {
	t, err := template.ParseFiles(templatePath, "web/templates/base.html")
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "base", data)
}
