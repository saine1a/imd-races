package main

import (
	"imd-races/racedata"
	"imd-races/analysis"
	"fmt"
	"strings"
)


func main() {

	races := []string{ "164315", "164466", "164383", "164431" }

	keyAthletes := []string{"Brain","Townshend","Haaijer","Stojsic"}

	raceResults := make ([]racedata.RaceResult,0,20)
	
	for r, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))

		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[r])
	
		u16Results := ageGroupResults["U16"]

		for _, v := range u16Results {

			for _, a := range keyAthletes {
				if strings.Contains(v.Athlete,a) {			
					if v.Dnf == false {
						fmt.Printf("%s %s %d %s %d\n", raceResults[r].RaceName, raceResults[r].RaceType, v.AgePosition, v.Athlete, v.Points)
					} else {
					fmt.Printf("%s : %s %d\n", v.DnfReason, v.Athlete, v.Points)
					}
				}
			}
		}

	}

	allPoints := analysis.PointsAnalysis(raceResults, "U16")

	for _, a := range allPoints {

		fmt.Printf("%d %s %d %d %d\n", a.OverallRank, a.Athlete, a.SLPointTotal, a.GSPointTotal,  a.OverallPoints)
	}

}
