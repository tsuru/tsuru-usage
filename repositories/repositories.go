// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"net/http"
	"os"
)

var Client = &http.Client{}
var apiHost = os.Getenv("API_HOST")
