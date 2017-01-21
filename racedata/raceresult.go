package racedata


type Result struct {

	Bib string
	Athlete string
	Club string
	Age string
	Ussa string
	R1 float64
	R2 float64
	Dnf bool
	DnfReason string
	AgePosition int
	Points int
	RaceType string
}

func TotalTime(result Result) float64 {
	if result.RaceType == "Super-G" {
		return result.R1
	} else {
		return result.R1 + result.R2
	}
}

type ResultArray []Result

type RaceResult struct {
	RaceName string
	RaceId string
	RaceType string
	Results ResultArray
}

func ( a ResultArray) Len() int {
	return len(a)
}

func ( a ResultArray) Swap (i, j int) {
	a[i], a[j] = a[j], a[i]
}

func ( a ResultArray) Less (i, j int) bool {

	var result bool

	if a[i].Dnf == false {

		if a[j].Dnf == false {
			result = TotalTime(a[i]) < TotalTime(a[j])
		} else {
			result = true // by convention, keep original order if both DNF
		}
	} else {
		result = a[j].Dnf
	}

	return result
}
