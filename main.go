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

var keyAthletes = []string{"Brain","Townshend","Haaijer","Stojsic","Macuga","Grossniklaus", "Hunt", "Jensen", "Combs", "Robertson", "Hooper", "Tanner"}

func initRaces() {
	
	raceResults = make ([]racedata.RaceResult,0,20)
	
	for _, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))

//		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[r])
	}

	allPoints = analysis.PointsAnalysis(raceResults, "U16")
}


type HomePage struct {
	Athletes []string
	AllPoints []*analysis.Points
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, &HomePage{Athletes: keyAthletes,AllPoints:allPoints} )
}

type RacePage struct {
	Athletes []string
	RaceId string
}

func handleRace(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("race.html")
	raceId := r.URL.Query().Get("raceId")
	t.Execute(w, &RacePage{Athletes:keyAthletes, RaceId : raceId} )
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

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/race", handleRace)	
	http.HandleFunc("/races", handleRaceList)
	http.ListenAndServe(":8080", nil)

}
