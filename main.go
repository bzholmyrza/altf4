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
	"strings"
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
	go volume(ch)
	for {
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

func volume(ch chan string) strings.Builder {
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

	for {
		bid, err := strconv.ParseFloat(res.Bids[0][1], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumBid += bid

		ask, err := strconv.ParseFloat(res.Bids[0][1], 64)
		if err != nil {
			log.Fatal(err)
		}
		sumAsk += ask

		ch <- fmt.Sprintf("ID: %v, sumBid = %v, sumAsk = %v", res.LastUpdateID, sumBid, sumBid)
		log.Printf("DATA:\t ID: %v, sumBid = %v, sumAsk = %v\n", res.LastUpdateID, sumBid, sumBid)
	}
}
