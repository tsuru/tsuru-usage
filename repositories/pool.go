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

type Pool struct {
	Name string
}

func FetchPools() ([]Pool, error) {
	url := fmt.Sprintf("%s/api/pools", apiHost)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		return nil, fmt.Errorf("Error fetching %s: %s", url, err)
	}
	defer response.Body.Close()
	var pools []Pool
	err = json.NewDecoder(response.Body).Decode(&pools)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response body: %s", err)
	}
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Name < pools[j].Name
	})
	return pools, nil
}
