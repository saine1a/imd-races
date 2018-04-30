package racedata

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/saine1a/imd-races/csvmapper"
	"github.com/saine1a/imd-races/racelisting"
)

var OurDateFormat = time.RFC3339 // "2006-1-2 15:4:5"

type USSARace struct {
	Code       string
	Name       string
	Date       string
	Discipline string
	Gender     string
	Vertical   int64
	Location   string
	Sport      string
	Division   string
}

type USSAResult struct {
	FinishPlace string
	AthleteID   string
	Athlete     string `csvField:"Full Name"`
	BirthYear   int64
	Division    string `csvField:"Division/Country"`
	R1          string `csvField:"FirstRun"`
	R2          string `csvField:"SecondRun"`
	RaceTime    string
	RacePoints  float64
	USSAResult  float64
}

func GetUSSAResults(definition racelisting.RaceDefinition) RaceResult {

	raceResult := RaceResult{}

	raceResult.Definition = definition

	// first get race definition info

	urlString := fmt.Sprintf("https://my.ussa.org/ussa-tools/events/results/%s/2018", definition.RaceId)

	url, err := url.Parse(urlString)

	if err != nil {
		log.Fatal(err)
	}

	query := url.Query()

	query.Set("csv", "0")
	url.RawQuery = query.Encode()

	fmt.Printf("Getting URL %s\n", url.String())

	resp, err := client.Get(url.String())

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	rawReader := strings.NewReader(string(body))

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(rawReader)

	// Skip first line
	record, err := csvReader.Read()
	csvReader.FieldsPerRecord = 0 // Hack to reset csvReader

	// Map CSV headers
	record, err = csvReader.Read() // header line
	raceInfoHeader := csvmapper.MapHeader(record, reflect.TypeOf(USSARace{}))

	// Now read money line
	record, err = csvReader.Read() // money line
	raceInfo := USSARace{}

	raceInfoHeader.ParseRecord(record, &raceInfo, OurDateFormat)

	if err == nil {
		raceResult.RaceName = raceInfo.Name
		raceResult.RaceType = raceInfo.Discipline
		raceResult.RaceDate = raceInfo.Date
	} else {
		log.Fatal(err)
	}

	// Now get results

	url, err = url.Parse(urlString)

	if err != nil {
		log.Fatal(err)
	}

	query = url.Query()

	query.Set("csv", "2")
	url.RawQuery = query.Encode()

	resp, err = client.Get(url.String())

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	csvReader = csv.NewReader(strings.NewReader(string(body)))

	// Skip first line
	record, err = csvReader.Read()
	csvReader.FieldsPerRecord = 0 // Hack to reset csvReader

	record, err = csvReader.Read() // second line = headers
	raceResultHeader := csvmapper.MapHeader(record, reflect.TypeOf(USSAResult{}))

	eof := false

	raceResult.Results = make(ResultArray, 0, 200)

	for !eof {

		record, err = csvReader.Read()

		if err != nil || len(record) == 0 {
			eof = true
		} else {
			result := USSAResult{}
			raceResultHeader.ParseRecord(record, &result, OurDateFormat)
			modifiedResult := Result{}

			modifiedResult.Age = calcAge(result.BirthYear)
			modifiedResult.BirthYear = fmt.Sprintf("%d", result.BirthYear)
			modifiedResult.Athlete = result.Athlete
			modifiedResult.Ussa = result.AthleteID
			modifiedResult.Bib = "-"
			modifiedResult.Club = "-"
			modifiedResult.Dnf = (strings.HasPrefix(result.FinishPlace, "DNF"))
			if modifiedResult.Dnf {
				modifiedResult.Position = 999
				modifiedResult.DnfReason = result.FinishPlace
			} else {
				modifiedResult.Position, _ = strconv.Atoi(result.FinishPlace)
			}
			modifiedResult.R1 = convTime(result.R1)
			modifiedResult.R2 = convTime(result.R2)
			modifiedResult.RaceType = raceResult.RaceType
			modifiedResult.USSAPoints = result.USSAResult
			raceResult.Results = append(raceResult.Results, &modifiedResult)
		}
	}

	raceResult.Results.SortResults()

	return raceResult
}

func calcAge(year int64) string {

	switch year {
	case 2002:
		return "U16"
	case 2003:
		return "U16"
	case 2004:
		return "U14"
	case 2005:
		return "U14"
	default:
		return "-"

	}
}

func convTime(raceTime string) float64 {

	// Special case for wierd bug with USSA site

	if strings.HasPrefix(raceTime, "00:60.") {
		raceTime = strings.Replace(raceTime, "00:60.", "01:00.", 1)
	}

	theTime, err :=
		time.Parse("04:05.00", raceTime)

	if err != nil {
		return 999
	} else {
		value := float64(theTime.Minute())*60.0 + float64(theTime.Second()) + float64(theTime.Nanosecond())/1000000000.0
		return value
	}
}
