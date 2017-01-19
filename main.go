package main

import (
	"imd-races/racedata"
	"imd-races/analysis"
	"net/http"
	"html/template"
)

var races = []string{ "165433", "165551", "165733"}

var raceResults [] racedata.RaceResult

var allPoints []*analysis.Points

var focusAthlete = "X6466759"

func initRaces() {
	
	raceResults = make ([]racedata.RaceResult,0,20)
	
	for _, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))


	}

	allPoints = analysis.PointsAnalysis(raceResults, "U16")
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
	Result racedata.ResultArray
}

func handleRace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")

	raceIndex := -1
	
	for i, race := range raceResults {
		if race.RaceId == raceId {
			raceIndex = i
		}
	}

	if raceIndex >= 0 {
		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[raceIndex])["U16"]
		t.Execute(w, &RacePage{FocusAthlete:focusAthlete, Result:ageGroupResults} )
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
}

func handleAthlete(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("athlete.html")
	athleteName := r.URL.Query().Get("athlete")

	t.Execute(w, &AthletePage{Athlete:athleteName})
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
