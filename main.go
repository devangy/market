package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	// env init
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// proxy server
	username := os.Getenv("username")
	password := os.Getenv("password")
	country := os.Getenv("country")
	entryPoint := os.Getenv("entryPoint")

	// Markets

	kalshiMarkets := os.Getenv("kalshi_markets_API")
	//poly_markets_API = os.Getenv("poly_markets_API")

	proxy, err := url.Parse(fmt.Sprintf("http://user-%s-country-%s:%s@%s", username, country, password, entryPoint))
	if err != nil {
		log.Fatalf("Error proxy parsing %v", err)
	}

	// creating a struct instance using a struct literal in memory
	apiClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	var wg sync.WaitGroup // init WaitGroup to wait for goroutine to finish instead of instnt exit
	wg.Add(1)

	go kalshi(kalshiMarkets, apiClient, &wg)

	wg.Wait() // wait for all go routines to finish here at end of all routines

}

func kalshi(apiURL string, apiClient *http.Client, wg *sync.WaitGroup) {
	// new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Err making get req: %v", err)
	}

	// query params
	params := req.URL.Query()
	params.Add("limit", "5")
	params.Add("with_nested_markets", "true")
	req.URL.RawQuery = params.Encode() // form full URL to make call

	res, err := apiClient.Do(req)
	if err != nil {
		log.Fatalf("Err getting res: %v ", err)
	}

	defer res.Body.Close() // close connection before exiting
	defer wg.Done()
	// read from the Body

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading from res Body")
	}

	// unmarshall body into json

	// empty interface will populate based on unmarshall by passing reference &
	//var Kmarket map[string] interface {}

	type Market struct {
		// OpenInterest int `json:"open_interest"`
		Liquidity       int `json:"liquidity"`
		Volume          int `json:"volume"`
		No_ask_dollars  int `json:"no_ask"`
		Yes_ask_dollars int `json:"yes_yes"`
		Status          int `json:"status"`
	}

	type Events struct {
		Title        string   `json:"title"`
		EventTicker  string   `json:"event_ticker"`
		SeriesTicker string   `json:"series_ticker"`
		Category     string   `json:"category"`
		Markets      []Market `json:"markets"`
	}

	// initial data struct
	type kmarketdata struct {
		Events []Events
	}
	var kdata kmarketdata

	// unmarshall
	if err = json.Unmarshal(body, &kdata); err != nil {
		log.Fatalf("Error unmarshalling: %v", err)
	}

	for _, event := range kdata.Events {
		for _, market := range event.Markets {
			fmt.Println("OI:", market.OpenInterest)
		}
	}

	//pretty print json
	prettyjson, err := json.MarshalIndent(kdata, " ", "  ")
	if err != nil {
		log.Fatal("Unable to prettyjson")
		panic(err)
	}

	fmt.Println("res:", string(prettyjson))
}
