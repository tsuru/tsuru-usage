// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type Group struct {
	Name  string
	Teams []string
	Pools []string
}

func FetchGroup(name string) (*Group, error) {
	url := fmt.Sprintf("%s/api/teamgroups/%s", apiHost, name)
	response, err := Client.Get(url)
	status := response.StatusCode
	if err != nil || (status != http.StatusOK && status != http.StatusNotFound) {
		errMsg := fmt.Sprintf("Error fetching %s: %s", url, err.Error())
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	defer response.Body.Close()
	var group Group
	err = json.NewDecoder(response.Body).Decode(&group)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding response body: %s", err.Error())
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	return &group, nil
}

func FetchGroups() ([]Group, error) {
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/teamgroups", host)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Error fetching %s: %s", url, err.Error())
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	defer response.Body.Close()
	var groups []Group
	err = json.NewDecoder(response.Body).Decode(&groups)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding response body: %s", err.Error())
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups, nil
}

func UpdateGroup(group Group) error {
	addr := fmt.Sprintf("%s/api/teamgroups/%s", apiHost, group.Name)
	v := url.Values{"teams": group.Teams, "pools": group.Pools}
	req, err := http.NewRequest(http.MethodPut, addr, strings.NewReader(v.Encode()))
	if err != nil {
		errMsg := fmt.Sprintf("Error in PUT %s: %s", addr, err.Error())
		log.Printf(errMsg)
		return errors.New(errMsg)
	}
	response, err := Client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Error in PUT %s: %s", addr, err.Error())
		log.Printf(errMsg)
		return errors.New(errMsg)
	}
	status := response.StatusCode
	if status != http.StatusOK && status != http.StatusCreated {
		errMsg := fmt.Sprintf("Error in PUT %s: HTTP %d", addr, status)
		log.Printf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}
