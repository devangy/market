package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"path/filepath"
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
	// kalshi_events_API := os.Getenv("kalshi_events_API")
	// poly_events_API := os.Getenv("poly_events_API")
	poly_trades_API := os.Getenv("poly_trades_API")
	// kalshi_trades_API := os.Getenv("kalshi_trades_API")

	proxy, err := url.Parse(fmt.Sprintf("http://user-%s-country-%s:%s@%s", username, country, password, entryPoint))
	if err != nil {
		log.Fatalf("Error proxy parsing %v", err)
	}

	// creating a struct instance using a struct literal in memory
	apiClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	// channel where both api will send the json
	json_chan := make(chan any, 200)

	// Teleg channel for clean and filtered data according to logic applied
	tgEventC := make(chan any, 200)

	// ctx := context.Background()

	go processJson(json_chan, tgEventC)

	// go Bot(tgEventC)

	// go kalshi(kalshi_events_API, apiClient, json_chan)
	// go poly(poly_events_API, apiClient, json_chan)
	go polyTrades(poly_trades_API, apiClient)
	// go kalshiTrades(kalshi_trades_API, apiClient)

	select {}

}

func kalshi(events_API string, apiClient *http.Client, json_chan chan any) {

	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	cursor := ""

	for range ticker.C {
		// ptr := &cursor

		req, err := http.NewRequest("GET", events_API, nil)
		if err != nil {
			log.Fatalf("Err making get req: %v", err)
		}

		// query params
		params := req.URL.Query()
		params.Add("limit", "200")
		params.Add("status", "open")
		params.Add("with_nested_markets", "true")
		params.Add("cursor", cursor)

		req.URL.RawQuery = params.Encode() // form full URL to make call\

		fmt.Println(req.URL.Query())
		res, err := apiClient.Do(req)
		// fmt.Println("url", res.Request.URL)
		if err != nil {
			log.Printf("Err getting res: %v", err)
			continue // Skip this loop iteration, try again next tick
		}

		defer res.Body.Close() // close connection before exiting
		// read from the Body

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Panic("Error reading from res Body:", err)
		}

		// type Market struct {
		// 	// OpenInterest int `json:"open_interest"`
		// 	Liquidity       int    `json:"liquidity"`
		// 	Volume          int    `json:"volume"`
		// 	No_ask_dollars  int    `json:"no_ask"`
		// 	Yes_ask_dollars int    `json:"yes_yes"`
		// 	Status          string `json:"status"`
		// }

		type Event struct {
			Name         string
			Title        string `json:"title"`
			EventTicker  string `json:"event_ticker"`
			SeriesTicker string `json:"series_ticker"`
			Category     string `json:"category"`
		}

		// initial data struct
		type kdump struct {
			Events []Event
			Cursor string `json:"cursor"`
		}
		var kdata kdump

		// if err := json.NewDecoder(res.Body).Decode(&kdata); err != nil {
		// 	log.Fatalf("Error decoding: %v", err)

		if err = json.Unmarshal(body, &kdata); err != nil {
			log.Fatalf("Error unmarshalling: %v", err)
		}

		// if err = json.Unmarshal(kdata , &kdatamain); err != nil {
		// 	log.Fatalf("err unmarshal")
		// }
		fmt.Println("RECEIVED CURSOR:", kdata.Cursor)
		// fmt.Println("Data:", kdata)

		// *ptr = kdata.Cursor
		params.Set("cursor", kdata.Cursor)
		// fmt.Println("ptr", ptr)
		cursor = kdata.Cursor
		fmt.Println("cursor", cursor)

		for _, event := range kdata.Events {
			event.Name = "kalshi"
			json_chan <- event
		}

	}

}

func poly(events_api string, apiClient *http.Client, json_chan chan any) {

	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("GET", events_api, nil)
		if err != nil {
			log.Fatal("err making poly GET request [events]", err)
		}

		// structs

		type polymarketdata struct {
			Name     string
			Title    string  `json:"title"`
			Category string  `json:"category"`
			Volume   float64 `json:"volume"`
			Image    string  `json:"image"`
		}

		params := req.URL.Query()

		params.Add("closed", "false")

		res, err := apiClient.Do(req)
		if err != nil {
			log.Fatal("err getting a res", err)
		}

		// creating a new decoder for incmin json data stream
		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			log.Fatal("err reading body", err)
		}

		// decoder := json.NewDecoder(req.Body)

		var pdata []polymarketdata

		// err = json.NewDecoder(res.Body).Decode(&pdata)

		// if err := decoder.Decode(&pdata); err != nil {
		// 	if err == io.EOF {
		// 		return
		// 	}
		// 	log.Println("decode err", err)
		// 	return
		// }

		if err = json.Unmarshal(body, &pdata); err != nil {
			log.Fatal("err unmarshal poly:", err)
		}

		for _, event := range pdata {
			event.Name = "poly"
			json_chan <- event
		}
	}
}

