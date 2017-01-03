package racedata

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
	"io/ioutil"
	"bufio"
	"strings"
	"strconv"
)

var client = &http.Client{}


func GetRace() [] Result {

	url, err := url.Parse("http://live-timing.com/includes/aj_race.php")

	if err != nil {
		log.Fatal(err)
	}

	query := url.Query()
	query.Set("r", "164315")
	query.Set("m", "1")
	query.Set("u", "0")
	url.RawQuery = query.Encode()

	fmt.Println(url)
	
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

	fmt.Println(string(body))

	s := bufio.NewScanner(strings.NewReader(string(body)))

	split := func(data []byte, atEOF bool) (int,[]byte,error) {

		var i int

		var b []byte

		var advance int
		var token [] byte
		var err error
		
		for advance, token, err = bufio.ScanRunes(data, atEOF) ; err == nil && token != nil && string(token) != "|" && (i + advance ) < len(data) ; {

			b = append(b, token...)

			i += advance

			if ( i <= len(data) ) {
				advance, token, err = bufio.ScanRunes(data[i:len(data)], atEOF)
			}
		}

		return i+advance, b, nil
	}

	s.Split(split)
	
	for ; s.Scan() && s.Text() != "hE" ; {
		// Scan until we find "hE"

//		fmt.Println("Skipping " + s.Text())
	}

	results := make([]Result, 0, 200)

	var i = 0

	for ; s.Scan() && s.Text() != "endC" ; {

		components := strings.Split(s.Text(),"=")

		switch ( components[0] ) {
		case "b" :
                        i += 1
			results = append(results,Result{})
			results[i-1].Dnf = false
			results[i-1].Bib = components[1]
			break;
	        case "m" :
			results[i-1].Athlete = components[1]
			break;
		case "c" :
			results[i-1].Club = components[1]
			break;
		case "s" :
			results[i-1].Age = components[1]
			break;
		case "un" :
			results[i-1].Ussa = components[1]
			break
		case "r1" :
			if ! strings.HasPrefix(components[1],"D")   {
				f, err := strconv.ParseFloat(components[2],64)
				if err == nil {
					results[i-1].R1 = f
				} else {
					results[i-1].R1 = -1
					results[i-1].Dnf = true
				}
			} else {
				results[i-1].DnfReason = components[1]
				results[i-1].Dnf = true
				results[i-1].R1 = -1
			}
			break
		case "r2" :
			if ! strings.HasPrefix(components[1],"D") {
				f, err := strconv.ParseFloat(components[2],64)
				if err == nil {
					results[i-1].R2 = f
				} else {
					results[i-1].R2 = -1
					results[i-1].Dnf = true
				}
			} else {
				if results[i-1].DnfReason == "" {
					results[i-1].DnfReason = components[1]
				}
				results[i-1].Dnf = true
				results[i-1].R2 = -1
			}
			break
		default :
		}
	}

	return results
}
