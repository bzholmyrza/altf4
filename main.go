package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Starts the program
func main() {
	doEvery(1 * time.Second, sum)
}

// Make request every second
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// Struct for data
type Response struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// Takes data from api, processes it and displays it in the console
func sum(t time.Time) {
	response, err := http.Get("https://api.binance.com/api/v3/depth?symbol=MANABTC&limit=10")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res Response
	json.Unmarshal(responseData, &res)

	var sumBid float64
	var sumAsk float64
	for i := 0; i < 7; i++ {
		bid, err := strconv.ParseFloat(res.Bids[0][0], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumBid += bid

		ask, err := strconv.ParseFloat(res.Bids[0][0], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumAsk += ask

		fmt.Printf("ID: %v, sumBid = %v, sumAsk = %v\n", res.LastUpdateID, sumBid, sumBid)
	}
	fmt.Println("------------------------------------------")
}
