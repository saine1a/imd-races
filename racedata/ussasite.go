package racedata

import (
	"encoding/csv"
	"fmt"
	"imd-races/csvmapper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
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

func GetUSSAResults(definition RaceDefinition) RaceResult {

	raceResult := RaceResult{}

	raceResult.Definition = definition

	// first get race definition info

	url, err := url.Parse("https://my.ussa.org/ussa-tools/events/results/U0173/2018")

	if err != nil {
		log.Fatal(err)
	}

	query := url.Query()

	query.Set("csv", "0")
	url.RawQuery = query.Encode()

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
	fmt.Println(len(record))
	raceInfoHeader := csvmapper.MapHeader(record, reflect.TypeOf(USSARace{}))

	// Now read money line
	record, err = csvReader.Read() // money line
	fmt.Println(len(record))
	fmt.Println(csvReader.FieldsPerRecord)
	raceInfo := USSARace{}

	for _, r := range record {
		fmt.Println(r)
	}

	raceInfoHeader.ParseRecord(record, &raceInfo, OurDateFormat)

	if err == nil {
		raceResult.RaceName = raceInfo.Name
		raceResult.RaceType = raceInfo.Discipline
		raceResult.RaceDate = raceInfo.Date
	} else {
		log.Fatal(err)
	}
	/*
		// Now get results

		url, err = url.Parse("https://my.ussa.org/ussa-tools/events/results/U0173/2018")

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

		record, err = csvReader.Read() // first line
		record, err = csvReader.Read() // second line

		if err == nil {
			raceResult.RaceName = record[1]
			raceResult.RaceType = record[3]
			raceResult.RaceDate = record[2]
		} else {
			log.Fatal(err)
		}
	*/
	return raceResult
}
