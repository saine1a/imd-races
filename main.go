package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saine1a/imd-races/analysis"
	"github.com/saine1a/imd-races/racedata"
	"github.com/saine1a/imd-races/racelisting"
)

/*
var focusAthlete = "Brain, Rebekah"

var ageGroup = "U16"
*/

var focusAthlete = "Brain, Jonathan"
var ageGroup = "U14"

var raceResults []racedata.RaceResult

var allPoints []*analysis.Points

func initRaces() {

	raceResults = make([]racedata.RaceResult, 0, 20)

	for _, race := range racelisting.Races {
		switch race.TimingSite {
		case racelisting.LIVETIMING:
			raceResults = append(raceResults, racedata.GetLiveTimingResults(race))
		case racelisting.USSA:
			raceResults = append(raceResults, racedata.GetUSSAResults(race))
		}

	}

	allPoints = analysis.PointsAnalysis(raceResults, ageGroup)
}

type HomePage struct {
	AllPoints    []*analysis.Points
	FocusAthlete string
}

func handleAthleteImpl() *HomePage {
	return &HomePage{AllPoints: allPoints, FocusAthlete: focusAthlete}
}

func handleGoHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")

	t.Execute(w, handleAthleteImpl())
}

type RacePage struct {
	RaceResult   *racedata.RaceResult
	FocusAthlete string
	Result       racedata.ResultArray
}

func handleRaceImpl(raceId string) *RacePage {

	raceIndex := -1
	raceResult := &racedata.RaceResult{}

	for i, race := range raceResults {
		if race.Definition.RaceId == raceId {
			raceIndex = i
			raceResult = &(raceResults[i])
		}
	}

	if raceIndex != -1 {
		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[raceIndex])[ageGroup]

		return &RacePage{RaceResult: raceResult, FocusAthlete: focusAthlete, Result: ageGroupResults}
	} else {
		return nil
	}
}

func handleRaceHtml(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")

	raceInfo := handleRaceImpl(raceId)

	if raceInfo != nil {
		t.Execute(w, raceInfo)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleRaceREST(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	raceId := params["raceId"]

	json.NewEncoder(w).Encode(handleRaceImpl(raceId))
}

type RaceListPage struct {
	RaceResults []racedata.RaceResult
}

func handleRaceListImpl() *RaceListPage {
	return &RaceListPage{RaceResults: raceResults}
}

func handleRaceListHtml(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("raceList.html")
	t.Execute(w, handleRaceListImpl())
}

func handleRaceListREST(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(handleRaceListImpl())
}

type AthletePage struct {
	AthleteName    string
	Ussa           string
	RaceResults    []racedata.RaceResult
	AthleteResults []*racedata.Result
}

func handleAthleteHtml(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("athlete.html")
	athleteName := r.URL.Query().Get("name")

	t.Execute(w, handleIndividualAthleteImpl(athleteName))
}

func handleAthletesREST(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(handleAthleteImpl())

}

func handleIndividualAthleteImpl(athleteName string) *AthletePage {

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

	return &AthletePage{RaceResults: raceResults, AthleteName: athleteName, AthleteResults: athleteResults}
}

func handleIndividualAthleteREST(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	athleteName := params["name"]

	json.NewEncoder(w).Encode(handleIndividualAthleteImpl(athleteName))
}

func main() {

	initRaces()

	router := mux.NewRouter()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	// HTML routes for Elm code

	router.Handle("/", http.FileServer(http.Dir("./Web")))

	// HTML routes for Go templates

	router.HandleFunc("/go", handleGoHome)
	router.HandleFunc("/go/athletePage", handleAthleteHtml)
	router.HandleFunc("/go/racePage", handleRaceHtml)
	router.HandleFunc("/go/racesPage", handleRaceListHtml)

	// REST API routes

	router.HandleFunc("/athlete/{name}", handleIndividualAthleteREST).Methods("GET")
	router.HandleFunc("/athlete", handleAthletesREST).Methods("GET")
	router.HandleFunc("/race/{raceId}", handleRaceREST).Methods("GET")
	router.HandleFunc("/race", handleRaceListREST).Methods("GET")

	fmt.Println("Ready & listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
