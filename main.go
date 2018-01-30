package main

import (
	"html/template"
	"imd-races/analysis"
	"imd-races/racedata"
	"net/http"
)

//var races = []racedata.RaceDefinition{{RaceId: "178531", Qualifier: true}, {RaceId: "178336", Qualifier: true}, {RaceId: "178251", Qualifier: true}}
var races = []racedata.RaceDefinition{{RaceId: "181356", Qualifier: true}, {RaceId: "181519", Qualifier: true}}

//var focusAthlete = "I6465959"
var focusAthlete = "X6466759"

//var ageGroup = "U14"
var ageGroup = "U16"

var raceResults []racedata.RaceResult

var allPoints []*analysis.Points

func initRaces() {

	raceResults = make([]racedata.RaceResult, 0, 20)

	for _, race := range races {
		raceResults = append(raceResults, racedata.GetRace(race))

	}

	allPoints = analysis.PointsAnalysis(raceResults, ageGroup)
}

type HomePage struct {
	AllPoints    []*analysis.Points
	FocusAthlete string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, &HomePage{AllPoints: allPoints, FocusAthlete: focusAthlete})
}

type RacePage struct {
	RaceResult   *racedata.RaceResult
	FocusAthlete string
	Result       racedata.ResultArray
}

func handleRace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")

	raceIndex := -1
	raceResult := &racedata.RaceResult{}

	for i, race := range raceResults {
		if race.Definition.RaceId == raceId {
			raceIndex = i
			raceResult = &(raceResults[i])
		}
	}

	if raceIndex >= 0 {
		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[raceIndex])[ageGroup]
		t.Execute(w, &RacePage{RaceResult: raceResult, FocusAthlete: focusAthlete, Result: ageGroupResults})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

type RaceListPage struct {
	RaceResults []racedata.RaceResult
}

func handleRaceList(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("raceList.html")
	t.Execute(w, &RaceListPage{RaceResults: raceResults})
}

type AthletePage struct {
	AthleteName    string
	Ussa           string
	RaceResults    []racedata.RaceResult
	AthleteResults []*racedata.Result
}

func handleAthlete(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("athlete.html")
	athleteName := r.URL.Query().Get("name")
	ussa := r.URL.Query().Get("ussa")

	athleteResults := make([]*racedata.Result, 0)

	for _, race := range raceResults {

		athleteFound := false

		for _, result := range race.Results {

			if result.Ussa == ussa {
				athleteResults = append(athleteResults, result)
				athleteFound = true
			}
		}

		if athleteFound == false {
			athleteResults = append(athleteResults, &racedata.Result{DnfReason: "DNS"})
		}
	}

	t.Execute(w, &AthletePage{RaceResults: raceResults, Ussa: ussa, AthleteName: athleteName, AthleteResults: athleteResults})
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
