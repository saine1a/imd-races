package main

import (
	"html/template"
	"imd-races/analysis"
	"imd-races/racedata"
	"net/http"
)

var focusAthlete = "Brain, Jonathan"
var ageGroup = "U14"

//var races = []racedata.RaceDefinition{{RaceId: "178531", Qualifier: true, TimingSite: racedata.LIVETIMING}, {RaceId: "178336", Qualifier: true, TimingSite: racedata.LIVETIMING}, {RaceId: "178251", Qualifier: true, TimingSite: racedata.LIVETIMING}}
var races = []racedata.RaceDefinition{
	{RaceId: "U0637", Qualifier: true, TimingSite: racedata.USSA}, 
	{RaceId: "U0638", Qualifier: true, TimingSite: racedata.USSA}, 
	{RaceId: "U0634", Qualifier: true, TimingSite: racedata.USSA},
	{RaceId: "U0051", Qualifier: true, TimingSite: racedata.USSA},
	{RaceId: "U0053", Qualifier: true, TimingSite: racedata.USSA}}

/* REBEKAH

var focusAthlete = "Brain, Rebekah"

var ageGroup = "U16"

var races = []racedata.RaceDefinition{
	{RaceId: "U0173", Qualifier: true, TimingSite: racedata.USSA},
	{RaceId: "U0175", Qualifier: true, TimingSite: racedata.USSA},
	{RaceId: "U0167", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0165", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0171", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0169", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0643", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0645", Qualifier: false, TimingSite: racedata.USSA},
	{RaceId: "U0647", Qualifier: false, TimingSite: racedata.USSA},
}*/

//var races = []racedata.RaceDefinition{{RaceId: "U0173", Qualifier: true, TimingSite: racedata.USSA}}

var raceResults []racedata.RaceResult

var allPoints []*analysis.Points

func initRaces() {

	raceResults = make([]racedata.RaceResult, 0, 20)

	for _, race := range races {
		switch race.TimingSite {
		case racedata.LIVETIMING:
			raceResults = append(raceResults, racedata.GetLiveTimingResults(race))
		case racedata.USSA:
			raceResults = append(raceResults, racedata.GetUSSAResults(race))
		}

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

	athleteResults := make([]*racedata.Result, 0)

	for _, race := range raceResults {

		athleteFound := false

		for _, result := range race.Results {

			if result.Athlete == athleteName {
				athleteResults = append(athleteResults, result)
				athleteFound = true
			}
		}

		if athleteFound == false {
			athleteResults = append(athleteResults, &racedata.Result{DnfReason: "DNS"})
		}
	}

	t.Execute(w, &AthletePage{RaceResults: raceResults, AthleteName: athleteName, AthleteResults: athleteResults})
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
