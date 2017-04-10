// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/tsuru/tsuru/db/storage"
)

const (
	// DefaultDatabaseURL represents the default database url
	DefaultDatabaseURL = "127.0.0.1:27017"
	// DefaultDatabaseName represents the default database name
	DefaultDatabaseName = "tsuru_usage"
)

type Storage struct {
	*storage.Storage
}

// conn reads the tsuru-autoscale config and calls storage.Open to get a database connection.
func conn() (*storage.Storage, error) {
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		url = DefaultDatabaseURL
	}
	dbname := os.Getenv("MONGODB_DATABASE_NAME")
	if dbname == "" {
		dbname = DefaultDatabaseName
	}
	return storage.Open(url, dbname)
}

// Conn creates a database connection
func Conn() (*Storage, error) {
	var (
		strg Storage
		err  error
	)
	strg.Storage, err = conn()
	return &strg, err
}

func (s *Storage) TeamGroups() *storage.Collection {
	c := s.Collection("team_groups")
	c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	return c
}

func (s *Storage) PlanCosts() *storage.Collection {
	c := s.Collection("plan_costs")
	c.EnsureIndex(mgo.Index{Key: []string{"plan", "service", "type"}, Unique: true})
	c.EnsureIndexKey("type")
	c.EnsureIndexKey("service")
	return c
}
