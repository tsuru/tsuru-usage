// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)

func render(w http.ResponseWriter, templatePath string, data interface{}) error {
	dir, _ := os.Getwd()
	if !strings.HasSuffix(dir, "/web") {
		dir += "/web"
	}
	t, err := template.ParseFiles(dir+"/"+templatePath, dir+"/templates/base.html")
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "base", data)
}
