package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Starts the program
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/ws", ws)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	ch := make(chan string)
	for {
		go volume(ch)
		// Write message back to browser
		err := conn.WriteMessage(1, []byte(<-ch))
		time.Sleep(time.Second)
		if err != nil {
			return
		}
	}
}

// Struct for data
type Response struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func volume(ch chan string) {
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
	for i := 0; i < 10; i++ {

		bidP, err := strconv.ParseFloat(res.Bids[i][0], 64)
		bidQ, err := strconv.ParseFloat(res.Bids[i][1], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumBid = bidP*bidQ + sumBid

		askP, err := strconv.ParseFloat(res.Asks[i][0], 64)
		askQ, err := strconv.ParseFloat(res.Asks[i][1], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumAsk = askP*askQ + sumAsk
	}
	ch <- fmt.Sprintf("ID: %v, sumBid = %v, sumAsk = %v", res.LastUpdateID, sumBid, sumBid)
	log.Printf("DATA:\t ID: %v, sumBid = %v, sumAsk = %v\n", res.LastUpdateID, sumBid, sumBid)
}
