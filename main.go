package main

import (
	"imd-races/racedata"
	"imd-races/analysis"
	"fmt"
)

func main() {

	results := racedata.GetRace()

	ageGroupResults := analysis.RaceAnalysis(results)
	
	u16Results := ageGroupResults["U16"]

	for i, v := range u16Results {
		if v.Dnf == false {
			fmt.Printf("%d : %s\n", i, v.Athlete)
		} else {
			fmt.Printf("%s : %s\n", v.DnfReason, v.Athlete)
		}
	}
}