func processJson(json_chan chan any, tgEventsC chan any) {
	// get current dir path
	//

	directory, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get current directory path", err)
	}

	// create a file or open existing output.jsonl file for writing data
	file, err := os.OpenFile(filepath.Join(directory, "hashes.bin"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Unable opening file for writing", err)
	}

	defer file.Close()

	// empty hashmap for keeping track of seen hashes with struct as it takes 0 bytes for storage and we care about only the key
	seenMap := make(map[uint64]struct{})

	// load already seen hashes into hashmap so that dupes dont get forwarded
	binFile, err := os.Open("hashes.bin")
	if err != nil {
		log.Fatal("failed to open bin file", err)

	}
	defer binFile.Close()

	// channel for sending only freshdata

	for {

		var buffBin [8]byte

		_, err = io.ReadFull(binFile, buffBin[:])
		if err != nil {
			log.Print("err reading bin file:", err)
		}
		if err == io.EOF {
			log.Print("EOF Reached!")
			break
		}

		var value uint64
		value = binary.BigEndian.Uint64(buffBin[:])
		// load hashes in hashmap
		seenMap[value] = struct{}{}

		log.Print("value", value)

		log.Print("buffBin", buffBin)
	}

	// start processing json
	for jdata := range json_chan {

		log.Print("JSONLoop:", jdata)

		// convert incoming json to bytes
		jsonBytes, err := json.Marshal(jdata)
		if err != nil {
			log.Fatal("failed converting to jsonBytes", err)
		}

		// init fnv-1a hashing state object
		fnvH := fnv.New64a()
		// hash each json data coming in
		fnvH.Write(jsonBytes)
		// output hashes
		jsonHashValue := fnvH.Sum64()

		log.Print("JsonHashValue", jsonHashValue)

		// if hash seen before jump to next item
		if _, exists := seenMap[jsonHashValue]; exists {
			log.Print("Duplicate found skipping to next")
			continue
		}

		// send only fresh data to Tg bot
		tgEventsC <- jdata

		// allocate buffer of size 8 bytes
		var buff [8]byte

		log.Print("empty buff", buff)

		// put the hash value in the buffer in BigEndian byte order
		binary.BigEndian.PutUint64(buff[:], jsonHashValue)

		// add hashes to our map to track seen keys
		// The empty struct takes zero bytes of memory. It has no fields, so it holds no data.
		// struct{}{} we care about only if key exists in collection
		// a way of creating a set data type in Go
		seenMap[jsonHashValue] = struct{}{}

		bywritten, err := file.Write(buff[:])
		log.Print("bytes written: ", bywritten)
		if err != nil {
			log.Fatal("failed to write buffer to file", err)
		}

	}
}

func polyTrades(api string, apiClient *http.Client) {
	ticker := time.NewTicker(150 * time.Millisecond)

	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			log.Fatal("Err making request poly [trades]")
		}

		res, err := apiClient.Do(req)
		if err != nil {
			log.Fatal("Failed to get a response", err)
		}
		defer res.Body.Close()

		type Trade struct {
			ProxyWallet           string  `json:"proxyWallet"`
			Side                  string  `json:"side"`
			Asset                 string  `json:"asset"`
			ConditionID           string  `json:"conditionId"`
			Size                  float64 `json:"size"`
			Price                 float64 `json:"price"`
			Timestamp             int64   `json:"timestamp"`
			Title                 string  `json:"title"`
			Slug                  string  `json:"slug"`
			Icon                  string  `json:"icon"`
			EventSlug             string  `json:"eventSlug"`
			Outcome               string  `json:"outcome"`
			OutcomeIndex          int     `json:"outcomeIndex"`
			Name                  string  `json:"name"`
			Pseudonym             string  `json:"pseudonym"`
			Bio                   string  `json:"bio"`
			ProfileImage          string  `json:"profileImage"`
			ProfileImageOptimized string  `json:"profileImageOptimized"`
			TransactionHash       string  `json:"transactionHash"`
		}

		var trades []Trade

		json.NewDecoder(res.Body).Decode(&trades)
		if err != nil {
			log.Fatal("Failed to decode json polytrades", err)
		}

		// prettyJson, _ := json.MarshalIndent(trades, "", "  ")
		// log.Print("ptrades: ", string(prettyJson))

		left := 0
		for right := 0; right < len(trades); right++ {
			// window invalidated remove the trades outside our time window
			for left < right && time.UnixMilli(trades[right].Timestamp).Sub(time.UnixMilli(trades[left].Timestamp)) > 1*time.Minute {
				left++
			}

			value := trades[right].Size * trades[right].Price
			if value >= 2000 {
				prettyJson, _ := json.MarshalIndent(trades[right], "", "  ")
				fmt.Print("LT ❤️:", string(prettyJson))
			}
		}
	}
}

func kalshiTrades(api string, apiClient *http.Client) {
	ticker := time.NewTicker(150 * time.Millisecond)

	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			log.Fatal("Err making request kalshi [trades]")
		}

		res, err := apiClient.Do(req)
		if err != nil {
			log.Fatal("Failed to get a response kalshi [trades]", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Failed to parse body", err)
		}

		log.Print("ktrades: ", string(body))
	}
}
