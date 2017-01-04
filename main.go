package main

import (
	"imd-races/racedata"
	"imd-races/analysis"
	"fmt"
	"strings"
)


func main() {

	races := []string{ "164315", "164466", "164383", "164431" }
	
	for _, race := range races {
		results := racedata.GetRace(race)

		ageGroupResults := analysis.RaceAnalysis(results.Results)
	
		u16Results := ageGroupResults["U16"]

		for _, v := range u16Results {

			if strings.Contains(v.Athlete,"Brain") {			
				if v.Dnf == false {
					fmt.Printf("%s %s %d : %s\n", results.RaceName, results.RaceType, v.AgePosition, v.Athlete)
				} else {
					fmt.Printf("%s : %s\n", v.DnfReason, v.Athlete)
				}
			}
		}
	}
}
