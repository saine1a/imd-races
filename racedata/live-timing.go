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

func GetRace() {

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
	
	for ; s.Scan() && s.Text() != "endC" ; {

		components := strings.Split(s.Text(),"=")

		fmt.Println(s.Text())
		
		switch ( components[0] ) {
		case "b" :
			fmt.Println("Bib " + components[1])
			break;
		case "m" :
			fmt.Println("Athlete " + components[1])
			break;
		case "c" :
			fmt.Println("Club " + components[1])
			break;
		case "s" :
			fmt.Println("Age " + components[1])
			break;
		case "un" :
			fmt.Println("USSA " + components[1])
			break
		case "r1" :
			if components[1] != "DNF" {
				f, err := strconv.ParseFloat(components[2],64)
				if err == nil {
					fmt.Println("R1 " + strconv.FormatFloat(f / 1000, 'f', -1, 64) )
				} else {
					fmt.Println("R1 bad format")
				}
			} else {
				fmt.Println("R1 DNF")
			}
			break
		case "r2" :
			if components[1] != "DNF" {
				f, err := strconv.ParseFloat(components[2],64)
				if err == nil {
					fmt.Println("R2 " + strconv.FormatFloat(f / 1000, 'f', -1, 64) )
				} else {
					fmt.Println("R2 bad format")
				}
			} else {
				fmt.Println("R2 DNF")
			}
			break
		default :
		}
		
	}
}
