package analysis

import (
	"fmt"
	"sort"

	"github.com/saine1a/imd-races/racedata"
)

var points = []int{100, 80, 60, 50, 45, 40, 36, 32, 29, 26, 24, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

var u14Points = []int{200, 175, 155, 140, 128, 118, 110, 103, 96, 90, 84, 79, 74, 70, 66, 63, 60, 57, 54, 51, 49, 47, 45, 43, 41, 39, 37, 35, 33, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

func SingleRaceAnalysis(results racedata.RaceResult) map[string]racedata.ResultArray {

	ageGroupMap := make(map[string]racedata.ResultArray)

	for _, v := range results.Results {

		ageGroupMap[v.Age] = append(ageGroupMap[v.Age], v)
	}

	for k := range ageGroupMap {
		sort.Sort(ageGroupMap[k])

		resultArray := ageGroupMap[k]

		for i, v := range resultArray {

			if resultArray[i].Club != "SBSTA" && resultArray[i].Club != "ASC" && resultArray[i].Club != "SWA" {
				if resultArray[i].Dnf == true {
					resultArray[i].Position = 999
					resultArray[i].AgePosition = 999
				} else {
					if i > 0 && racedata.TotalTime(resultArray[i-1]) == racedata.TotalTime(v) { // Exact same time
						resultArray[i].Position = resultArray[i-1].Position
						resultArray[i].AgePosition = resultArray[i-1].AgePosition
					} else {
						resultArray[i].AgePosition = i + 1
						resultArray[i].AgePosition = i + 1
					}
				}

				if v.Age == "U14" || v.Age == "U16" {
					if resultArray[i].AgePosition <= 60 {
						resultArray[i].Points = u14Points[resultArray[i].AgePosition-1]
					} else {
						resultArray[i].Points = 0
					}
				} else {
					if resultArray[i].AgePosition <= 30 {
						resultArray[i].Points = points[resultArray[i].AgePosition-1]
					} else {
						resultArray[i].Points = 0
					}
				}

			}
		}
	}

	return ageGroupMap
}

type Points struct {
	Athlete       string
	Club          string
	BirthYear     string
	GSPoints      []int
	SLPoints      []int
	SGPoints      []int
	GSPointTotal  int
	SLPointTotal  int
	SGPointTotal  int
	OverallPoints int
	OverallRank   int
	GSResults     []*racedata.Result
	SLResults     []*racedata.Result
	SGResults     []*racedata.Result
}

type PointsArray []*Points

func (a PointsArray) Len() int {
	return len(a)
}

func (a PointsArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a PointsArray) Less(i, j int) bool {

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

	sort.Sort(sort.Reverse(sort.IntSlice(p.SGPoints)))

	sg := 0

	for i, pts := range p.SGPoints {
		if i < 2 {
			sg += pts
		}
	}

	p.SGPointTotal = sg

	p.OverallPoints = gs + sl + sg
}

func PointsAnalysis(races []racedata.RaceResult, ageGroup string) []*Points {

	var athletePoints = make(map[string]*Points)

	for _, r := range races {

		ageGroupResults := SingleRaceAnalysis(r)

		ageSpecificResults := ageGroupResults[ageGroup]

		for _, v := range ageSpecificResults {

			if athletePoints[v.Athlete] == nil {
				athletePoints[v.Athlete] = &Points{}
			}

			athletePoints[v.Athlete].Athlete = v.Athlete
			athletePoints[v.Athlete].Club = v.Club
			athletePoints[v.Athlete].BirthYear = v.BirthYear

			if r.Definition.Qualifier {
				if r.RaceType == "Slalom" {
					athletePoints[v.Athlete].SLPoints = append(athletePoints[v.Athlete].SLPoints, v.Points)
					CalculatePoints(athletePoints[v.Athlete])
					athletePoints[v.Athlete].SLResults = append(athletePoints[v.Athlete].SLResults, v)
				} else {
					if r.RaceType == "Giant Slalom" {
						athletePoints[v.Athlete].GSPoints = append(athletePoints[v.Athlete].GSPoints, v.Points)
						CalculatePoints(athletePoints[v.Athlete])
						athletePoints[v.Athlete].GSResults = append(athletePoints[v.Athlete].GSResults, v)
					} else {
						if r.RaceType == "Super-G" || r.RaceType == "Super G" || r.RaceType == "Speed Training" {
							athletePoints[v.Athlete].SGPoints = append(athletePoints[v.Athlete].SGPoints, v.Points)
							CalculatePoints(athletePoints[v.Athlete])
							athletePoints[v.Athlete].SGResults = append(athletePoints[v.Athlete].SGResults, v)
						} else {
							fmt.Println("UNKNOWN RACE TYPE " + r.RaceType)
						}
					}
				}
			}

		}
	}

	points := make(PointsArray, 0, 100)

	for _, v := range athletePoints {

		points = append(points, v)

	}

	sort.Sort(sort.Reverse(points))

	for i, a := range points {

		if i > 0 && a.OverallPoints == points[i-1].OverallPoints {
			a.OverallRank = points[i-1].OverallRank
		} else {
			a.OverallRank = i + 1
		}

	}

	return points
}
