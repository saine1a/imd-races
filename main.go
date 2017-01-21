package main

import (
	"imd-races/racedata"
	"imd-races/analysis"
	"net/http"
	"html/template"
	"fmt"
)

var races = []string{ "165433", "165551", "165733"}

var raceResults [] racedata.RaceResult

var allPoints []*analysis.Points

var focusAthlete = "X6466759"

var ageGroup = "U16"

func initRaces() {
	
	raceResults = make ([]racedata.RaceResult,0,20)
	
	for _, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))


	}

	allPoints = analysis.PointsAnalysis(raceResults, ageGroup)
}


type HomePage struct {
	AllPoints []*analysis.Points
	FocusAthlete string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, &HomePage{AllPoints:allPoints, FocusAthlete:focusAthlete} )
}

type RacePage struct {
	FocusAthlete string
	RaceType string
	Result racedata.ResultArray
}

func handleRace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")

	raceIndex := -1
	raceType := ""
	
	for i, race := range raceResults {
		if race.RaceId == raceId {
			raceIndex = i
			raceType = race.RaceType
		}
	}

	fmt.Println("Hello world")
	fmt.Println(raceType)
	
	if raceIndex >= 0 {
		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[raceIndex])[ageGroup]
		t.Execute(w, &RacePage{RaceType:raceType, FocusAthlete:focusAthlete, Result:ageGroupResults} )
	}
}

type RaceListPage struct {
	RaceResults [] racedata.RaceResult
}

func handleRaceList(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("raceList.html")
	t.Execute(w, &RaceListPage{RaceResults:raceResults} )
}

type AthletePage struct {
	Athlete string
	Points *analysis.Points
}

func handleAthlete(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("athlete.html")
	athleteName := r.URL.Query().Get("athlete")

	athleteIndex := -1
	
	for i, athlete := range allPoints {
		if athlete.Ussa == focusAthlete {
			athleteIndex = i
		}
	}

	if athleteIndex >= 0 {
		t.Execute(w, &AthletePage{Athlete:athleteName,Points:allPoints[athleteIndex]})
	}
}

func main() {

	initRaces()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources")))) 

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/athlete", handleAthlete)
	http.HandleFunc("/race", handleRace)	
	http.HandleFunc("/races", handleRaceList)
	http.ListenAndServe(":8080", nil)

}
