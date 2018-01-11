package racedata

import (
	"fmt"
)

type Dnf struct {
	Ussa   string
	RaceId string
}

var dnfOverrides = []Dnf{{RaceId: "178251", Ussa: "I6779097"}}

func IsDnf(result *Result, raceId string) bool {

	if result.Dnf == true {
		return true
	}

	for _, v := range dnfOverrides {
		if v.RaceId == raceId && v.Ussa == result.Ussa {
			fmt.Println("DNF Override for %s %s\n", result.Athlete, raceId)
			return true
		}
	}

	return false
}
