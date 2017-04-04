# Copyright 2017 tsuru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build:
	go build -o bin/tsuru-usage

test:
	go test ./...

coverage:
	./go.test.bash

linux-build:
	GOOS=linux GOARCH=amd64 go build -o tsuru-usage

deploy: linux-build
	tsuru app-deploy -a tsuru-usage tsuru-usage Procfile tsuru.yaml