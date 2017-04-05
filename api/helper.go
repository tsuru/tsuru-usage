package api

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/common/model"
)

func runForYear(f func(month time.Month) (model.Vector, error)) map[time.Month]model.Vector {
	results := make(chan monthUsage, 12)
	usage := make(map[time.Month]model.Vector)
	wg := sync.WaitGroup{}
	wg.Add(12)
	for m := 1; m <= 12; m++ {
		go func(m int) {
			result, err := f(time.Month(m))
			if err != nil {
				log.Printf("failed to get month %s usage: %s", time.Month(m).String(), err)
			}
			results <- monthUsage{month: time.Month(m), value: result}
			wg.Done()
		}(m)
	}
	wg.Wait()
	close(results)
	for u := range results {
		usage[u.month] = u.value
	}
	return usage
}

func selectorForGroup(groupName string) (string, error) {
	group, err := FindTeamGroup(groupName)
	if err != nil {
		return "", err
	}
	return strings.Join(group.Teams, "|"), nil
}
