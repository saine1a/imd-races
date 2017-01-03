package racedata

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
	"io/ioutil"
	"bufio"
	"strings"
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

	s := bufio.NewScanner(strings.NewReader(string(body)))

	var b []byte
	
	for ; s.Scan() ; {
		str := s.Text()
		for _, tok := range str {
			switch tok {
			case '|':
				b = append(b, "\",\""...)
			case '=':
				b = append(b, "\"=\""...)
			default:
				b = append(b, string(tok)...)
			}
		}
	}


	fmt.Println(string(b))
}
