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

var focusAthlete = "Brain"

func initRaces() {
	
	raceResults = make ([]racedata.RaceResult,0,20)
	
	for _, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))


	}

	allPoints = analysis.PointsAnalysis(raceResults, "U16")
}


type HomePage struct {
	Athlete string
	AllPoints []*analysis.Points
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, &HomePage{Athlete: focusAthlete,AllPoints:allPoints} )
}

type RacePage struct {
	Athlete string
	Result racedata.ResultArray
}

func handleRace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")

	fmt.Println(raceId)
	
	raceIndex := -1
	
	for i, race := range raceResults {
		if race.RaceId == raceId {
			raceIndex = i
		}
	}

	if raceIndex >= 0 {
		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[raceIndex])["U16"]
		t.Execute(w, &RacePage{Athlete:focusAthlete, Result:ageGroupResults} )
	}
}

type RaceListPage struct {
	RaceResults [] racedata.RaceResult
}

func handleRaceList(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("raceList.html")
	t.Execute(w, &RaceListPage{RaceResults:raceResults} )
}

func main() {

	initRaces()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources")))) 

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/race", handleRace)	
	http.HandleFunc("/races", handleRaceList)
	http.ListenAndServe(":8080", nil)

}
