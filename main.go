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

type Page struct {
	Athletes []string
	AllPoints []*analysis.Points
}

var page = Page{}

func initRaces() {

	page.Athletes = []string{"Brain","Townshend","Haaijer","Stojsic","Macuga","Grossniklaus", "Hunt", "Jensen", "Combs", "Robertson", "Hooper", "Tanner"}
	
	raceResults = make ([]racedata.RaceResult,0,20)
	
	for _, race := range races {

		raceResults = append(raceResults,racedata.GetRace(race))

//		ageGroupResults := analysis.SingleRaceAnalysis(raceResults[r])
	}

	page.AllPoints = analysis.PointsAnalysis(raceResults, "U16")
}


func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, &page )
}


func main() {

	initRaces()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}
