package racedata

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
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

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
