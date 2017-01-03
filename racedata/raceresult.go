package racedata

import (
	"fmt"
	"strconv"
)

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
}

type ResultArray []Result

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
			result = (a[i].R1 + a[i].R2) < (a[j].R1 + a[j].R2)
		} else {
			result = true // by convention, keep original order if both DNF
		}
	} else {
		result = a[j].Dnf
	}

	fmt.Println("Compare " + a[i].Athlete + " (" + strconv.FormatFloat(a[i].R1 + a[i].R2,'f',-1,64) + ") and " + a[j].Athlete + " (" + strconv.FormatFloat(a[j].R1 + a[j].R2,'f',-1,64) + "), result ")

	if result {
		fmt.Println("i < j")
	} else {
		fmt.Println("i >= j")
	}
	return result
}
