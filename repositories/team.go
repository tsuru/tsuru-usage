// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
)

type Team struct {
	Name string
}

func FetchTeams() ([]Team, error) {
	url := fmt.Sprintf("%s/api/teams", apiHost)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		return nil, fmt.Errorf("Error fetching %s: %s", url, err)
	}
	defer response.Body.Close()
	var teams []Team
	err = json.NewDecoder(response.Body).Decode(&teams)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: %s", err)
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams, nil
}
