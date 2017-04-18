// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
)

type Group struct {
	Name  string
	Teams []string
	Pools []string
}

func fetchGroup(name string) (*Group, error) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teamgroups/%s", host, name)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		return nil, fmt.Errorf("Error fetching %s: %s", url, err)
	}
	defer response.Body.Close()
	var group Group
	err = json.NewDecoder(response.Body).Decode(&group)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: %s", err)
	}
	return &group, nil
}

func fetchGroups() ([]Group, error) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teamgroups", host)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		return nil, fmt.Errorf("Error fetching %s: %s", url, err)
	}
	defer response.Body.Close()
	var groups []Group
	err = json.NewDecoder(response.Body).Decode(&groups)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: %s", err)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups, nil
}
