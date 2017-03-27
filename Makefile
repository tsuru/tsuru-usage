# Copyright 2017 tsuru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build:
	go build -o bin/tsuru-usage

test:
	go test ./...
