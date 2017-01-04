package analysis

import (
	"imd-races/racedata"
	"sort"
)

var points = []int{100,80,60,50,45,40,36,32,29,26,24,22,20,18,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1}

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

				if resultArray[i].AgePosition < 30 {
					resultArray[i].Points = points[resultArray[i].AgePosition-1]
				} else {
					resultArray[i].Points = 0 
				}
			}
		}
	}

	return ageGroupMap
}
