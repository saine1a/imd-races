package analysis

import (
	"imd-races/racedata"
	"sort"
)



func RaceAnalysis(results [] racedata.Result) map[string]racedata.ResultArray {

	ageGroupMap := make(map[string]racedata.ResultArray)
	
	for _, v := range results {

		ageGroupMap[v.Age] = append(ageGroupMap[v.Age],v)
	}

	for k := range ageGroupMap {
		sort.Sort(ageGroupMap[k])

		resultArray := ageGroupMap[k]
		
		for i, v := range resultArray {

			if v.Dnf == false {
				if i > 0 && racedata.TotalTime(resultArray[i-1]) == racedata.TotalTime(v) { // Exact same time
					resultArray[i].AgePosition = resultArray[i-1].AgePosition
				} else {
					resultArray[i].AgePosition = i+1
				}
			}
		}
	}

	return ageGroupMap
}
