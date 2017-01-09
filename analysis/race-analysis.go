package analysis

import (
	"imd-races/racedata"
	"sort"
)

var points = []int{100,80,60,50,45,40,36,32,29,26,24,22,20,18,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1}

func SingleRaceAnalysis(results racedata.RaceResult) map[string]racedata.ResultArray {

	ageGroupMap := make(map[string]racedata.ResultArray)
	
	for _, v := range results.Results {

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

				if resultArray[i].AgePosition <= 30 {
					resultArray[i].Points = points[resultArray[i].AgePosition-1]
				} else {
					resultArray[i].Points = 0 
				}
			}
		}
	}

	return ageGroupMap
}

type Points struct {

	Athlete string
	Club string
	GSPoints []int
	SLPoints []int
	GSPointTotal int
	SLPointTotal int
	OverallPoints int
	OverallRank int
}

type PointsArray [] *Points

func ( a PointsArray) Len() int {
	return len(a)
}

func ( a PointsArray) Swap (i, j int) {
	a[i], a[j] = a[j], a[i]
}

func ( a PointsArray) Less (i, j int) bool {

	return a[i].OverallPoints < a[j].OverallPoints
}

func CalculatePoints(p *Points) {

	sort.Sort(sort.Reverse(sort.IntSlice(p.SLPoints)))

	sl := 0
	
	for i, pts := range p.SLPoints {
		if i < 2 {
			sl += pts
		}
	}

	p.SLPointTotal = sl

	sort.Sort(sort.Reverse(sort.IntSlice(p.GSPoints)))

	gs := 0

	for i, pts := range p.GSPoints {
		if i < 2 {
			gs += pts
		}
	}

	p.GSPointTotal = gs

	p.OverallPoints = gs + sl
}

func PointsAnalysis ( races [] racedata.RaceResult, ageGroup string ) []*Points {

	var athletePoints = make(map[string]*Points)
	
	for _, r := range races {

		ageGroupResults := SingleRaceAnalysis(r)

		ageSpecificResults := ageGroupResults[ageGroup]

		for _, v := range ageSpecificResults {

			if athletePoints[v.Ussa] == nil {
				athletePoints[v.Ussa] = &Points{}
			}
			
			athletePoints[v.Ussa].Athlete = v.Athlete
			athletePoints[v.Ussa].Club = v.Club

			if r.RaceType == "Slalom" {
				athletePoints[v.Ussa].SLPoints = append(athletePoints[v.Ussa].SLPoints,v.Points)
				CalculatePoints(athletePoints[v.Ussa])
			} else {
				if r.RaceType == "Giant Slalom" {
					athletePoints[v.Ussa].GSPoints = append(athletePoints[v.Ussa].GSPoints,v.Points)
					CalculatePoints(athletePoints[v.Ussa])
				}
			}

			
		}
	}

	points := make(PointsArray, 0, 100)

	for _, v := range athletePoints {

		points = append(points,v)
		
	}

	sort.Sort(sort.Reverse(points))

	for i, a := range points {

		if i > 0 && a.OverallPoints == points[i-1].OverallPoints {
				a.OverallRank = points[i-1].OverallRank
		} else {
			a.OverallRank = i+1
		}

	}
	
	return points
}
