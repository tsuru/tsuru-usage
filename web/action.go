// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func render(w http.ResponseWriter, templatePath string, data interface{}) error {
	dir, _ := os.Getwd()
	if !strings.HasSuffix(dir, "/web") {
		dir += "/web"
	}
	templates := []string{
		dir + "/" + templatePath,
		dir + "/templates/base.html",
		dir + "/templates/back.html",
		dir + "/templates/group_tabs.html",
	}
	t, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("Error parsing template %s: %s", templatePath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error rendering template %s: %s", templatePath, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err
}
