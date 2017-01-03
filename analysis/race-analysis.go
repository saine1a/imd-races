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
	}

	return ageGroupMap
}
